package imageservice_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/stretchr/testify/assert"
)

var mockImageData = []byte("testdatatestdatatestdata")

func TestFetcherReceive(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost", bytes.NewReader(mockImageData))
	assert.NoError(t, err)
	req.Header.Add("X-Image-Type", "kvm")
	req.Header.Add("X-Image-Comment", "test image")

	image, err := ctx.Fetcher.Receive(req)
	assert.NoError(t, err)
	assert.NotNil(t, image)
	assert.EqualValues(t, len(mockImageData), image.Size)
	assert.Equal(t, "kvm", image.Type)
	assert.Equal(t, "test image", image.Comment)
}

func TestFetcherFetch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(mockImageData)
	}))
	defer ts.Close()

	image := &metadata.Image{
		Source: ts.URL,
		Type:   "kvm",
	}

	image, err := ctx.Fetcher.Fetch(image)
	assert.NoError(t, err)
	assert.NotNil(t, image.ID)
	for i := 0; i < 10; i++ {
		image, err := ctx.MetadataStore.GetByID(image.ID)
		assert.NoError(t, err)
		if image.Status == metadata.StatusComplete || image.Status == metadata.StatusError {
			break
		}
		time.Sleep(1 * time.Second)
	}
	assert.Equal(t, metadata.StatusComplete, image.Status)
	assert.Equal(t, "kvm", image.Type)
	assert.Equal(t, ts.URL, image.Source)
	assert.EqualValues(t, len(mockImageData), image.Size)
}
