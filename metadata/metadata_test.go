package metadata_test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/mistifyio/mistify-image-service/metadata/mocks"
	"github.com/stretchr/testify/suite"
)

type MetadataTestSuite struct {
	suite.Suite
	MockStoreName string
}

func (s *MetadataTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)
	s.MockStoreName = "mock"
}

func TestMetadataTestSuite(t *testing.T) {
	suite.Run(t, new(MetadataTestSuite))
}

func (s *MetadataTestSuite) TestList() {
	list := metadata.List()
	s.NotNil(list)
}

func (s *MetadataTestSuite) TestRegister() {
	metadata.Register(s.MockStoreName, func() metadata.Store {
		return &mocks.Store{}
	})

	s.Contains(metadata.List(), s.MockStoreName, "should contain registered store")
}

func (s *MetadataTestSuite) TestNewStore() {
	metadata.Register(s.MockStoreName, func() metadata.Store {
		return &mocks.Store{}
	})

	s.NotNil(metadata.NewStore(s.MockStoreName), "should create registered store")
	s.Nil(metadata.NewStore("asdf"), "shouldn't create unregistered store")
}
