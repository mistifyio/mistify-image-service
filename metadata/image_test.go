package metadata_test

import (
	"errors"
	"testing"
	"time"

	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/mistifyio/mistify-image-service/metadata/mocks"
	"github.com/stretchr/testify/suite"
)

type ImageTestSuite struct {
	suite.Suite
	MockStoreName string
	TestImage     *metadata.Image
	Store         metadata.Store
}

func (s *ImageTestSuite) SetupSuite() {
	s.MockStoreName = "mock"
	metadata.Register(s.MockStoreName, func() metadata.Store {
		return &mocks.Store{}
	})
}

func (s *ImageTestSuite) SetupTest() {
	s.Store = metadata.NewStore(s.MockStoreName)
	s.TestImage = &metadata.Image{
		ID:    metadata.NewID(),
		Store: s.Store,
	}

	s.Store.(*mocks.Store).On("Put", s.TestImage).Return(nil)
}

func TestImageTestSuite(t *testing.T) {
	suite.Run(t, new(ImageTestSuite))
}

func (s *ImageTestSuite) TestNewID() {
	s.NotNil(metadata.NewID())
}

func (s *ImageTestSuite) TestIsValidImageType() {
	for imageType := range metadata.ValidImageTypes {
		s.True(metadata.IsValidImageType(imageType), "should be a valid image type")
	}

	s.False(metadata.IsValidImageType("foobar"), "should be an invalid image type")
}

func (s *ImageTestSuite) TestSetPending() {
	s.NoError(s.TestImage.SetPending())
	s.Equal(metadata.StatusPending, s.TestImage.Status)
}

func (s *ImageTestSuite) TestSetDownloading() {
	size := int64(10000)

	s.NoError(s.TestImage.SetDownloading(size))
	s.Equal(metadata.StatusDownloading, s.TestImage.Status)
	s.WithinDuration(s.TestImage.DownloadStart, time.Now(), 1*time.Second, "downloadstart should be set to now")
	s.Equal(size, s.TestImage.ExpectedSize, "expected size should be set to the size")
}

func (s *ImageTestSuite) TestUpdateSize() {
	size := int64(10000)

	for i := 0; i < 5; i++ {
		newSize := size + int64(i*1000)
		s.NoError(s.TestImage.UpdateSize(newSize))
		s.Equal(newSize, s.TestImage.Size, "image size should be updated")
	}
}

func (s *ImageTestSuite) TestSetFinished() {
	s.NoError(s.TestImage.SetFinished(nil))
	s.Equal(metadata.StatusComplete, s.TestImage.Status, "finishing with no error should set complete status")
	s.WithinDuration(s.TestImage.DownloadEnd, time.Now(), 1*time.Second, "downloadend should be set to now")

	image := &metadata.Image{
		ID:    metadata.NewID(),
		Store: metadata.NewStore(s.MockStoreName),
	}
	image.Store.(*mocks.Store).On("Put", image).Return(nil)
	s.NoError(image.SetFinished(errors.New("An Error")))
	s.Equal(metadata.StatusError, image.Status, "finishing with error should set error status")
	s.WithinDuration(image.DownloadEnd, time.Now(), 1*time.Minute, "downloadend should be set to now")
}
