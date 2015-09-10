package metadata

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	id := NewID()
	assert.NotEmpty(t, id)
}

func TestIsValidImageType(t *testing.T) {
	for imageType := range ValidImageTypes {
		assert.True(t, IsValidImageType(imageType))
	}
	assert.False(t, IsValidImageType("foobar"))
}

func newTestImage() *Image {
	registerMockStore()
	image := &Image{
		ID:    NewID(),
		Store: NewStore("mock"),
	}
	return image
}

func TestSetPending(t *testing.T) {
	image := newTestImage()

	assert.NoError(t, image.SetPending())
	assert.Equal(t, StatusPending, image.Status)
}

func TestSetDownloading(t *testing.T) {
	image := newTestImage()
	size := int64(10000)

	assert.NoError(t, image.SetDownloading(size))
	assert.Equal(t, StatusDownloading, image.Status)
	assert.WithinDuration(t, image.DownloadStart, time.Now(), 1*time.Minute)
	assert.Equal(t, size, image.ExpectedSize)
}

func TestUpdateSize(t *testing.T) {
	image := newTestImage()
	size := int64(10000)

	for i := 0; i < 5; i++ {
		newSize := size + int64(i*1000)
		assert.NoError(t, image.UpdateSize(newSize))
		assert.Equal(t, newSize, image.Size)
	}
}

func TestSetFinished(t *testing.T) {
	image := newTestImage()

	assert.Nil(t, image.SetFinished(nil))
	assert.Equal(t, StatusComplete, image.Status)
	assert.WithinDuration(t, image.DownloadEnd, time.Now(), 1*time.Minute)

	image = newTestImage()
	imgErr := errors.New("An Error")
	assert.Nil(t, image.SetFinished(imgErr))
	assert.Equal(t, StatusError, image.Status)
	assert.WithinDuration(t, image.DownloadEnd, time.Now(), 1*time.Minute)
}
