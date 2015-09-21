package images_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mistifyio/mistify-image-service/images"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/suite"
)

type FSTestSuite struct {
	StoreTestSuite
	FSConfig *images.FSConfig
}

func (s *FSTestSuite) SetupTest() {
	// FS specific test setup
	dir, _ := ioutil.TempDir("", "fsTest-"+uuid.New())
	s.FSConfig = &images.FSConfig{
		Dir: dir,
	}
	s.StoreConfig, _ = json.Marshal(s.FSConfig)

	// General store test setup
	s.StoreTestSuite.SetupTest()
}

func (s *FSTestSuite) TearDownTest() {
	s.NoError(os.RemoveAll(s.FSConfig.Dir))
}

func TestFSTestSuite(t *testing.T) {
	s := new(FSTestSuite)
	s.StoreName = "fs"
	suite.Run(t, s)
}

func (s *FSTestSuite) TestConfigValidate() {
	tests := []struct {
		description string
		config      *images.FSConfig
		expectedErr error
	}{
		{"empty config should be invalid",
			&images.FSConfig{}, images.ErrMissingDir},
		{"config to use for tests should be valid",
			s.FSConfig, nil},
	}

	for _, test := range tests {
		s.Equal(test.expectedErr, test.config.Validate(), test.description)
	}
}

func (s *FSTestSuite) TestInit() {
	tests := []struct {
		description string
		configJSON  string
		expectedErr bool
	}{
		{"bad json should fail",
			"not actually json", true},
		{"invalid config should fail",
			`{}`, true},
		{"bad dir should fail",
			`{"dir":"/dev/null"}`, true},
		{"config to use for tests should succeed",
			string(s.StoreConfig), false},
	}

	for _, test := range tests {
		store := images.NewStore("fs")
		config := []byte(test.configJSON)
		err := store.Init(config)
		if test.expectedErr {
			s.Error(err, test.description)
		} else {
			s.NoError(err, test.description)
		}
	}
}

func (s *FSTestSuite) TestDelete() {
	// General Store.Delete tests
	s.StoreTestSuite.TestDelete()

	// FS specific tests

	// Make sure delete isn't able to delete the directory
	_ = s.Store.Delete("")
	_, err := os.Stat(s.FSConfig.Dir)
	s.NoError(err, "should not delete base directory")
}
