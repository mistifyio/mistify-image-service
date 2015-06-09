package metadata

import (
	"io/ioutil"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MockStore struct{}

func (ms *MockStore) Init(b []byte) error                  { return nil }
func (ms *MockStore) Shutdown() error                      { return nil }
func (ms *MockStore) List(s string) ([]*Image, error)      { return []*Image{}, nil }
func (ms *MockStore) GetByID(s string) (*Image, error)     { return &Image{}, nil }
func (ms *MockStore) GetBySource(s string) (*Image, error) { return &Image{}, nil }
func (ms *MockStore) Put(i *Image) error                   { return nil }
func (ms *MockStore) Delete(s string) error                { return nil }

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()

	log.SetLevel(log.WarnLevel)

	// Store-specific setup

	// KVite
	testKviteFile, err := ioutil.TempFile("", "kvitetest.db")
	if err != nil {
		log.WithField("error", err).Fatal("failed to create kvite temp file")
		return
	}
	testKviteConfig.Filename = testKviteFile.Name()
	defer func() {
		if err := os.RemoveAll(testKviteConfig.Filename); err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"filename": testKviteConfig.Filename,
			}).Error("could not clean up kvite file")
		}
	}()

	// Run the tests
	code = m.Run()
}

func TestList(t *testing.T) {
	list := List()
	assert.NotNil(t, list)
}

func TestNewStore(t *testing.T) {
	list := List()
	if len(list) == 0 {
		return
	}

	assert.NotNil(t, NewStore(list[0]))
	assert.Nil(t, NewStore("qweryasdf"))
}

func TestRegister(t *testing.T) {
	registerMockStore()

	assert.NotNil(t, NewStore("mock"))
}

func registerMockStore() {
	Register("mock", func() Store {
		return &MockStore{}
	})
}
