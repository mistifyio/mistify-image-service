package imageservice

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/mistifyio/mistify-image-service/metadata"
)

// RegisterImageRoutes registers the image routes and handlers
func RegisterImageRoutes(prefix string, router *mux.Router) {
	router.HandleFunc(prefix, ListImages).Queries("type", "{imageType:[a-zA-Z]}").Methods("GET")
	router.HandleFunc(prefix, ListImages).Methods("GET")
	//router.HandleFunc(prefix, UploadImage).Methods("PUT")
	router.HandleFunc(prefix, FetchImage).Methods("POST")
	sub := router.PathPrefix(prefix).Subrouter()
	sub.HandleFunc("/{imageID}", GetImage).Methods("GET")
	sub.HandleFunc("/{imageID}", DeleteImage).Methods("DELETE")
	sub.HandleFunc("/{imageID}/download", DownloadImage).Methods("GET")
}

// ListImages gets a list of images, optionally filtered by type
func ListImages(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)
	vars := mux.Vars(r)

	images, err := ctx.MetadataStore.List(vars["imageType"])
	if err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}
	hr.JSON(http.StatusOK, images)
}

// FetchImage asynchronously retrieves and adds an image to the system from an
// external source. If image has already been downloaded (same source), the
// existing image data will be returned. Getting the image information after a
// successful fetch has been initiated will show current download status.
func FetchImage(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)

	image := &metadata.Image{}
	if err := json.NewDecoder(r.Body).Decode(image); err != nil {
		hr.JSON(http.StatusBadRequest, err.Error())
		return
	}

	image, err := ctx.Fetcher.Fetch(image)
	if err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}
	hr.JSON(http.StatusAccepted, image)
}

// GetImage retrieves information about an image.
func GetImage(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}

	image, err := getImage(w, r)
	if err != nil {
		return
	}

	hr.JSON(http.StatusOK, image)
}

// DeleteImage removes an image.
func DeleteImage(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)

	image, err := getImage(w, r)
	if err != nil {
		return
	}

	if err := ctx.ImageStore.Delete(image.ID); err != nil && !images.IsErrNotFound(err) {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}
	if err := ctx.MetadataStore.Delete(image.ID); err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}

	hr.JSON(http.StatusOK, image)
}

// DownloadImage streams an image data
func DownloadImage(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	image, err := getImage(w, r)
	if err != nil {
		return
	}

	w.Header().Set("Content-Length", string(image.Size))
	w.Header().Set("Content-Type", "application/octet-stream")
	if err := ctx.ImageStore.Get(image.ID, w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getImage(w http.ResponseWriter, r *http.Request) (*metadata.Image, error) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)
	vars := mux.Vars(r)

	imageID := vars["imageID"]
	image, err := ctx.MetadataStore.GetByID(imageID)
	if err != nil {
		if metadata.IsErrNotFound(err) {
			hr.JSON(http.StatusNotFound, err.Error())
			return nil, err
		}
		hr.JSONError(http.StatusInternalServerError, err)
		return nil, err
	}
	return image, nil
}
