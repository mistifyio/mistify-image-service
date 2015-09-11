package images_test

import (
	"bytes"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/stretchr/testify/suite"
)

type StoreTestSuite struct {
	suite.Suite
	StoreName   string
	StoreConfig []byte
	Store       images.Store
	ImageID     string
	ImageData   []byte
}

func (s *StoreTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)
	s.ImageID = "foobar"
	s.ImageData = []byte("testdatatestdatatestdata")
}

func (s *StoreTestSuite) SetupTest() {
	s.Store = images.NewStore(s.StoreName)
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

func (s *StoreTestSuite) TestPut() {
	in := bytes.NewReader(s.ImageData)
	s.NoError(s.Store.Put(s.ImageID, in))
}

func (s *StoreTestSuite) TestStat() {
	in := bytes.NewReader(s.ImageData)
	_ = s.Store.Put(s.ImageID, in)

	stat, err := s.Store.Stat(s.ImageID)
	s.NoError(err)
	s.NotNil(stat)
	s.EqualValues(len(s.ImageData), stat.Size())

	stat, err = s.Store.Stat("asdf")
	s.Error(err)
}

func (s *StoreTestSuite) TestGet() {
	in := bytes.NewReader(s.ImageData)
	_ = s.Store.Put(s.ImageID, in)

	out := bytes.NewBuffer(make([]byte, 0, len(s.ImageData)))
	s.NoError(s.Store.Get(s.ImageID, out))
	s.Equal(s.ImageData, out.Bytes())

	s.Error(s.Store.Get("asdf", out))
}

func (s *StoreTestSuite) TestDelete() {
	in := bytes.NewReader(s.ImageData)
	_ = s.Store.Put(s.ImageID, in)

	s.NoError(s.Store.Delete(s.ImageID))
}

func (s *StoreTestSuite) TestShutdown() {
	s.NoError(s.Store.Shutdown())
}
