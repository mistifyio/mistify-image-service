package metadata_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/suite"
)

type KViteTestSuite struct {
	StoreTestSuite
	KViteConfig *metadata.KViteConfig
}

func (s *KViteTestSuite) SetupSuite() {
	// General store suite setup
	s.StoreTestSuite.SetupSuite()
}

func (s *KViteTestSuite) SetupTest() {
	// KVite specific test setup
	dbfile, _ := ioutil.TempFile("", "kviteTest-"+uuid.New()+".db")
	_ = dbfile.Close()

	s.KViteConfig = &metadata.KViteConfig{
		Filename: dbfile.Name(),
		Table:    "test",
	}
	configBytes, _ := json.Marshal(s.KViteConfig)
	s.StoreConfig = configBytes

	// General store test setup
	s.StoreTestSuite.SetupTest()
}

func (s *KViteTestSuite) TearDownTest() {
	// Clean up kvite file
	s.NoError(os.Remove(s.KViteConfig.Filename))
}

func TestKViteTestSuite(t *testing.T) {
	s := new(KViteTestSuite)
	s.StoreName = "kvite"
	suite.Run(t, s)
}

func (s *KViteTestSuite) TestConfigValidate() {
	tests := []struct {
		description string
		config      *metadata.KViteConfig
		expectedErr error
	}{
		{"empty config should be invalid",
			&metadata.KViteConfig{}, metadata.ErrMissingFilename},
		{"filename-only config should be invalid",
			&metadata.KViteConfig{Filename: "/foo"}, metadata.ErrMissingTable},
		{"table-only config should be valid",
			&metadata.KViteConfig{Table: "foo"}, metadata.ErrMissingFilename},
		{"config to use for tests should be valid",
			s.KViteConfig, nil},
	}

	for _, test := range tests {
		s.Equal(test.expectedErr, test.config.Validate(), test.description)
	}
}

func (s *KViteTestSuite) TestInit() {
	tests := []struct {
		description string
		configJSON  string
		expectedErr bool
	}{
		{"bad json should fail",
			"not actually json", true},
		{"incomplete config should fail",
			`{"table":"blah"}`, true},
		{"invalid filepath should fail",
			`{"filepath":"/dev/null/foo","table":"foo"}`, true},
		{"valid config should succeed",
			string(s.StoreConfig), false},
	}

	for _, test := range tests {
		store := metadata.NewStore("kvite")
		config := []byte(test.configJSON)
		err := store.Init(config)
		if test.expectedErr {
			s.Error(err, test.description)
		} else {
			s.NoError(err, test.description)
		}
	}
}
