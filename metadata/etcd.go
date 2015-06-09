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

type (
	// Etcd is a metadata store using etcd
	Etcd struct {
		client *etcd.Client
		Prefix string
		Config *EtcdConfig
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

// Validate checks whether the config is valid and
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
			return errors.New("incomplete tls config")
		}
		ec.clientNewType = "tls"
	}
	return nil
}

// Init parses the config and creates an etcd client
func (ec *Etcd) Init(configBytes []byte) error {
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
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
		}).Error("failed config validation")

		return err
	}

	ec.Config = config
	log.WithFields(etcdLogFields).WithFields(log.Fields{
		"config": ec.Config,
	}).Info("config loaded")

	ec.Prefix = path.Join(ec.Config.Prefix, "images")

	// Create the etcd client
	var client *etcd.Client
	var err error
	switch ec.Config.clientNewType {
	case "file":
		client, err = etcd.NewClientFromFile(ec.Config.Filepath)
	case "tls":
		client, err = etcd.NewTLSClient(ec.Config.Machines, ec.Config.Cert, ec.Config.Key, ec.Config.CaCert)
	default:
		client = etcd.NewClient(ec.Config.Machines)
	}
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error":  err,
			"config": ec.Config,
		}).Error("failed to create client")
		return err
	}

	ec.client = client

	if _, err := ec.client.CreateDir(ec.Prefix, 0); err != nil {
		etcdErr := err.(*etcd.EtcdError)
		if etcdErr.ErrorCode != etcderr.EcodeNodeExist {
			log.WithFields(etcdLogFields).WithFields(log.Fields{
				"error": err,
				"key":   ec.Prefix,
			}).Error("failed to create images dir")
			return err
		}
	}
	return nil
}

// Shutdown closes the etcd client connection
func (ec *Etcd) Shutdown() error {
	ec.client.Close()
	return nil
}

// List retrieves a list of images from etcd
func (ec *Etcd) List(imageType string) ([]*Image, error) {
	var images []*Image

	// Look up the prefix to get a list of imageIDs
	resp, err := ec.client.Get(ec.Prefix, false, false)
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   ec.Prefix,
		}).Error("failed to look up images dir")
		return nil, err
	}

	// Look up metadata for each imageID and filter by type
	for _, node := range resp.Node.Nodes {
		imageID := path.Base(node.Key)
		image, err := ec.GetByID(imageID)
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
func (ec *Etcd) GetByID(imageID string) (*Image, error) {
	image := &Image{}

	metadataKey := ec.metadataKey(imageID)
	resp, err := ec.client.Get(metadataKey, false, false)
	if err != nil {
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

	image.Store = ec
	return image, nil
}

// GetBySource retrieves an image from etcd using the image source
func (ec *Etcd) GetBySource(imageSource string) (*Image, error) {
	// Look up the prefix to get a list of imageIDs
	resp, err := ec.client.Get(ec.Prefix, false, false)
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   ec.Prefix,
		}).Error("failed to look up images dir")
		return nil, err
	}

	// Look up metadata for each imageID and return if the right image is found
	for _, node := range resp.Node.Nodes {
		imageID := path.Base(node.Key)
		image, err := ec.GetByID(imageID)
		if err != nil {
			return nil, err
		}
		if imageSource == image.Source {
			return image, nil
		}
	}

	return nil, nil
}

// Put stores an image in etcd
func (ec *Etcd) Put(image *Image) error {
	imageJSON, err := json.Marshal(image)
	if err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"image": fmt.Sprintf("%+v", image),
		}).Error("failed to marshal image to json")
		return err
	}

	metadataKey := ec.metadataKey(image.ID)
	if _, err := ec.client.Set(metadataKey, string(imageJSON), 0); err != nil {
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
func (ec *Etcd) Delete(imageID string) error {
	key := path.Join(ec.Prefix, imageID)
	if _, err := ec.client.Delete(key, true); err != nil {
		log.WithFields(etcdLogFields).WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).Error("failed to delete image")
		return err
	}

	return nil
}

func (ec *Etcd) metadataKey(imageID string) string {
	return path.Join(ec.Prefix, imageID, "metadata")
}

func init() {
	Register("etcd", func() Store {
		return &Etcd{}
	})
}
