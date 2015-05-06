package imageservice

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/metadata"
)

type (
	// Fetcher handles fetching new images and updating metadata accordingly
	Fetcher struct {
		ctx *Context
	}
)

// NewFetcher creates a new Fetcher
func NewFetcher(ctx *Context) *Fetcher {
	fetcher := &Fetcher{
		ctx: ctx,
	}
	return fetcher
}

// Fetch runs pre-flight checks and kicks off an asynchronous image download
func (fetcher *Fetcher) Fetch(image *metadata.Image) (*metadata.Image, error) {
	// Ensure sufficient information for fetching
	if image.Source == "" {
		return nil, errors.New("missing image source")
	}
	if image.Type == "" {
		return nil, errors.New("missing image type")
	}

	// Avoid re-downloading the same image. If a redownload is desired, first
	// delete the existing image.
	existingImage, err := fetcher.ctx.MetadataStore.GetBySource(image.Source)
	if existingImage != nil || err != nil {
		return existingImage, err
	}

	// Additional metadata preparation and initial save
	image.NewID()
	image.Store = fetcher.ctx.MetadataStore
	if err := image.SetPending(); err != nil {
		return nil, err
	}

	// Kick off the download
	go fetcher.fetchImage(image)

	return image, nil
}

// fetchImage downloads a remote image
func (fetcher *Fetcher) fetchImage(image *metadata.Image) {
	var err error
	monitorStop := make(chan struct{}, 1)

	defer func() {
		// Stop size monitoring
		monitorStop <- struct{}{}
		// Last size update
		_ = fetcher.updateImageSize(image)
		// Set final status
		_ = image.SetFinished(err)
	}()

	// Start the download
	resp, err := http.Get(image.Source)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("unexpected response status")
		log.WithFields(log.Fields{
			"error":        err,
			"expectedCode": http.StatusOK,
			"statusCode":   resp.StatusCode,
			"image":        image,
		}).Error(err)
		return
	}

	// Update status to indicate download has begun
	if err = image.SetDownloading(resp.ContentLength); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"image": image,
		}).Error("failed to SetDownloading")
		return
	}

	// Start watching the size
	go fetcher.monitorDownload(image, monitorStop)

	// Stream the download
	if err = fetcher.ctx.ImageStore.Put(image.ID, resp.Body); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"image": image,
		}).Error("failed to download")
		return
	}
}

// monitorDownload periodically get the file stats from the image store and
// updates the size in the metadata.
func (fetcher *Fetcher) monitorDownload(image *metadata.Image, stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			// Periodic size update
			_ = fetcher.updateImageSize(image)
			time.Sleep(5 * time.Second)
		}
	}
}

// updateImageSize updates the image size in metadata
func (fetcher *Fetcher) updateImageSize(image *metadata.Image) error {
	stat, err := fetcher.ctx.ImageStore.Stat(image.ID)
	if err != nil {
		return err
	}
	return image.UpdateSize(stat.Size())
}
