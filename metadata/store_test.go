package metadata_test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/mistifyio/mistify-image-service/metadata/mocks"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	MockStoreName string
}

func (s *StoreTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)
	s.MockStoreName = "mock"
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (s *StoreTestSuite) TestList() {
	list := metadata.List()
	s.NotNil(list)
}

func (s *StoreTestSuite) TestRegister() {
	metadata.Register(s.MockStoreName, func() metadata.Store {
		return &mocks.Store{}
	})

	s.Contains(metadata.List(), s.MockStoreName)
}

func (s *StoreTestSuite) TestNewStore() {
	metadata.Register(s.MockStoreName, func() metadata.Store {
		return &mocks.Store{}
	})

	s.NotNil(metadata.NewStore(s.MockStoreName))
}
