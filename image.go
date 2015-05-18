package imageservice

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mistifyio/mistify-image-service/metadata"
)

// RegisterImageRoutes registers the image routes and handlers
func RegisterImageRoutes(prefix string, router *mux.Router) {
	router.HandleFunc(prefix, listImagesHandler).Queries("type", "{imageType:[a-zA-Z]+}").Methods("GET")
	router.HandleFunc(prefix, listImagesHandler).Methods("GET")
	router.HandleFunc(prefix, receiveImageHandler).Methods("PUT")
	router.HandleFunc(prefix, fetchImageHandler).Methods("POST")
	sub := router.PathPrefix(prefix).Subrouter()
	sub.HandleFunc("/{imageID}", getImageHandler).Methods("GET")
	sub.HandleFunc("/{imageID}", deleteImageHandler).Methods("DELETE")
	sub.HandleFunc("/{imageID}/download", downloadImageHandler).Methods("GET")
}

// listImagesHandler gets a list of images, optionally filtered by type
func listImagesHandler(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)

	vars := mux.Vars(r)
	imageType := vars["imageType"]
	if imageType != "" && !metadata.IsValidImageType(imageType) {
		hr.JSON(http.StatusBadRequest, "invalid type")
		return
	}

	images, err := ctx.MetadataStore.List(imageType)
	if err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}

	if images == nil {
		images = make([]*metadata.Image, 0)
	}
	hr.JSON(http.StatusOK, images)
}

// receiveImageHandler adds and stores an image from the request body
func receiveImageHandler(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)

	imageType := r.Header.Get("X-Image-Type")
	if !metadata.IsValidImageType(imageType) {
		hr.JSON(http.StatusBadRequest, "invalid X-Image-Type header")
		return
	}

	image, err := ctx.Fetcher.Receive(r)
	if err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}

	hr.JSON(http.StatusOK, image)
}

// fetchImageHandler asynchronously retrieves and adds an image to the system
// from an external source. If image has already been downloaded (same source),
// the existing image data will be returned. Getting the image information
// after a successful fetch has been initiated will show current download
// status.
func fetchImageHandler(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)

	image := &metadata.Image{}
	if err := json.NewDecoder(r.Body).Decode(image); err != nil {
		hr.JSON(http.StatusBadRequest, err.Error())
		return
	}
	image.ID = metadata.NewID()

	// Ensure sufficient information for fetching
	if image.Source == "" {
		hr.JSON(http.StatusBadRequest, "missing image source")
		return
	}
	if !metadata.IsValidImageType(image.Type) {
		hr.JSON(http.StatusBadRequest, "invalid image type")
		return
	}

	image, err := ctx.Fetcher.Fetch(image)
	if err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}
	hr.JSON(http.StatusAccepted, image)
}

// getImageHandler retrieves information about an image.
func getImageHandler(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}

	image := getImage(w, r)
	if image == nil {
		return
	}

	hr.JSON(http.StatusOK, image)
}

// deleteImageHandler removes an image.
func deleteImageHandler(w http.ResponseWriter, r *http.Request) {
	hr := HTTPResponse{w}
	ctx := GetContext(r)

	image := getImage(w, r)
	if image == nil {
		return
	}

	if err := ctx.ImageStore.Delete(image.ID); err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}
	if err := ctx.MetadataStore.Delete(image.ID); err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return
	}

	hr.JSON(http.StatusOK, image)
}

// downloadImageHandler streams an image data
func downloadImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := GetContext(r)

	image := getImage(w, r)
	if image == nil {
		return
	}

	if image.Status != metadata.StatusComplete {
		http.Error(w, "incomplete", 404)
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(image.Size, 10))
	w.Header().Set("Content-Type", "application/octet-stream")
	if err := ctx.ImageStore.Get(image.ID, w); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getImage(w http.ResponseWriter, r *http.Request) *metadata.Image {
	hr := HTTPResponse{w}
	ctx := GetContext(r)
	vars := mux.Vars(r)

	imageID := vars["imageID"]
	image, err := ctx.MetadataStore.GetByID(imageID)
	if err != nil {
		hr.JSONError(http.StatusInternalServerError, err)
		return nil
	}
	if image == nil {
		hr.JSON(http.StatusNotFound, nil)
		return nil
	}

	return image
}
