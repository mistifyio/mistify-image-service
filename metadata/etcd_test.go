package metadata_test

import (
	"encoding/json"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-etcd/etcd"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/suite"
)

type EtcdTestSuite struct {
	suite.Suite
	EtcdConfig *metadata.EtcdConfig
	// EtcdClient is for post-test cleanup of etcd
	EtcdClient *etcd.Client
	Store      metadata.Store
	Image      *metadata.Image
}

func (s *EtcdTestSuite) SetupSuite() {
	s.Image = &metadata.Image{
		ID:     metadata.NewID(),
		Type:   "kvm",
		Source: "http://localhost",
	}

	s.EtcdClient = etcd.NewClient(nil)
}

func (s *EtcdTestSuite) SetupTest() {
	// Configure and initialize a new etcd Store
	s.EtcdConfig = &metadata.EtcdConfig{
		Prefix: "etcdTest-" + uuid.New(),
	}
	s.Store = metadata.NewStore("etcd")
	configBytes, _ := json.Marshal(s.EtcdConfig)
	_ = s.Store.Init(configBytes)
}

func (s *EtcdTestSuite) TearDownTest() {
	// Clean up etcd
	_, err := s.EtcdClient.Delete(s.EtcdConfig.Prefix, true)
	s.NoError(err)
}

func TestEtcdTestSuite(t *testing.T) {
	suite.Run(t, new(EtcdTestSuite))
}

func (s *EtcdTestSuite) TestEtcdConfigValidate() {
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

func (s *EtcdTestSuite) TestEtcdInit() {
	goodConfigBytes, _ := json.Marshal(s.EtcdConfig)

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
			string(goodConfigBytes), false},
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

func (s *EtcdTestSuite) TestEtcdPut() {
	s.NoError(s.Store.Put(s.Image), "complete image should be put")
}

func (s *EtcdTestSuite) TestEtcdGetBySource() {
	_ = s.Store.Put(s.Image)

	// Image exists
	image, err := s.Store.GetBySource(s.Image.Source)
	s.NoError(err, "retrieving existing image should not fail")
	s.NotNil(image, "image should be found")
	s.Equal(s.Image.ID, image.ID, "image should be what we expect")

	// Image doesn't exist
	image, err = s.Store.GetBySource("foobar")
	s.Error(err, "image shouldn't be found")
}

func (s *EtcdTestSuite) TestEtcdGetByID() {
	_ = s.Store.Put(s.Image)

	// Image exists
	image, err := s.Store.GetByID(s.Image.ID)
	s.NoError(err, "retrieving existing image should not fail")
	s.NotNil(image, "image should be found")
	s.Equal(s.Image.ID, image.ID, "image should be what we expect")

	// Image doesn't exist
	image, err = s.Store.GetByID("foobar")
	s.Error(err, "image shouldn't be found")
}

func (s *EtcdTestSuite) TestEtcdList() {
	_ = s.Store.Put(s.Image)

	images, err := s.Store.List("")
	s.NoError(err, "listing all images shouldn't error")
	s.NotNil(images)
	s.Len(images, 1, "list should only contain the one image added")

	var found bool
	for _, image := range images {
		if image.ID == s.Image.ID {
			found = true
			break
		}
	}
	s.True(found, "image should be in list")
}

func (s *EtcdTestSuite) TestEtcdDelete() {
	_ = s.Store.Put(s.Image)

	s.NoError(s.Store.Delete(s.Image.ID), "deleting existing image shouldn't error")
	image, _ := s.Store.GetByID(s.Image.ID)
	s.Nil(image, "image should be deleted")

	s.Error(s.Store.Delete(s.Image.ID), "deleting missing image should error")
}

func (s *EtcdTestSuite) TestEtcdShutdown() {
	s.NoError(s.Store.Shutdown(), "shutdown shouldn't error")
}

func init() {
	log.SetLevel(log.FatalLevel)
}
