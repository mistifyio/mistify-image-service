package metadata_test

import (
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
	log.SetLevel(log.WarnLevel)

	os.Exit(m.Run())
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
