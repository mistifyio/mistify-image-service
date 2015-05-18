package imageservice

import (
	"encoding/json"
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/spf13/viper"
)

type (
	// Context holds the initialized stores
	Context struct {
		ImageStore    images.Store
		MetadataStore metadata.Store
		Fetcher       *Fetcher
	}
)

// NewContext creates a new context from configuration
func NewContext() (*Context, error) {
	ctx := &Context{}

	// Image Storage
	imageStoreType := viper.GetString("imageStoreType")
	// json errors would have been caught by viper when loading the file
	imageStoreConfig, _ := json.Marshal(viper.Get("imageStoreConfig"))
	if err := ctx.InitImageStore(imageStoreType, imageStoreConfig); err != nil {
		return nil, err
	}

	// Metadata Storage
	metadataStoreType := viper.GetString("metadataStoreType")
	// json errors would have been caught by viper when loading the file
	metadataStoreConfig, _ := json.Marshal(viper.Get("metadataStoreConfig"))
	if err := ctx.InitMetadataStore(metadataStoreType, metadataStoreConfig); err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"type":   metadataStoreType,
			"config": metadataStoreConfig,
		}).Error("failed to initialize metadata store")
		return nil, err
	}

	// Image Fetcher
	ctx.Fetcher = NewFetcher(ctx)

	return ctx, nil
}

// NewImageStore creates a new image store for the context
func (ctx *Context) InitImageStore(storeType string, configBytes []byte) error {
	store := images.NewStore(storeType)
	if store == nil {
		err := errors.New("unknown image store type")
		log.WithFields(log.Fields{
			"error": err,
			"type":  storeType,
		}).Error("failed to create image store")
		return err
	}

	if err := store.Init(configBytes); err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"type":   storeType,
			"config": string(configBytes),
		}).Error("failed to initialize image store")
		return err
	}

	ctx.ImageStore = store

	return nil
}

// NewMetadataStore creates a new metadata store for the context
func (ctx *Context) InitMetadataStore(storeType string, configBytes []byte) error {
	store := metadata.NewStore(storeType)
	if store == nil {
		err := errors.New("unknown metadata store type")
		log.WithFields(log.Fields{
			"error": err,
			"type":  storeType,
		}).Error("failed to create metadata store")
		return err
	}

	if err := store.Init(configBytes); err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"type":   storeType,
			"config": string(configBytes),
		}).Error("failed to initialize metadata store")
		return err
	}

	ctx.MetadataStore = store

	return nil
}
