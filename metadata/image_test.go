package metadata_test

import (
	"testing"
	"time"

	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	id := metadata.NewID()
	assert.NotEmpty(t, id)
}

func TestIsValidImageType(t *testing.T) {
	for imageType := range metadata.ValidImageTypes {
		assert.True(t, metadata.IsValidImageType(imageType))
	}
	assert.False(t, metadata.IsValidImageType("foobar"))
}

func newImage() *metadata.Image {
	registerMockStore()
	image := &metadata.Image{
		ID:    metadata.NewID(),
		Store: metadata.NewStore("mock"),
	}
	return image
}

func TestSetPending(t *testing.T) {
	image := newImage()

	assert.NoError(t, image.SetPending())
	assert.Equal(t, metadata.StatusPending, image.Status)
}

func TestSetDownloading(t *testing.T) {
	image := newImage()
	size := int64(10000)

	assert.NoError(t, image.SetDownloading(size))
	assert.Equal(t, metadata.StatusDownloading, image.Status)
	assert.WithinDuration(t, image.DownloadStart, time.Now(), 1*time.Minute)
	assert.Equal(t, size, image.ExpectedSize)
}

func TestUpdateSize(t *testing.T) {
	image := newImage()
	size := int64(10000)

	for i := 0; i < 5; i++ {
		newSize := size + int64(i*1000)
		assert.NoError(t, image.UpdateSize(newSize))
		assert.Equal(t, newSize, image.Size)
	}
}

func TestSetFinished(t *testing.T) {
	image := newImage()

	assert.Nil(t, image.SetFinished(nil))
	assert.Equal(t, metadata.StatusComplete, image.Status)
	assert.WithinDuration(t, image.DownloadEnd, time.Now(), 1*time.Minute)
}
