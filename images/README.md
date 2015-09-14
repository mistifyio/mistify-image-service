# images

[![images](https://godoc.org/github.com/mistifyio/mistify-image-service/images?status.png)](https://godoc.org/github.com/mistifyio/mistify-image-service/images)

Package images handles the storing and retrieval of raw image data.

## Usage

```go
var ErrMissingDir = errors.New("missing dir")
```
ErrMissingDir is used when the required dir is omitted from the config

#### func  List

```go
func List() []string
```
List registered store names

#### func  Register

```go
func Register(name string, newFunc func() Store)
```
Register adds a new Store type under a name

#### type FS

```go
type FS struct {
	Config *FSConfig
}
```

FS is an image store using the filesystem

#### func (*FS) Delete

```go
func (fs *FS) Delete(imageID string) error
```
Delete removes an image from the filesystem

#### func (*FS) Get

```go
func (fs *FS) Get(imageID string, out io.Writer) error
```
Get retrieves an image from the filesystem

#### func (*FS) Init

```go
func (fs *FS) Init(configBytes []byte) error
```
Init parses the config and ensures the directory exists

#### func (*FS) Put

```go
func (fs *FS) Put(imageID string, in io.Reader) error
```
Put stores an image in the filesystem

#### func (*FS) Shutdown

```go
func (fs *FS) Shutdown() error
```
Shutdown is a noop

#### func (*FS) Stat

```go
func (fs *FS) Stat(imageID string) (os.FileInfo, error)
```
Stat retrieves file information about an image

#### type FSConfig

```go
type FSConfig struct {
	Dir string
}
```

FSConfig contains necessary config options to set up the fs store

#### func (*FSConfig) Validate

```go
func (fsc *FSConfig) Validate() error
```
Validate checks whether the config is valid

#### type Store

```go
type Store interface {
	// Init handles casting to the appropriate config struct and then
	// performing any connection / initialization needed for the Store
	Init([]byte) error
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
```

Store provides a common API for image storage backends

#### func  NewStore

```go
func NewStore(name string) Store
```
NewStore create a new instance of a Store from a name

--
*Generated with [godocdown](https://github.com/robertkrimen/godocdown)*
