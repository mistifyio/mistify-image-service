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
    	* POST - Fetch and store an image from an external http source
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

#### func (*Fetcher) Receive

```go
func (fetcher *Fetcher) Receive(r *http.Request) (*metadata.Image, error)
```
Receive adds and saves an image synchronously from the request body

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
