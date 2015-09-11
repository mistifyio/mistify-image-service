package metadata

import (
	"encoding/json"
	"errors"
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	etcderr "github.com/coreos/etcd/error"
	"github.com/coreos/go-etcd/etcd"
)

// ErrIncompleteTLSConfig is used when something is missing from an etcd tls
// configuration
var ErrIncompleteTLSConfig = errors.New("incomplete tls config")

type (
	// etcdStore is a metadata store using etcd
	etcdStore struct {
		client *etcd.Client
		prefix string
		config *EtcdConfig
	}

	// EtcdConfig contains config options to set up an etcd client
	EtcdConfig struct {
		Machines      []string
		Cert          string
		Key           string
		CaCert        string
		Filepath      string
		Prefix        string
		clientNewType string
	}
)

// etcdLogFields contain fields to include on all logs
var etcdLogFields = log.Fields{
	"type":  "metadata",
	"store": "etcd",
}

// Validate checks whether the config is valid and determines what method
// is required to create the new client based on what is provided
func (ec *EtcdConfig) Validate() error {
	if ec.Filepath != "" {
		ec.clientNewType = "file"
		return nil
	}

	// All tls related properties should be empty or all should be defined
	tlsPresent := ec.Cert != "" || ec.Key != "" || ec.CaCert != ""
	tlsMissing := ec.Cert == "" || ec.Key == "" || ec.CaCert == ""
	if tlsPresent {
		if tlsMissing {
			log.WithFields(etcdLogFields).WithFields(log.Fields{
				"error":  ErrIncompleteTLSConfig,
				"cert":   ec.Cert,
				"key":    ec.Key,
				"caCert": ec.CaCert,
			}).Error(ErrIncompleteTLSConfig)
			return ErrIncompleteTLSConfig
		}
		ec.clientNewType = "tls"
	}
	return nil
}

// Init parses the config and creates an etcd client
func (es *etcdStore) Init(configBytes []byte) error {
	config := &EtcdConfig{}

	// Parse the config json
	if err := json.Unmarshal(configBytes, config); err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"json":  string(configBytes),
		}).Error("failed to unmarshal config json")
		return err
	}

	if err := config.Validate(); err != nil {
		return err
	}

	es.config = config
	log.WithFields(etcdLogFields).WithFields(log.Fields{
		"config": es.config,
	}).Info("config loaded")

	es.prefix = path.Join(es.config.Prefix, "images")

	// Create the etcd client
	var client *etcd.Client
	var err error
	switch es.config.clientNewType {
	case "file":
		client, err = etcd.NewClientFromFile(es.config.Filepath)
	case "tls":
		client, err = etcd.NewTLSClient(es.config.Machines, es.config.Cert, es.config.Key, es.config.CaCert)
	default:
		client = etcd.NewClient(es.config.Machines)
	}
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error":  err,
			"config": es.config,
		}).Error("failed to create client")
		return err
	}

	es.client = client

	if _, err := es.client.CreateDir(es.prefix, 0); err != nil {
		etcdErr := err.(*etcd.EtcdError)
		if etcdErr.ErrorCode != etcderr.EcodeNodeExist {
			log.WithFields(etcdLogFields).WithFields(log.Fields{
				"error": err,
				"key":   es.prefix,
			}).Error("failed to create images dir")
			return err
		}
	}
	return nil
}

// Shutdown closes the etcd client connection
func (es *etcdStore) Shutdown() error {
	es.client.Close()
	return nil
}

// List retrieves a list of images from etcd
func (es *etcdStore) List(imageType string) ([]*Image, error) {
	var images []*Image

	// Look up the prefix to get a list of imageIDs
	resp, err := es.client.Get(es.prefix, false, false)
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   es.prefix,
		}).Error("failed to look up images dir")
		return nil, err
	}

	// Look up metadata for each imageID and filter by type
	for _, node := range resp.Node.Nodes {
		imageID := path.Base(node.Key)
		image, err := es.GetByID(imageID)
		if err != nil {
			return nil, err
		}
		if imageType == "" || imageType == image.Type {
			images = append(images, image)
		}
	}

	return images, nil
}

// GetByID retrieves an image from etcd using the image id
func (es *etcdStore) GetByID(imageID string) (*Image, error) {
	image := &Image{}

	metadataKey := es.metadataKey(imageID)
	resp, err := es.client.Get(metadataKey, false, false)
	if err != nil {
		etcdErr := err.(*etcd.EtcdError)
		if etcdErr.ErrorCode != etcderr.EcodeKeyNotFound {
			return nil, ErrNotFound
		}

		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   metadataKey,
		}).Error("failed to look up image")
		return nil, err
	}

	if err := json.Unmarshal([]byte(resp.Node.Value), image); err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   metadataKey,
			"value": resp.Node.Value,
		}).Error("invalid image json")
		return nil, err
	}

	image.Store = es
	return image, nil
}

// GetBySource retrieves an image from etcd using the image source
func (es *etcdStore) GetBySource(imageSource string) (*Image, error) {
	// Look up the prefix to get a list of imageIDs
	resp, err := es.client.Get(es.prefix, false, false)
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   es.prefix,
		}).Error("failed to look up images dir")
		return nil, err
	}

	// Look up metadata for each imageID and return if the right image is found
	for _, node := range resp.Node.Nodes {
		imageID := path.Base(node.Key)
		image, err := es.GetByID(imageID)
		if err != nil {
			return nil, err
		}
		if imageSource == image.Source {
			return image, nil
		}
	}

	return nil, ErrNotFound
}

// Put stores an image in etcd
func (es *etcdStore) Put(image *Image) error {
	imageJSON, err := json.Marshal(image)
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"image": fmt.Sprintf("%+v", image),
		}).Error("failed to marshal image to json")
		return err
	}

	metadataKey := es.metadataKey(image.ID)
	if _, err := es.client.Set(metadataKey, string(imageJSON), 0); err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   metadataKey,
			"value": string(imageJSON),
		}).Error("failed to store image")
		return err
	}

	return nil
}

// Delete removs an image from etcd
func (es *etcdStore) Delete(imageID string) error {
	key := path.Join(es.prefix, imageID)
	if _, err := es.client.Delete(key, true); err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).Error("failed to delete image")
		return err
	}

	return nil
}

func (es *etcdStore) metadataKey(imageID string) string {
	return path.Join(es.prefix, imageID, "metadata")
}

func init() {
	Register("etcd", func() Store {
		return &etcdStore{}
	})
}
