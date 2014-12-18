package riak

import (
    "github.com/mistifyio/mistify-image-service/metadata"
)

type (

    // Riak metadata store backend
    Riak struct {
        metadata.MetadataStore
    }

)

func (self *Riak) Init(config map[string]string) error {
    return nil
}
