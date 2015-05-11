package imageservice_test

import (
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var ctx *imageservice.Context

func TestMain(m *testing.M) {
	log.SetLevel(log.WarnLevel)

	imgDir := "/tmp/test_image_service"
	mdFile := "/tmp/test_image_service.db"

	// Clean up
	if err := cleanup(imgDir, mdFile); err != nil {
		log.Fatal("failed to clean up")
	}

	// Set up config
	viper.SetDefault("imageStoreType", "fs")
	viper.SetDefault("imageStoreConfig", &images.FSConfig{
		Dir: imgDir,
	})

	viper.SetDefault("metadataStoreType", "kvite")
	viper.SetDefault("metadataStoreConfig", &metadata.KViteConfig{
		Filename: mdFile,
		Table:    "test_image_service",
	})

	// Run the tests
	code := m.Run()

	// Clean up
	if err := cleanup(imgDir, mdFile); err != nil {
		log.Fatal("failed to clean up")
	}

	// Exit with test run exit code
	os.Exit(code)
}

func cleanup(imgDir, mdFile string) error {
	if err := os.RemoveAll(imgDir); err != nil && !os.IsNotExist(err) {
		log.Error(err)
		return err
	}

	if err := os.Remove(mdFile); err != nil && !os.IsNotExist(err) {
		log.Error(err)
		return err
	}

	return nil
}

func TestContextNewImageStore(t *testing.T) {
	ctx := &imageservice.Context{}

	assert.Error(t, ctx.NewImageStore("asdfwqfas"))
	assert.Nil(t, ctx.ImageStore)

	assert.NoError(t, ctx.NewImageStore("fs"))
	assert.NotNil(t, ctx.ImageStore)
}

func TestContextNewMetadataStore(t *testing.T) {
	ctx := &imageservice.Context{}

	assert.Error(t, ctx.NewMetadataStore("asdfwqfas"))
	assert.Nil(t, ctx.MetadataStore)

	assert.NoError(t, ctx.NewMetadataStore("kvite"))
	assert.NotNil(t, ctx.MetadataStore)
}

func TestContextNewFetcher(t *testing.T) {
	ctx := &imageservice.Context{}
	ctx.NewFetcher()

	assert.NotNil(t, ctx.Fetcher)
}

func TestNewContext(t *testing.T) {
	context, err := imageservice.NewContext()
	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.NotNil(t, context.ImageStore)
	assert.NotNil(t, context.MetadataStore)
	assert.NotNil(t, context.Fetcher)
	ctx = context
}
