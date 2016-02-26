package metadata

import (
	"time"

	"github.com/pborman/uuid"
)

// Image statuses
const (
	StatusPending     = "pending"
	StatusDownloading = "downloading"
	StatusComplete    = "complete"
	StatusError       = "error"
)

// Valid image types
const (
	ImageTypeKVM       = "kvm"
	ImageTypeContainer = "container"
)

// ValidImageTypes is a map of valid image types for quick lookups
var ValidImageTypes = map[string]struct{}{
	ImageTypeKVM:       {},
	ImageTypeContainer: {},
}

type (
	// Image is metadata for an image
	Image struct {
		ID            string    `json:"id"`
		Source        string    `json:"source"`
		Type          string    `json:"type"`
		Comment       string    `json:"comment"`
		Status        string    `json:"status"`
		Size          int64     `json:"size"`
		ExpectedSize  int64     `json:"expected_size"`
		DownloadStart time.Time `json:"download_start"`
		DownloadEnd   time.Time `json:"download_end"`
		Store         Store     `json:"-"`
	}
)

// NewID generates a new unique uuid
func NewID() string {
	return uuid.New()
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

// IsValidImageType tests whether the image type is valid
func IsValidImageType(imageType string) bool {
	_, ok := ValidImageTypes[imageType]
	return ok
}
