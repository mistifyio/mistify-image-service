package imageservice

type (

    Config struct {
        ImageStoreType string `json:"imageStoreType"`
        ImageStoreConfig map[string]string `json:"imageStoreConfig"`
        MetadataStoreType string `json:"metadataStoreType"`
        MetadataStoreConfig map[string] string `json:"metadataStoreConfig"`
    }
    
)

func NewConfig() *Config {
    c := &Config{
            
    }
    
    return c
}