package metadata_test

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/stretchr/testify/assert"
)

type MockStore struct{}

func (ms *MockStore) Init(b []byte) error                           { return nil }
func (ms *MockStore) Shutdown() error                               { return nil }
func (ms *MockStore) List(s string) ([]*metadata.Image, error)      { return []*metadata.Image{}, nil }
func (ms *MockStore) GetByID(s string) (*metadata.Image, error)     { return &metadata.Image{}, nil }
func (ms *MockStore) GetBySource(s string) (*metadata.Image, error) { return &metadata.Image{}, nil }
func (ms *MockStore) Put(i *metadata.Image) error                   { return nil }
func (ms *MockStore) Delete(s string) error                         { return nil }

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log.SetLevel(log.WarnLevel)

	// Store-specific setup

	// KVite
	kviteFile, err := ioutil.TempFile("", "kvitetest.db")
	if err != nil {
		log.WithField("error", err).Fatal("failed to create kvite temp file")
		return
	}
	kviteConfig.Filename = kviteFile.Name()
	defer func() {
		if err := os.RemoveAll(kviteConfig.Filename); err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"filename": kviteConfig.Filename,
			}).Error("could not clean up kvite file")
		}
	}()

	// Run the tests
	code = m.Run()
}

func TestList(t *testing.T) {
	list := metadata.List()
	assert.NotNil(t, list)
}

func TestNewStore(t *testing.T) {
	list := metadata.List()
	if len(list) == 0 {
		return
	}

	assert.NotNil(t, metadata.NewStore(list[0]))
	assert.Nil(t, metadata.NewStore("qweryasdf"))
}

func TestRegister(t *testing.T) {
	registerMockStore()

	assert.NotNil(t, metadata.NewStore("mock"))
}

func registerMockStore() {
	metadata.Register("mock", func() metadata.Store {
		return &MockStore{}
	})
}
