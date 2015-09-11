package images_test

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/mistifyio/mistify-image-service/images/mocks"
	"github.com/stretchr/testify/suite"
)

type ImagesTestSuite struct {
	suite.Suite
	MockStoreName string
}

func (s *ImagesTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)
	s.MockStoreName = "mock"
}

func TestImagesTestSuite(t *testing.T) {
	suite.Run(t, new(ImagesTestSuite))
}

func (s *ImagesTestSuite) TestList() {
	list := images.List()
	s.NotNil(list)
}

func (s *ImagesTestSuite) TestRegister() {
	images.Register(s.MockStoreName, func() images.Store {
		return &mocks.Store{}
	})

	s.Contains(images.List(), s.MockStoreName, "should contain registered store")
}

func (s *ImagesTestSuite) TestNewStore() {
	images.Register(s.MockStoreName, func() images.Store {
		return &mocks.Store{}
	})

	s.NotNil(images.NewStore(s.MockStoreName), "should create registered store")
	s.Nil(images.NewStore("asdf"), "shouldn't create unregistered store")
}
