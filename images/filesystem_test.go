package images_test

import (
	"bytes"
	"testing"

	"github.com/mistifyio/mistify-image-service/images"
	"github.com/stretchr/testify/assert"
)

var fsConfig = &images.FSConfig{
	Dir: "/tmp/testimages",
}

var fsStore images.Store

func TestFSStore(t *testing.T) {
	var fsc *images.FSConfig

	fsc = &images.FSConfig{}
	assert.Error(t, fsc.Validate())

	assert.NoError(t, fsConfig.Validate())
}

func TestFSInit(t *testing.T) {
	fs := images.NewStore("fs")
	assert.NoError(t, fs.Init(fsConfig))

	fsStore = fs
}

func TestFSPut(t *testing.T) {
	in := bytes.NewReader(mockImageData)
	assert.NoError(t, fsStore.Put(mockImageID, in))
}

func TestFSGet(t *testing.T) {
	out := bytes.NewBuffer(make([]byte, 0, len(mockImageData)))
	assert.NoError(t, fsStore.Get(mockImageID, out))
	assert.Equal(t, string(mockImageData), out.String())
}

func TestFSStat(t *testing.T) {
	stat, err := fsStore.Stat(mockImageID)
	assert.NoError(t, err)
	assert.NotNil(t, stat)
	assert.EqualValues(t, len(mockImageData), stat.Size())
}

func TestFSDelete(t *testing.T) {
	assert.NoError(t, fsStore.Delete(mockImageID))
}

func TestFSShutdown(t *testing.T) {
	assert.NoError(t, fsStore.Shutdown())
}
