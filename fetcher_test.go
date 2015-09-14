package imageservice_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service"
	"github.com/mistifyio/mistify-image-service/images"
	"github.com/mistifyio/mistify-image-service/metadata"
	"github.com/pborman/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type FetcherTestSuite struct {
	suite.Suite
	Context     *imageservice.Context
	ImageData   []byte
	FetchServer *httptest.Server
	StoreDir    string
}

func (s *FetcherTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)

	s.ImageData = []byte("testdatatestdatatestdata")

	// Test Server to serve image data for fetching
	s.FetchServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/404" {
			http.NotFound(w, r)
			return
		}
		if _, err := w.Write(s.ImageData); err != nil {
			log.WithField("error", err).Error("Failed to write mock image data to response")
		}
	}))
}

func (s *FetcherTestSuite) SetupTest() {
	// NOTE: Using the mocks here would require somewhat complicated logic,
	// approaching a real in-memory Store. Might as well use actual stores
	// (which are tested themselves in their respective packages)

	s.StoreDir, _ = ioutil.TempDir("", "fetcherTest-"+uuid.New())
	// Images Store Setup
	imageStoreConfig := &images.FSConfig{
		Dir: s.StoreDir,
	}
	viper.Set("imageStoreType", "fs")
	viper.Set("imageStoreConfig", imageStoreConfig)

	// Metadata Store Setup
	metadataStoreConfig := &metadata.KViteConfig{
		Filename: filepath.Join(s.StoreDir, "kvite.db"),
		Table:    "test",
	}
	viper.Set("metadataStoreType", "kvite")
	viper.Set("metadataStoreConfig", metadataStoreConfig)

	// Set up context
	ctx, err := imageservice.NewContext()
	s.Require().NoError(err)
	s.Context = ctx
}

func (s *FetcherTestSuite) TearDownTest() {
	s.NoError(os.RemoveAll(s.StoreDir))
}

func (s *FetcherTestSuite) TearDownSuite() {
	s.FetchServer.Close()
}

func TestFetcherTestSuite(t *testing.T) {
	suite.Run(t, new(FetcherTestSuite))
}

func (s *FetcherTestSuite) TestFetcherReceive() {
	imageType := "kvm"
	imageComment := "test image"

	// Create the upload request
	req, _ := http.NewRequest("GET", "http://localhost", bytes.NewReader(s.ImageData))

	_, err := s.Context.Fetcher.Receive(req)
	s.Error(err, "missing type header should fail")

	req.Header.Add("X-Image-Type", "kvm")
	req.Header.Add("X-Image-Comment", "test image")

	image, err := s.Context.Fetcher.Receive(req)
	s.NoError(err)
	s.EqualValues(len(s.ImageData), image.Size, "sizes should match")
	s.Equal(imageType, image.Type, "types should match")
	s.Equal(imageComment, image.Comment, "comments should match")
	s.NotEmpty(image.ID, "id should have been assigned")
}

func (s *FetcherTestSuite) TestFetcherFetch() {
	imageReq := &metadata.Image{
		ID: metadata.NewID(),
	}

	_, err := s.Context.Fetcher.Fetch(imageReq)
	s.Error(err, "missing source should error")
	imageReq.Source = "asdf"

	_, err = s.Context.Fetcher.Fetch(imageReq)
	s.Error(err, "missing type should error")
	imageReq.Type = "kvm"

	tests := []struct {
		source      string
		finalStatus string
	}{
		{"asdf", metadata.StatusError},
		{s.FetchServer.URL + "/404", metadata.StatusError},
		{s.FetchServer.URL, metadata.StatusComplete},
	}

	for _, test := range tests {
		imageReq = &metadata.Image{
			ID:     metadata.NewID(),
			Source: test.source,
			Type:   "kvm",
		}
		imageReq.Source = test.source
		image, err := s.Context.Fetcher.Fetch(imageReq)
		s.NoError(err, "valid config should have no initial error")
		s.Equal(metadata.StatusPending, image.Status, "new image should start out pending")
		for i := 0; i < 300; i++ {
			image, err := s.Context.MetadataStore.GetByID(image.ID)
			s.NotNil(image, "image should not be nil")
			s.NoError(err)

			if image.Status == metadata.StatusComplete || image.Status == metadata.StatusError {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		s.Equal(test.finalStatus, image.Status, "final status should be expected for source %s", test.source)
		if test.finalStatus == metadata.StatusComplete {
			s.EqualValues(len(s.ImageData), image.Size, "final size should be expected")
		}
	}

	image, err := s.Context.Fetcher.Fetch(imageReq)
	s.NoError(err)
	s.NotNil(image)
	s.Equal(metadata.StatusComplete, image.Status, "previously fetched image should be returned ready")
}
