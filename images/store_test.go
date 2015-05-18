package images_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/stretchr/testify/assert"
)

type MockStore struct{}

func (ms *MockStore) Init(b []byte) error                 { return nil }
func (ms *MockStore) Shutdown() error                     { return nil }
func (ms *MockStore) Stat(id string) (os.FileInfo, error) { return *new(os.FileInfo), nil }
func (ms *MockStore) Get(id string, out io.Writer) error  { return nil }
func (ms *MockStore) Put(id string, in io.Reader) error   { return nil }
func (ms *MockStore) Delete(id string) error              { return nil }

var mockImageID = "foobar"
var mockImageData = []byte("testdatatestdatatestdata")

func TestMain(m *testing.M) {
	// So that defers run when we exit
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log.SetLevel(log.WarnLevel)

	// Store-specific setup
	var err error

	// Filesystem
	if fsConfig.Dir, err = ioutil.TempDir("", "testimages"); err != nil {
		log.WithField("error", err).Fatal("failed to create filesystem temp dir")
		return
	}
	defer func() {
		if err := os.RemoveAll(fsConfig.Dir); err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"dir":   fsConfig.Dir,
			}).Error("could not clean up filesystem")
		}
	}()

	// Run the tests
	code = m.Run()
}

func TestList(t *testing.T) {
	list := images.List()
	assert.NotNil(t, list)
}

func TestNewStore(t *testing.T) {
	list := images.List()
	if len(list) == 0 {
		return
	}

	assert.NotNil(t, images.NewStore(list[0]))
	assert.Nil(t, images.NewStore("qwertyasdf"))
}

func TestRegister(t *testing.T) {
	registerMockStore()

	assert.NotNil(t, images.NewStore("mock"))
}

func registerMockStore() {
	images.Register("mock", func() images.Store {
		return &MockStore{}
	})
}
