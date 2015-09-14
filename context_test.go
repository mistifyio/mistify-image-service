package imageservice_test

import (
	"encoding/json"
	"errors"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service"
	"github.com/mistifyio/mistify-image-service/images"
	imocks "github.com/mistifyio/mistify-image-service/images/mocks"
	"github.com/mistifyio/mistify-image-service/metadata"
	mmocks "github.com/mistifyio/mistify-image-service/metadata/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type ContextTestSuite struct {
	suite.Suite
	ValidConfig   interface{}
	InvalidConfig interface{}
}

func (s *ContextTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)

	s.ValidConfig = map[string]string{"foo": "bar"}
	s.InvalidConfig = struct{}{}

	// Images Store Setup
	images.Register("mock", func() images.Store {
		m := &imocks.Store{}
		vj, _ := json.Marshal(s.ValidConfig)
		m.On("Init", vj).Return(nil)
		ij, _ := json.Marshal(s.InvalidConfig)
		m.On("Init", ij).Return(errors.New("asdf"))
		return m
	})

	// Metadata Store Setup
	metadata.Register("mock", func() metadata.Store {
		m := &mmocks.Store{}
		vj, _ := json.Marshal(s.ValidConfig)
		m.On("Init", vj).Return(nil)
		ij, _ := json.Marshal(s.InvalidConfig)
		m.On("Init", ij).Return(errors.New("asdf"))
		return m
	})
}

func (s *ContextTestSuite) SetupTest() {
	// Images Store Setup
	viper.Set("imageStoreType", "mock")
	viper.Set("imageStoreConfig", s.ValidConfig)

	// Metadata Store Setup
	viper.Set("metadataStoreType", "mock")
	viper.Set("metadataStoreConfig", s.ValidConfig)
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

func (s *ContextTestSuite) TestContextInitImageStore() {
	context := &imageservice.Context{}

	s.Error(context.InitImageStore("asdfwqfas", nil), "unknown type should fail")

	imageStoreType := viper.GetString("imageStoreType")
	imageStoreConfig, _ := json.Marshal(viper.Get("imageStoreConfig"))
	s.NoError(context.InitImageStore(imageStoreType, imageStoreConfig), "known type, valid config should succeed")
	s.NotNil(context.ImageStore, "known type, valid config should succeed")

	ij, _ := json.Marshal(s.InvalidConfig)
	s.Error(context.InitImageStore(imageStoreType, ij), "known type, invalid config should fail")
}

func (s *ContextTestSuite) TestContextNewMetadataStore() {
	context := &imageservice.Context{}

	s.Error(context.InitMetadataStore("asdfwqfas", nil), "unknown type should fail")
	s.Nil(context.MetadataStore)

	metadataStoreType := viper.GetString("metadataStoreType")
	metadataStoreConfig, _ := json.Marshal(viper.Get("metadataStoreConfig"))
	s.NoError(context.InitMetadataStore(metadataStoreType, metadataStoreConfig), "known type, valid config should succeed")
	s.NotNil(context.MetadataStore)

	ij, _ := json.Marshal(s.InvalidConfig)
	s.Error(context.InitMetadataStore(metadataStoreType, ij), "valid type, invalid config should fail")
}

func (s *ContextTestSuite) TestNewContext() {
	context, err := imageservice.NewContext()
	s.NoError(err, "valid store configs should succeed")
	s.NotNil(context)
	s.NotNil(context.ImageStore)
	s.NotNil(context.MetadataStore)
	s.NotNil(context.Fetcher)

	viper.Set("metadataStoreConfig", s.InvalidConfig)
	_, err = imageservice.NewContext()
	s.Error(err, "bad metadata config should fail")

	viper.Set("imageStoreConfig", s.InvalidConfig)
	_, err = imageservice.NewContext()
	s.Error(err, "bad image config should fail")

}
