package metadata

import (
	"time"

	"code.google.com/p/go-uuid/uuid"
)

// Image statuses
const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusComplete    = "complete"
	StatusError       = "error"
)

type (
	// Image is metadata for an image
	Image struct {
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
)

// NewID generates a new unique ID for the image
func (image *Image) NewID() {
	image.ID = uuid.New()
}

// SetPending updates an image to pending status
func (image *Image) SetPending() error {
	image.Status = StatusPending
	return image.Store.Put(image)
}

// SetDownloading updates an image to downloading status with estimated size
func (image *Image) SetDownloading(size int64) error {
	image.Status = StatusDownloading
	image.DownloadStart = time.Now()
	image.ExpectedSize = size
	return image.Store.Put(image)
}

// UpdateSize upates an image's current size
func (image *Image) UpdateSize(size int64) error {
	image.Size = size
	return image.Store.Put(image)
}

// SetFinished updates an image to the final status
func (image *Image) SetFinished(err error) error {
	if err != nil {
		image.Status = StatusError
	} else {
		image.Status = StatusComplete
	}

	image.DownloadEnd = time.Now()
	return image.Store.Put(image)
}
