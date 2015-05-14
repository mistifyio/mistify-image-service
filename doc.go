/*
Package imageservice is the mistify guest image server. In order to remove
dependence on external sources, which may be unavailable or tampered with, a
mistify-agent hypervisor will instead fetch images from the
mistify-image-service. An operator will load images into mistify-image-service
creating by either direct upload or by having the service fetch an image from
an external source over http.

HTTP API Endpoints

	/images
		* GET  - Retrieve a list of images, optionally filtered by type.
		* POST - Fetch and store an image
		* PUT  - Upload and store image

	/images/{imageID}
		* GET    - Retrieves information for an image
		* DELETE - Deletes an image

	/images/{imageID}/download
		* GET - Download an image

Image information uses the metadata.Image struct.  When directly uploading an
image, the body should be the raw image data, with the image type and optional
comment provided via headers X-Image-Type and X-Image-Comment, respectively.
*/
package imageservice
