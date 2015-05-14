# imageservice

[![imageservice](https://godoc.org/github.com/mistifyio/mistify-image-service?status.png)](https://godoc.org/github.com/mistifyio/mistify-image-service)

Package imageservice is the mistify guest image server. In order to remove
dependence on external sources, which may be unavailable or tampered with, a
mistify-agent hypervisor will instead fetch images from the
mistify-image-service. An operator will load images into mistify-image-service
creating by either direct upload or by having the service fetch an image from an
external source over http.

### HTTP API Endpoints

    /images
    	* GET  - Retrieve a list of images, optionally filtered by type.
    	* POST - Fetch and store an image
    	* PUT  - Upload and store image

    /images/{imageID}
    	* GET    - Retrieves information for an image
    	* DELETE - Deletes an image

    /images/{imageID}/download
    	* GET - Download an image

Image information uses the metadata.Image struct. When directly uploading an
image, the body should be the raw image data, with the image type and optional
comment provided via headers X-Image-Type and X-Image-Comment, respectively.

## Usage

#### func  DeleteImage

```go
func DeleteImage(w http.ResponseWriter, r *http.Request)
```
DeleteImage removes an image.

#### func  DownloadImage

```go
func DownloadImage(w http.ResponseWriter, r *http.Request)
```
DownloadImage streams an image data

#### func  FetchImage

```go
func FetchImage(w http.ResponseWriter, r *http.Request)
```
FetchImage asynchronously retrieves and adds an image to the system from an
external source. If image has already been downloaded (same source), the
existing image data will be returned. Getting the image information after a
successful fetch has been initiated will show current download status.

#### func  GetImage

```go
func GetImage(w http.ResponseWriter, r *http.Request)
```
GetImage retrieves information about an image.

#### func  ListImages

```go
func ListImages(w http.ResponseWriter, r *http.Request)
```
ListImages gets a list of images, optionally filtered by type

#### func  RegisterImageRoutes

```go
func RegisterImageRoutes(prefix string, router *mux.Router)
```
RegisterImageRoutes registers the image routes and handlers

#### func  Run

```go
func Run(ctx *Context, port int) error
```
Run starts the server

#### func  SetContext

```go
func SetContext(r *http.Request, ctx *Context)
```
SetContext sets a Context value for a request

#### func  UploadImage

```go
func UploadImage(w http.ResponseWriter, r *http.Request)
```
UploadImage adds and stores an image from the request body

#### type Context

```go
type Context struct {
	ImageStore    images.Store
	MetadataStore metadata.Store
	Fetcher       *Fetcher
}
```

Context holds the initialized stores

#### func  GetContext

```go
func GetContext(r *http.Request) *Context
```
GetContext retrieves a Context value for a request

#### func  NewContext

```go
func NewContext() (*Context, error)
```
NewContext creates a new context from configuration

#### func (*Context) NewFetcher

```go
func (ctx *Context) NewFetcher()
```
NewFetcher creates a new image fetcher for the context

#### func (*Context) NewImageStore

```go
func (ctx *Context) NewImageStore(storeType string) error
```
NewImageStore creates a new image store for the context

#### func (*Context) NewMetadataStore

```go
func (ctx *Context) NewMetadataStore(storeType string) error
```
NewMetadataStore creates a new metadata store for the context

#### type Fetcher

```go
type Fetcher struct {
}
```

Fetcher handles fetching new images and updating metadata accordingly

#### func  NewFetcher

```go
func NewFetcher(ctx *Context) *Fetcher
```
NewFetcher creates a new Fetcher

#### func (*Fetcher) Fetch

```go
func (fetcher *Fetcher) Fetch(image *metadata.Image) (*metadata.Image, error)
```
Fetch runs pre-flight checks and kicks off an asynchronous image download

#### func (*Fetcher) Upload

```go
func (fetcher *Fetcher) Upload(r *http.Request) (*metadata.Image, error)
```
Upload uploads an image synchronously

#### type HTTPError

```go
type HTTPError struct {
	Message string   `json:"message"`
	Code    int      `json:"code"`
	Stack   []string `json:"stack"`
}
```

HTTPError contains information for http error responses

#### type HTTPResponse

```go
type HTTPResponse struct {
	http.ResponseWriter
}
```

HTTPResponse is a wrapper for http.ResponseWriter which provides access to
several convenience methods

#### func (*HTTPResponse) JSON

```go
func (hr *HTTPResponse) JSON(code int, obj interface{})
```
JSON writes appropriate headers and JSON body to the http response

#### func (*HTTPResponse) JSONError

```go
func (hr *HTTPResponse) JSONError(code int, err error)
```
JSONError prepares an HTTPError with a stack trace and writes it with
HTTPResponse.JSON

#### func (*HTTPResponse) JSONMsg

```go
func (hr *HTTPResponse) JSONMsg(code int, msg string)
```
JSONMsg is a convenience method to write a JSON response with just a message
string

--
*Generated with [godocdown](https://github.com/robertkrimen/godocdown)*
