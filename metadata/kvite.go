package metadata

import (
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/kvite"
)

type (
	// KVite is a metadata store using kvite
	KVite struct {
		db     *kvite.DB
		Config *KViteConfig
	}

	// KViteConfig contains necessary config options to set up kvite
	KViteConfig struct {
		Filename string
		Table    string
	}
)

// kviteLogFields contain fields to include on all logs
var kviteLogFields = log.Fields{
	"type":  "metadata",
	"store": "kvite",
}

const kviteBucket = "images"

// Validate checks whether the config is valid
func (kvc *KViteConfig) Validate() error {
	if kvc.Filename == "" {
		return errors.New("empty filename")
	}
	if kvc.Table == "" {
		return errors.New("empty table")
	}
	return nil
}

// Init parses the config and opens a connection to kvite
func (kv *KVite) Init(configBytes []byte) error {
	config := &KViteConfig{}

	// Parse the config json
	if err := json.Unmarshal(configBytes, config); err != nil {
		log.WithFields(kviteLogFields).WithFields(log.Fields{
			"error": err,
			"json":  string(configBytes),
		}).Error("failed to unmarshal config json")
		return err
	}

	if err := config.Validate(); err != nil {
		log.WithFields(kviteLogFields).WithFields(log.Fields{
			"error": err,
		}).Error("failed config validation")

		return err
	}

	kv.Config = config
	log.WithFields(kviteLogFields).WithFields(log.Fields{
		"config": kv.Config,
	}).Info("config loaded")

	// Create the kvite database connection
	db, err := kvite.Open(kv.Config.Filename, kv.Config.Table)
	if err != nil {
		log.WithFields(kviteLogFields).WithFields(log.Fields{
			"error":  err,
			"config": kv.Config,
		}).Error("failed to open db connection")
		return err
	}

	kv.db = db
	return nil
}

// Shutdown closes the connection to kvite
func (kv *KVite) Shutdown() error {
	if err := kv.db.Close(); err != nil {
		log.WithFields(kviteLogFields).WithFields(log.Fields{
			"error": err,
		}).Error("failed to close db connection")
		return err
	}
	return nil
}

// List retrieves a list of images from kvite
func (kv *KVite) List(imageType string) ([]*Image, error) {
	var images []*Image
	err := kv.db.Transaction(func(tx *kvite.Tx) error {
		// Setup the bucket
		bucket, err := kv.bucketSetup(tx)
		if bucket == nil || err != nil {
			return err
		}

		// Parse each image json and append to the images array
		return bucket.ForEach(func(key string, value []byte) error {
			image := &Image{}
			if err := json.Unmarshal(value, image); err != nil {
				log.WithFields(kviteLogFields).WithFields(log.Fields{
					"error":  err,
					"bucket": kviteBucket,
					"key":    key,
					"value":  string(value),
				}).Error("failed to parse image json")
				return err
			}
			if imageType == "" || image.Type == imageType {
				images = append(images, image)
			}
			return nil
		})
	})

	return images, err
}

// GetByID retrieves an image from kvite using the image id
func (kv *KVite) GetByID(imageID string) (*Image, error) {
	var image Image
	err := kv.db.Transaction(func(tx *kvite.Tx) error {
		// Setup the bucket
		bucket, err := kv.bucketSetup(tx)
		if bucket == nil || err != nil {
			return err
		}

		value, err := bucket.Get(imageID)
		if err != nil {
			log.WithFields(kviteLogFields).WithFields(log.Fields{
				"error":   err,
				"imageID": imageID,
			}).Error("failed to retrieve image")
		}
		if value == nil {
			return nil
		}

		if err := json.Unmarshal(value, &image); err != nil {
			log.WithFields(kviteLogFields).WithFields(log.Fields{
				"error":  err,
				"bucket": kviteBucket,
				"key":    imageID,
				"value":  string(value),
			}).Error("failed to parse image json")
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if image.ID == "" {
		return nil, nil
	}
	return &image, nil
}

// GetBySource retrieves an image from kvite using the image source
func (kv *KVite) GetBySource(imageSource string) (*Image, error) {
	var foundImage Image
	err := kv.db.Transaction(func(tx *kvite.Tx) error {
		// Setup the bucket
		bucket, err := kv.bucketSetup(tx)
		if bucket == nil || err != nil {
			return err
		}

		// Go through the image and check the sources
		// Is there a more efficient way with kvite?
		return bucket.ForEach(func(key string, value []byte) error {
			image := &Image{}
			if err := json.Unmarshal(value, image); err != nil {
				log.WithFields(kviteLogFields).WithFields(log.Fields{
					"error":  err,
					"bucket": kviteBucket,
					"key":    key,
					"value":  string(value),
				}).Error("failed to parse image json")
				return err
			}
			if image.Source == imageSource {
				foundImage = *image
			}
			return nil
		})
	})

	if err != nil {
		return nil, err
	}
	if foundImage.ID == "" {
		return nil, nil
	}
	return &foundImage, nil
}

// Put stores an image in kvite
func (kv *KVite) Put(image *Image) error {
	err := kv.db.Transaction(func(tx *kvite.Tx) error {
		// Setup the bucket
		bucket, err := kv.bucketSetup(tx)
		if bucket == nil || err != nil {
			return err
		}

		value, err := json.Marshal(image)
		if err != nil {
			log.WithFields(kviteLogFields).WithFields(log.Fields{
				"error": err,
				"image": image,
			}).Error("failed to marshal image")
		}

		if err := bucket.Put(image.ID, value); err != nil {
			log.WithFields(kviteLogFields).WithFields(log.Fields{
				"error": err,
				"key":   image.ID,
				"value": string(value),
			}).Error("failed to store image")
			return err
		}
		return nil
	})
	return err
}

// Delete removes an image from kvite
func (kv *KVite) Delete(imageID string) error {
	err := kv.db.Transaction(func(tx *kvite.Tx) error {
		// Setup the bucket
		bucket, err := kv.bucketSetup(tx)
		if bucket == nil || err != nil {
			return err
		}

		if err := bucket.Delete(imageID); err != nil {
			log.WithFields(kviteLogFields).WithFields(log.Fields{
				"error": err,
				"key":   imageID,
			}).Error("failed to deleteimage")
			return err
		}
		return nil
	})
	return err
}

// bucketSetup gets a kvite bucket and logs any issues/errors
func (kv *KVite) bucketSetup(tx *kvite.Tx) (*kvite.Bucket, error) {
	// Setup the bucket
	bucket, err := tx.Bucket(kviteBucket)
	if err != nil {
		log.WithFields(kviteLogFields).WithFields(log.Fields{
			"error":  err,
			"bucket": kviteBucket,
		}).Error("failed to retrieve bucket")
		return nil, err
	}
	if bucket == nil {
		log.WithFields(kviteLogFields).WithFields(log.Fields{
			"bucket": kviteBucket,
		}).Info("bucket does not exist")
	}
	return bucket, err
}

func init() {
	Register("kvite", func() Store {
		return &KVite{}
	})
}
