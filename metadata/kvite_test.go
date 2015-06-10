package metadata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testKviteConfig = &KViteConfig{
	Table: "kvitetest",
}

var testKviteStore Store
var testKviteImage *Image

func TestKViteConfigValidate(t *testing.T) {
	var kvc *KViteConfig

	kvc = &KViteConfig{}
	assert.Error(t, kvc.Validate())

	kvc = &KViteConfig{
		Filename: "/foo",
	}
	assert.Error(t, kvc.Validate())

	kvc = &KViteConfig{
		Table: "foobar",
	}
	assert.Error(t, kvc.Validate())

	assert.NoError(t, testKviteConfig.Validate())
}

func TestKViteInit(t *testing.T) {
	kv := NewStore("kvite")
	configBytes, _ := json.Marshal(testKviteConfig)
	assert.NoError(t, kv.Init(configBytes))

	testKviteStore = kv
}

func TestKVitePut(t *testing.T) {
	assert.NoError(t, testKviteStore.Put(testKviteImage))
}

func TestKViteGetByID(t *testing.T) {
	image, err := testKviteStore.GetByID(testKviteImage.ID)
	assert.NoError(t, err)
	assert.Equal(t, testKviteImage.ID, image.ID)
}

func TestKViteGetBySource(t *testing.T) {
	image, err := testKviteStore.GetBySource(testKviteImage.Source)
	assert.NoError(t, err)
	assert.Equal(t, testKviteImage.ID, image.ID)
}

func TestKViteList(t *testing.T) {
	images, err := testKviteStore.List("")
	assert.NoError(t, err)
	var found bool
	for _, image := range images {
		if image.ID == testKviteImage.ID {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestKViteDelete(t *testing.T) {
	assert.NoError(t, testKviteStore.Delete(testKviteImage.ID))
}

func TestKViteShutdown(t *testing.T) {
	assert.NoError(t, testKviteStore.Shutdown())
}

func init() {
	testKviteImage = &Image{
		ID:     NewID(),
		Type:   "kvm",
		Source: "http://localhost",
	}
}
