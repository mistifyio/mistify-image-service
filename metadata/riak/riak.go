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
