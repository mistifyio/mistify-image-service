package metadata_test

import (
	"encoding/json"
	"testing"

	"github.com/coreos/go-etcd/etcd"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/suite"
)

type EtcdTestSuite struct {
	StoreTestSuite
	EtcdConfig *metadata.EtcdConfig
	// EtcdClient is for post-test cleanup of etcd
	EtcdClient *etcd.Client
}

func (s *EtcdTestSuite) SetupSuite() {
	// Etcd specific suite setup
	s.EtcdClient = etcd.NewClient(nil)

	// General store suite setup
	s.StoreTestSuite.SetupSuite()
}

func (s *EtcdTestSuite) SetupTest() {
	// Etcd specific test setup
	s.EtcdConfig = &metadata.EtcdConfig{
		Prefix: "etcdTest-" + uuid.New(),
	}
	configBytes, _ := json.Marshal(s.EtcdConfig)
	s.StoreConfig = configBytes

	// General store test setup
	s.StoreTestSuite.SetupTest()
}

func (s *EtcdTestSuite) TearDownTest() {
	// Clean up etcd
	_, err := s.EtcdClient.Delete(s.EtcdConfig.Prefix, true)
	s.NoError(err)
}

func TestEtcdTestSuite(t *testing.T) {
	s := new(EtcdTestSuite)
	s.StoreName = "etcd"
	suite.Run(t, s)
}

func (s *EtcdTestSuite) TestConfigValidate() {
	tests := []struct {
		description string
		config      *metadata.EtcdConfig
		expectedErr error
	}{
		{"empty config should be valid",
			&metadata.EtcdConfig{}, nil},
		{"filepath-only config should be valid",
			&metadata.EtcdConfig{Filepath: "/foo"}, nil},
		{"prefix-only config should be valid",
			&metadata.EtcdConfig{Prefix: "foo"}, nil},
		{"config to use for tests should be valid",
			s.EtcdConfig, nil},
		{"cert-only config should be invalid",
			&metadata.EtcdConfig{Cert: "foo"}, metadata.ErrIncompleteTLSConfig},
		{"key-only config should be invalid",
			&metadata.EtcdConfig{Key: "bar"}, metadata.ErrIncompleteTLSConfig},
		{"cacert-only config should be invalid",
			&metadata.EtcdConfig{CaCert: "baz"}, metadata.ErrIncompleteTLSConfig},
		{"complete tls config should be valid",
			&metadata.EtcdConfig{
				Cert:   "foo",
				Key:    "bar",
				CaCert: "baz",
			},
			nil,
		},
	}

	for _, test := range tests {
		s.Equal(test.expectedErr, test.config.Validate(), test.description)
	}
}

func (s *EtcdTestSuite) TestInit() {
	tests := []struct {
		description string
		configJSON  string
		expectedErr bool
	}{
		{"bad json should fail",
			"not actually json", true},
		{"incomplete tls config should fail",
			`{"cert":"blah"}`, true},
		{"invalid filepath should fail",
			`{"filepath":"/dev/null/foo"}`, true},
		{"bad tls config should fail",
			`{"cert":"/dev/null/foo", "key":"asdf", "cacert":"asdf"}`, true},
		{"valid config should succeed",
			string(s.StoreConfig), false},
	}

	for _, test := range tests {
		store := metadata.NewStore("etcd")
		config := []byte(test.configJSON)
		err := store.Init(config)
		if test.expectedErr {
			s.Error(err, test.description)
		} else {
			s.NoError(err, test.description)
		}
	}
}
