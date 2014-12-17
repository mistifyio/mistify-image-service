package imageservice

import (
    "github.com/mistifyio/mistify-image-service/images"
    "github.com/mistifyio/mistify-image-service/images/riakcs"

    "github.com/mistifyio/mistify-image-service/metadata"
    "github.com/mistifyio/mistify-image-service/metadata/riak"
    
    "fmt"
)

type (
    
    // context is the core of everything
    Context struct {
        ImagesBackend images.ImageStore
        MetadataBackend metadata.MetadataStore
	}
    
)

const (
    // imagestore backend types 
    
    // Riak-CS backend
    IMAGE_STORE_RIAKCS = "riakcs"
    
    // Riak metadata backend
    METADATA_STORE_RIAK = "riak"
)

// Create a new context from configuration
func NewContext(config *Config) (*Context, error) {
    // create context
    ctx := &Context{}
    
    // create backends
    err := ctx.createImageStore(config.ImageStoreType)
    if nil != err {
        return nil, err
    }
    err = ctx.createMetadataStore(config.MetadataStoreType)
    if nil != err {
        return nil, err
    }
    
    // initialize the backends
    err = ctx.ImagesBackend.Init(config.ImageStoreConfig)
    if nil != err {
        return nil, err
    }
    err = ctx.MetadataBackend.Init(config.MetadataStoreConfig)
    if nil != err {
        return nil, err
    }
    
    return ctx, nil
}

// Create an imagestore backend instance by type
func (c *Context) createImageStore(imageStoreType string) error {
    switch imageStoreType {
        case IMAGE_STORE_RIAKCS:
            c.ImagesBackend = new(riakcs.RiakCS);
        default:
            return fmt.Errorf("Invalid image store backend type: %s", imageStoreType)
    }
    
    return nil
}

// Create a metadata store backend by type
func (c *Context) createMetadataStore(metadataStoreType string) error {
    switch metadataStoreType {
        case METADATA_STORE_RIAK:
            c.MetadataBackend = new(riak.Riak);
        default:
            return fmt.Errorf("Invalid metadata store backend type: %s", metadataStoreType)
    }
    
    return nil
}
