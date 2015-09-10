package metadata

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testEtcdConfig = &etcdConfig{
	Prefix: "etcdTest",
}

var testEtcdStore Store
var testEtcdImage *Image

func TestEtcdConfigValidate(t *testing.T) {
	var ec *etcdConfig

	ec = &etcdConfig{}
	assert.NoError(t, ec.Validate())

	ec = &etcdConfig{
		Filepath: "/foo",
	}
	assert.NoError(t, ec.Validate())

	ec = &etcdConfig{
		Cert: "foobar",
	}
	assert.Error(t, ec.Validate())

	ec = &etcdConfig{
		Cert:   "foobar",
		Key:    "foobar",
		CaCert: "foobar",
	}
	assert.NoError(t, ec.Validate())

	assert.NoError(t, testEtcdConfig.Validate())
}

func TestEtcdInit(t *testing.T) {
	ec := NewStore("etcd")
	assert.Error(t, ec.Init([]byte("not actually json")))

	ec = NewStore("etcd")
	assert.Error(t, ec.Init([]byte(`{"cert":"blah"}`)))

	ec = NewStore("etcd")
	assert.Error(t, ec.Init([]byte(`{"filepath":"/dev/null/foo"}`)))

	ec = NewStore("etcd")
	assert.Error(t, ec.Init([]byte(`{"cert":"/dev/null/foo", "key":"asdf", "cacert":"asdf"}`)))

	ec = NewStore("etcd")
	configBytes, _ := json.Marshal(testEtcdConfig)
	assert.NoError(t, ec.Init(configBytes))

	testEtcdStore = ec

	ec = NewStore("etcd")
	assert.NoError(t, ec.Init(configBytes))
}

func TestEtcdPut(t *testing.T) {
	assert.NoError(t, testEtcdStore.Put(testEtcdImage))
}

func TestEtcdGetBySource(t *testing.T) {
	image, err := testEtcdStore.GetBySource(testEtcdImage.Source)
	assert.NoError(t, err)
	assert.Equal(t, testEtcdImage.ID, image.ID)

	image, err = testEtcdStore.GetBySource("foobar")
	assert.Nil(t, image)
}

func TestEtcdGetByID(t *testing.T) {
	image, err := testEtcdStore.GetByID(testEtcdImage.ID)
	assert.NoError(t, err)
	assert.Equal(t, testEtcdImage.ID, image.ID)

	image, err = testEtcdStore.GetByID("foobar")
	assert.Nil(t, image)
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
	assert.Error(t, testEtcdStore.Delete(testEtcdImage.ID))
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
