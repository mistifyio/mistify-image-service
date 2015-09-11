package metadata_test

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	StoreName   string
	StoreConfig []byte
	Store       metadata.Store
	Image       *metadata.Image
}

func (s *StoreTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)
	s.Image = &metadata.Image{
		ID:     metadata.NewID(),
		Type:   "kvm",
		Source: "http://localhost",
	}
}

func (s *StoreTestSuite) SetupTest() {
	s.Store = metadata.NewStore(s.StoreName)
	_ = s.Store.Init(s.StoreConfig)
}

func (s *StoreTestSuite) TestConfigValidate() {
	// This is going to be unique to each store type
	s.Fail("test suite does not define TestConfigValidate", s.StoreName)
}

func (s *StoreTestSuite) TestInit() {
	// This is going to be unique to each store type based on the config
	s.Fail("test suite does not define TestInit", s.StoreName)
}

func (s *StoreTestSuite) Put() {
	s.NoError(s.Store.Put(s.Image), "complete image should be put")
}

func (s *StoreTestSuite) TestGetBySource() {
	_ = s.Store.Put(s.Image)

	// Image exists
	image, err := s.Store.GetBySource(s.Image.Source)
	s.NoError(err, "retrieving existing image should not fail")
	s.NotNil(image, "image should be found")
	s.Equal(s.Image.ID, image.ID, "image should be what we expect")

	// Image doesn't exist
	image, err = s.Store.GetBySource("foobar")
	s.Equal(metadata.ErrNotFound, err, "image shouldn't be found")
}

func (s *StoreTestSuite) TestGetByID() {
	_ = s.Store.Put(s.Image)

	// Image exists
	image, err := s.Store.GetByID(s.Image.ID)
	s.NoError(err, "retrieving existing image should not fail")
	s.NotNil(image, "image should be found")
	s.Equal(s.Image.ID, image.ID, "image should be what we expect")

	// Image doesn't exist
	image, err = s.Store.GetByID("foobar")
	s.Equal(metadata.ErrNotFound, err, "image shouldn't be found")
}

func (s *StoreTestSuite) TestList() {
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

func (s *StoreTestSuite) TestDelete() {
	_ = s.Store.Put(s.Image)

	s.NoError(s.Store.Delete(s.Image.ID), "deleting existing image shouldn't error")
	image, _ := s.Store.GetByID(s.Image.ID)
	s.Nil(image, "image should be deleted")

	s.NoError(s.Store.Delete(s.Image.ID), "deleting missing image shouldn't error")
}

func (s *StoreTestSuite) TestShutdown() {
	s.NoError(s.Store.Shutdown(), "shutdown shouldn't error")
	s.NoError(s.Store.Shutdown(), "second shutdown shouldn't error")
}
