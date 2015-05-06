// Package images handles the storing and retrieval of raw image data.
package images

import (
	"io"
	"os"
)

// stores maps names to functions that generate a new Store of that type.
// New Store types can register themselves, eliminating the need to hardcode
// new switch cases for new instance creation. The function should just return
// a pointer to a new Store instance, with any connection/configuration handled
// separately via Store.Init().

var stores = map[string]func() Store{}

const configKey = "imageStoreConfig"

type (
	// Store provides a common API for image storage backends
	Store interface {
		// Init handles casting to the appropriate config struct and then
		// performing any connection / initialization needed for the Store
		Init(interface{}) error
		// Shutdown handles disconnection and cleanup for the Store
		Shutdown() error

		// Stat retrieves file information about an image
		Stat(string) (os.FileInfo, error)
		// Get retrieves an image from the Store
		Get(string, io.Writer) error
		// Put stores an image in the Store
		Put(string, io.Reader) error
		// Delete removes an image from the Store
		Delete(string) error
	}
)

// Register adds a new Store type under a name
func Register(name string, newFunc func() Store) {
	stores[name] = newFunc
}

// NewStore create a new instance of a Store from a name
func NewStore(name string) Store {
	newFunc, ok := stores[name]
	if !ok {
		return nil
	}
	return newFunc()
}
