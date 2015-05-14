# metadata

[![metadata](https://godoc.org/github.com/mistifyio/mistify-image-service/metadata?status.png)](https://godoc.org/github.com/mistifyio/mistify-image-service/metadata)

Package metadata handles the storing and retrieval of image metadata.

## Usage

```go
const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusComplete    = "complete"
	StatusError       = "error"
)
```
Image statuses

```go
const (
	ImageTypeKVM       = "kvm"
	ImageTypeContainer = "container"
)
```
Valid image types

```go
var ValidImageTypes = map[string]struct{}{
	ImageTypeKVM:       struct{}{},
	ImageTypeContainer: struct{}{},
}
```
Map of valid image types for quick lookups

#### func  IsValidImageType

```go
func IsValidImageType(imageType string) bool
```
IsValidImageType tests whether the image type is valid

#### func  Register

```go
func Register(name string, newFunc func() Store)
```
Register adds a new Store under a name

#### type Image

```go
type Image struct {
	ID            string
	Source        string
	Type          string
	Comment       string
	Status        string
	Size          int64
	ExpectedSize  int64
	DownloadStart time.Time
	DownloadEnd   time.Time
	Store         Store `json:"-"`
}
```

Image is metadata for an image

#### func (*Image) NewID

```go
func (image *Image) NewID()
```
NewID generates a new unique ID for the image

#### func (*Image) SetDownloading

```go
func (image *Image) SetDownloading(size int64) error
```
SetDownloading updates an image to downloading status with estimated size

#### func (*Image) SetFinished

```go
func (image *Image) SetFinished(err error) error
```
SetFinished updates an image to the final status

#### func (*Image) SetPending

```go
func (image *Image) SetPending() error
```
SetPending updates an image to pending status

#### func (*Image) UpdateSize

```go
func (image *Image) UpdateSize(size int64) error
```
UpdateSize upates an image's current size

#### type KVite

```go
type KVite struct {
	Config *KViteConfig
}
```

KVite is a metadata store using kvite

#### func (*KVite) Delete

```go
func (kv *KVite) Delete(imageID string) error
```
Delete removes an image from kvite

#### func (*KVite) GetByID

```go
func (kv *KVite) GetByID(imageID string) (*Image, error)
```
GetByID retrieves an image from kvite using the image id

#### func (*KVite) GetBySource

```go
func (kv *KVite) GetBySource(imageSource string) (*Image, error)
```
GetBySource retrieves an image from kvite using the image source

#### func (*KVite) Init

```go
func (kv *KVite) Init(rawConfig interface{}) error
```
Init parses the config and opens a connection to kvite

#### func (*KVite) List

```go
func (kv *KVite) List(imageType string) ([]*Image, error)
```
List retrieves a list of images from kvite

#### func (*KVite) Put

```go
func (kv *KVite) Put(image *Image) error
```
Put stores an image in kvite

#### func (*KVite) Shutdown

```go
func (kv *KVite) Shutdown() error
```
Shutdown closes the connection to kvite

#### type KViteConfig

```go
type KViteConfig struct {
	Filename string
	Table    string
}
```

KViteConfig contains necessary config options to set up kvite

#### func (*KViteConfig) Validate

```go
func (kvc *KViteConfig) Validate() error
```
Validate checks whether the config is valid

#### type Store

```go
type Store interface {
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
```

Store provides a common API for image storage backends

#### func  NewStore

```go
func NewStore(name string) Store
```
NewStore creates a new instance of a Store from a name

--
*Generated with [godocdown](https://github.com/robertkrimen/godocdown)*