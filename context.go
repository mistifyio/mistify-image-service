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
	if err := ctx.NewImageStore(imageStoreType); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"type":  imageStoreType,
		}).Error("failed to create image store")
		return nil, err
	}

	// json errors would have been caught by viper when loading the file
	imageStoreConfig, _ := json.Marshal(viper.Get("imageStoreConfig"))
	if err := ctx.ImageStore.Init(imageStoreConfig); err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"type":   imageStoreType,
			"config": imageStoreConfig,
		}).Error("failed to initialize image store")
		return nil, err
	}

	// Metadata Storage
	metadataStoreType := viper.GetString("metadataStoreType")
	if err := ctx.NewMetadataStore(metadataStoreType); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"type":  metadataStoreType,
		}).Error("failed to create metadata store")
		return nil, err
	}

	// json errors would have been caught by viper when loading the file
	metadataStoreConfig, _ := json.Marshal(viper.Get("metadataStoreConfig"))
	if err := ctx.MetadataStore.Init(metadataStoreConfig); err != nil {
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
func (ctx *Context) NewImageStore(storeType string) error {
	store := images.NewStore(storeType)
	if store == nil {
		return errors.New("unknown image store type")
	}
	ctx.ImageStore = store
	return nil
}

// NewMetadataStore creates a new metadata store for the context
func (ctx *Context) NewMetadataStore(storeType string) error {
	store := metadata.NewStore(storeType)
	if store == nil {
		return errors.New("unknown metadata store type")
	}
	ctx.MetadataStore = store
	return nil
}
