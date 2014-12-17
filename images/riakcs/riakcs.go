package riakcs

import (
    "github.com/mistifyio/mistify-image-service/images"
)

type (
    // Riak-CS image store backend
    RiakCS struct {
        images.ImageStore
    }
)

// Initialize the backend
func (self *RiakCS) Init(config map[string]string) error {
    return nil
}

// Shut down
func (self *RiakCS) Shutdown() error {
    return nil
}

func (self *RiakCS) ListImages() { }
func (self *RiakCS) GetImage() { }
func (self *RiakCS) PutImage() { }
func (self *RiakCS) DeleteImage() { }

