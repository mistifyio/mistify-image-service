package metadata_test

import (
	"testing"

	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/stretchr/testify/assert"
)

var kviteConfig = &metadata.KViteConfig{
	Filename: "/tmp/kvitetest.db",
	Table:    "kvitetest",
}

var kviteStore metadata.Store
var kviteImage *metadata.Image

func TestKViteConfigValidate(t *testing.T) {
	var kvc *metadata.KViteConfig

	kvc = &metadata.KViteConfig{}
	assert.Error(t, kvc.Validate())

	kvc = &metadata.KViteConfig{
		Filename: "/foo",
	}
	assert.Error(t, kvc.Validate())

	kvc = &metadata.KViteConfig{
		Table: "foobar",
	}
	assert.Error(t, kvc.Validate())

	assert.NoError(t, kviteConfig.Validate())
}

func TestKViteInit(t *testing.T) {
	kv := metadata.NewStore("kvite")
	assert.NoError(t, kv.Init(kviteConfig))

	kviteStore = kv
}

func TestKVitePut(t *testing.T) {
	assert.NoError(t, kviteStore.Put(kviteImage))
}

func TestKViteGetByID(t *testing.T) {
	image, err := kviteStore.GetByID(kviteImage.ID)
	assert.NoError(t, err)
	assert.Equal(t, kviteImage.ID, image.ID)
}

func TestKViteGetBySource(t *testing.T) {
	image, err := kviteStore.GetBySource(kviteImage.Source)
	assert.NoError(t, err)
	assert.Equal(t, kviteImage.ID, image.ID)
}

func TestKViteList(t *testing.T) {
	images, err := kviteStore.List("")
	assert.NoError(t, err)
	var found bool
	for _, image := range images {
		if image.ID == kviteImage.ID {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestKViteDelete(t *testing.T) {
	assert.NoError(t, kviteStore.Delete(kviteImage.ID))
}

func TestKViteShutdown(t *testing.T) {
	assert.NoError(t, kviteStore.Shutdown())
}

func init() {
	kviteImage = &metadata.Image{
		Type:   "kvm",
		Source: "http://localhost",
	}
	kviteImage.NewID()
}
