package imageservice_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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
	logx "github.com/mistifyio/mistify-logrus-ext"
	"github.com/pborman/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/tylerb/graceful"
)

type APITestSuite struct {
	suite.Suite
	Port        int
	StoreDir    string
	ImageData   []byte
	APIServer   *graceful.Server
	FetchServer *httptest.Server
	APIURL      string
}

func (s *APITestSuite) SetupSuite() {
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

	s.Port = 54321
	s.APIURL = fmt.Sprintf("http://localhost:%d/images", s.Port)
}

func (s *APITestSuite) SetupTest() {
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

	// Start API server
	s.APIServer = imageservice.Run(ctx, s.Port)
	time.Sleep(100 * time.Millisecond)
}

func (s *APITestSuite) TearDownTest() {
	// Stop API server
	stopChan := s.APIServer.StopChan()
	s.APIServer.Stop(5 * time.Second)
	<-stopChan

	// Cleanup store
	s.NoError(os.RemoveAll(s.StoreDir))
}

func (s *APITestSuite) TearDownSuite() {
	s.FetchServer.Close()
}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (s *APITestSuite) TestReceiveImage() {
	tests := []struct {
		description        string
		imageType          string
		expectedStatusCode int
	}{
		{"missing type should fail", "", http.StatusBadRequest},
		{"invalid type should fail", "asdf", http.StatusBadRequest},
		{"valid type kvm should succeed", "kvm", http.StatusOK},
		{"valid type container should succeed", "container", http.StatusOK},
	}

	for _, test := range tests {
		msg := func(val string) string {
			return test.description + " : " + val
		}

		image, resp, err := s.uploadImage(test.imageType)
		s.Equal(test.expectedStatusCode, resp.StatusCode, msg("status code should be expected"))
		if test.expectedStatusCode != http.StatusOK {
			s.Nil(image)
			continue
		}

		s.NoError(err, msg("upload shouldn't error"))
		s.NotEmpty(image.ID, msg("should have ID assigned"))
		s.Equal(metadata.StatusComplete, image.Status, msg("final status should be complete"))
		s.EqualValues(len(s.ImageData), image.Size, msg("final size should be expected"))
	}
}

func (s *APITestSuite) TestFetchImage() {
	tests := []struct {
		description        string
		requestData        []byte
		expectedStatusCode int
	}{
		{"bad json should fail",
			[]byte("asdf"), http.StatusBadRequest},
		{"missing source should fail",
			[]byte(`{}`), http.StatusBadRequest},
		{"missing image type should fail",
			[]byte(fmt.Sprintf(`{"source":"%s"}`, s.FetchServer.URL)), http.StatusBadRequest},
		{"invalid image type should fail",
			[]byte(fmt.Sprintf(`{"source":"%s","type":"asdf"}`, s.FetchServer.URL)), http.StatusBadRequest},
		{"complete kvm request should succeed",
			[]byte(fmt.Sprintf(`{"source":"%s","type":"kvm"}`, s.FetchServer.URL)), http.StatusAccepted},
		{"complete container request should succeed",
			// Modify url slightly to prevent it using prefetched image
			[]byte(fmt.Sprintf(`{"source":"%s","type":"container"}`, s.FetchServer.URL+"?")), http.StatusAccepted},
	}

	for _, test := range tests {
		msg := func(val string) string {
			return test.description + " : " + val
		}

		resp, err := http.Post(s.APIURL, "application/json", bytes.NewBuffer(test.requestData))
		s.NoError(err, msg("fetch request shouldn't error"))
		s.Equal(test.expectedStatusCode, resp.StatusCode, msg("fetch response should be accepted"))
		defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")

		if test.expectedStatusCode != http.StatusAccepted {
			continue
		}

		image, err := unmarshalImageResp(resp)
		s.NoError(err, msg("resp body should be valid image json"))
		s.NotEmpty(image.ID, msg("should have ID assigned"))
		s.Equal(metadata.StatusPending, image.Status, msg("status should start out pending"))

		// Poll until the fetch request completes or errors
		finalImage := &metadata.Image{}
		for i := 0; i < 300; i++ {
			finalImage, _, err = s.getImage(image.ID)
			s.NoError(err)
			s.NotNil(image)

			if finalImage.Status == metadata.StatusComplete || finalImage.Status == metadata.StatusError {
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		s.Equal(metadata.StatusComplete, finalImage.Status, msg("final status should be complete"))
		s.EqualValues(len(s.ImageData), finalImage.Size, msg("final size should be expected"))

	}
}

func (s *APITestSuite) TestListImages() {
	imageKVM, _, _ := s.uploadImage("kvm")
	imageContainer, _, _ := s.uploadImage("container")

	tests := []struct {
		description        string
		imageType          string
		expectedStatusCode int
		expectedImages     []*metadata.Image
	}{
		{"no filter should list all images",
			"", http.StatusOK, []*metadata.Image{imageKVM, imageContainer}},
		{"kvm filter should list only kvm images",
			"kvm", http.StatusOK, []*metadata.Image{imageKVM}},
		{"container filter should list only container images",
			"container", http.StatusOK, []*metadata.Image{imageContainer}},
		{"invalid filter should error",
			"asdf", http.StatusBadRequest, nil},
	}

	for _, test := range tests {
		msg := func(val string) string {
			return test.description + " : " + val
		}

		resp, err := http.Get(fmt.Sprintf("%s?type=%s", s.APIURL, test.imageType))
		s.NoError(err, msg("request should not error"))
		s.Equal(test.expectedStatusCode, resp.StatusCode, msg("http status codes should match"))
		defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")

		if test.expectedImages == nil {
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		s.NoError(err, msg("reading resp body should not error"))

		var images []*metadata.Image
		s.NoError(json.Unmarshal(body, &images), msg("response body should be valid image json array"))
		s.Equal(len(test.expectedImages), len(images), msg("list should match expected in length"))
		for _, expectedImage := range test.expectedImages {
			found := false
			for _, image := range images {
				if expectedImage.ID == image.ID {
					found = true
					break
				}
			}
			s.True(found, msg("expected image should be found in list"))
		}
	}
}

func (s *APITestSuite) TestGetImage() {
	imageKVM, _, err := s.uploadImage("kvm")
	s.NoError(err)

	image, resp, err := s.getImage("asdf")
	s.Error(err)
	s.Equal(http.StatusNotFound, resp.StatusCode)

	image, _, err = s.getImage(imageKVM.ID)
	s.NoError(err)
	s.NotNil(image)
	s.Equal(imageKVM.ID, image.ID, "id of image retrieved should match requested")
}

func (s *APITestSuite) TestDeleteImage() {
	imageKVM, _, _ := s.uploadImage("kvm")
	req, _ := http.NewRequest("DELETE", s.imageURL(imageKVM.ID), nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	s.NoError(err)
	defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")
	s.Equal(http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequest("DELETE", s.imageURL("asdf"), nil)
	resp, err = client.Do(req)
	s.NoError(err)
	defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")
	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *APITestSuite) TestDownloadImage() {
	imageKVM, _, _ := s.uploadImage("kvm")
	resp, err := http.Get(s.imageURL(imageKVM.ID) + "/download")
	s.NoError(err)
	defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")
	s.Equal(http.StatusOK, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	s.NoError(err)
	s.Equal(s.ImageData, body)

	resp, err = http.Get(s.imageURL("asdf") + "/download")
	s.NoError(err)
	defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")
	s.Equal(http.StatusNotFound, resp.StatusCode)

}

// uploadImage uploads the ImageData with valid properties
func (s *APITestSuite) uploadImage(imageType string) (*metadata.Image, *http.Response, error) {
	req, err := http.NewRequest("PUT", s.APIURL, bytes.NewBuffer(s.ImageData))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("X-Image-Type", imageType)
	req.Header.Add("X-Image-Comment", "uploaded image")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")

	if resp.StatusCode != http.StatusOK {
		return nil, resp, errors.New(fmt.Sprintf("unexpected status code %d", resp.StatusCode))
	}

	image, err := unmarshalImageResp(resp)
	return image, resp, err
}

// getImage retrieves image metadata
func (s *APITestSuite) getImage(id string) (*metadata.Image, *http.Response, error) {
	resp, err := http.Get(s.imageURL(id))
	if err != nil {
		return nil, resp, err
	}
	defer logx.LogReturnedErr(resp.Body.Close, nil, "failed to close fetch response body")

	if resp.StatusCode != http.StatusOK {
		return nil, resp, errors.New(fmt.Sprintf("unexpected status code %d", resp.StatusCode))
	}

	image, err := unmarshalImageResp(resp)
	return image, resp, err
}

// unmarshalImageResp turns a response body into a *metadata.Image
func unmarshalImageResp(resp *http.Response) (*metadata.Image, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	image := &metadata.Image{}
	if err := json.Unmarshal(body, image); err != nil {
		return nil, err
	}

	return image, nil
}

func (s *APITestSuite) imageURL(id string) string {
	return fmt.Sprintf("%s/%s", s.APIURL, id)
}
