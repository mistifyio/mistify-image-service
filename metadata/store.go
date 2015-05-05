// Package metadata handles the storing and retrieval of image metadata
package metadata

// stores maps names to functions that generate a new Store of that type.
// New Store types can register themselves, eliminating the need to hardcode
// new switch cases for new instance creation. The function should just return
// a pointer to a new Store instance, with any connection/configuration handled
// separately via Store.Init().
var stores = map[string]func() Store{}

type (
	// Store provides a common API for image storage backends
	Store interface {
		// Init handles casting to the appropriate config struct and then
		// performing any connection / initialization needed for the Store
		Init(interface{}) error
		// Shutdown handles disconnection and cleanup for the Store
		Shutdown() error

		// List retrieves a list of metadata for all available images,
		// optionally filtered by type.
		List(string) ([]*Image, error)
		// GetByID retrieves metadata for an image from the Store by ID
		GetByID(string) (*Image, error)
		// GetBySource retrieves metadata for an image from the Store by source
		GetBySource(string) (*Image, error)
		// Put stores metadata for an image form the Store
		Put(*Image) error
		// Delete removes metadata for an image from the Store
		Delete(string) error
	}
)

// Register adds a new Store under a name
func Register(name string, newFunc func() Store) {
	stores[name] = newFunc
}

// NewStore creates a new instance of a Store from a name
func NewStore(name string) Store {
	newFunc, ok := stores[name]
	if !ok {
		return nil
	}
	return newFunc()
}
