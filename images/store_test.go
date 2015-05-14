package images_test

import (
	"io"
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
	log.SetLevel(log.WarnLevel)

	os.Exit(m.Run())
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
