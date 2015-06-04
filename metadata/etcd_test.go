package metadata_test

import (
	"encoding/json"
	"testing"

	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/stretchr/testify/assert"
)

var etcdConfig = &metadata.EtcdConfig{
	Prefix: "etcdTest",
}

var etcdStore metadata.Store
var etcdImage *metadata.Image

func TestEtcdConfigValidate(t *testing.T) {
	var ec *metadata.EtcdConfig

	ec = &metadata.EtcdConfig{}
	assert.NoError(t, ec.Validate())

	ec = &metadata.EtcdConfig{
		Filepath: "/foo",
	}
	assert.NoError(t, ec.Validate())

	ec = &metadata.EtcdConfig{
		Cert: "foobar",
	}
	assert.Error(t, ec.Validate())

	ec = &metadata.EtcdConfig{
		Cert:   "foobar",
		Key:    "foobar",
		CaCert: "foobar",
	}
	assert.NoError(t, ec.Validate())

	assert.NoError(t, etcdConfig.Validate())
}

func TestEtcdInit(t *testing.T) {
	ec := metadata.NewStore("etcd")
	configBytes, _ := json.Marshal(etcdConfig)
	assert.NoError(t, ec.Init(configBytes))

	etcdStore = ec
}

func TestEtcdPut(t *testing.T) {
	assert.NoError(t, etcdStore.Put(etcdImage))
}

func TestEtcdGetBySource(t *testing.T) {
	image, err := etcdStore.GetBySource(etcdImage.Source)
	assert.NoError(t, err)
	assert.Equal(t, etcdImage.ID, image.ID)
}

func TestEtcdGetByID(t *testing.T) {
	image, err := etcdStore.GetByID(etcdImage.ID)
	assert.NoError(t, err)
	assert.Equal(t, etcdImage.ID, image.ID)
}

func TestEtcdList(t *testing.T) {
	images, err := etcdStore.List("")
	assert.NoError(t, err)
	var found bool
	for _, image := range images {
		if image.ID == etcdImage.ID {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestEtcdDelete(t *testing.T) {
	assert.NoError(t, etcdStore.Delete(etcdImage.ID))
}

func TestEtcdShutdown(t *testing.T) {
	assert.NoError(t, etcdStore.Shutdown())
}

func init() {
	etcdImage = &metadata.Image{
		ID:     metadata.NewID(),
		Type:   "kvm",
		Source: "http://localhost",
	}
}
