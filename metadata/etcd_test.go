package metadata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testEtcdConfig = &EtcdConfig{
	Prefix: "etcdTest",
}

var testEtcdStore Store
var testEtcdImage *Image

func TestEtcdConfigValidate(t *testing.T) {
	var ec *EtcdConfig

	ec = &EtcdConfig{}
	assert.NoError(t, ec.Validate())

	ec = &EtcdConfig{
		Filepath: "/foo",
	}
	assert.NoError(t, ec.Validate())

	ec = &EtcdConfig{
		Cert: "foobar",
	}
	assert.Error(t, ec.Validate())

	ec = &EtcdConfig{
		Cert:   "foobar",
		Key:    "foobar",
		CaCert: "foobar",
	}
	assert.NoError(t, ec.Validate())

	assert.NoError(t, testEtcdConfig.Validate())
}

func TestEtcdInit(t *testing.T) {
	ec := NewStore("etcd")
	configBytes, _ := json.Marshal(testEtcdConfig)
	assert.NoError(t, ec.Init(configBytes))

	testEtcdStore = ec
}

func TestEtcdPut(t *testing.T) {
	assert.NoError(t, testEtcdStore.Put(testEtcdImage))
}

func TestEtcdGetBySource(t *testing.T) {
	image, err := testEtcdStore.GetBySource(testEtcdImage.Source)
	assert.NoError(t, err)
	assert.Equal(t, testEtcdImage.ID, image.ID)
}

func TestEtcdGetByID(t *testing.T) {
	image, err := testEtcdStore.GetByID(testEtcdImage.ID)
	assert.NoError(t, err)
	assert.Equal(t, testEtcdImage.ID, image.ID)
}

func TestEtcdList(t *testing.T) {
	images, err := testEtcdStore.List("")
	assert.NoError(t, err)
	var found bool
	for _, image := range images {
		if image.ID == testEtcdImage.ID {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestEtcdDelete(t *testing.T) {
	assert.NoError(t, testEtcdStore.Delete(testEtcdImage.ID))
}

func TestEtcdShutdown(t *testing.T) {
	assert.NoError(t, testEtcdStore.Shutdown())
}

func init() {
	testEtcdImage = &Image{
		ID:     NewID(),
		Type:   "kvm",
		Source: "http://localhost",
	}
}
