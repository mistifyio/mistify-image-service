package imageservice_test

import (
	"encoding/json"

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
}

func (s *ContextTestSuite) SetupSuite() {
	log.SetLevel(log.FatalLevel)

	// Images Store Setup
	images.Register("mock", func() images.Store {
		return &imocks.Store{}
	})
	viper.SetDefault("imageStoreType", "mock")
	viper.SetDefault("imageStoreConfig", nil)

	// Metadata Store Setup
	metadata.Register("mock", func() metadata.Store {
		return &mmocks.Store{}
	})
	viper.SetDefault("metadataStoreType", "mock")
	viper.SetDefault("metadataStoreConfig", nil)
}

func (s *ContextTestSuite) TestContextInitImageStore() {
	context := &imageservice.Context{}

	s.Error(context.InitImageStore("asdfwqfas", nil))
	s.Nil(context.ImageStore)

	imageStoreType := viper.GetString("imageStoreType")
	imageStoreConfig, _ := json.Marshal(viper.Get("imageStoreConfig"))
	s.NoError(context.InitImageStore(imageStoreType, imageStoreConfig))
	s.NotNil(context.ImageStore)
}

func (s *ContextTestSuite) TestContextNewMetadataStore() {
	context := &imageservice.Context{}

	s.Error(context.InitMetadataStore("asdfwqfas", nil))
	s.Nil(context.MetadataStore)

	metadataStoreType := viper.GetString("metadataStoreType")
	metadataStoreConfig, _ := json.Marshal(viper.Get("metadataStoreConfig"))
	s.NoError(context.InitMetadataStore(metadataStoreType, metadataStoreConfig))
	s.NotNil(context.MetadataStore)
}

func (s *ContextTestSuite) TestNewContext() {
	context, err := imageservice.NewContext()
	s.NoError(err)
	s.NotNil(context)
	s.NotNil(context.ImageStore)
	s.NotNil(context.MetadataStore)
	s.NotNil(context.Fetcher)
}
