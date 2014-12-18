package imageservice

import (
    "io/ioutil"
    "encoding/json"
)

type (

    Config struct {
        ImageStoreType string `json:"imageStoreType"`
        ImageStoreConfig map[string]string `json:"imageStoreConfig"`
        MetadataStoreType string `json:"metadataStoreType"`
        MetadataStoreConfig map[string] string `json:"metadataStoreConfig"`
    }

)

func ConfigFromFile(path string) (*Config, error) {
    data, err := ioutil.ReadFile(path)
    if nil != err {
        return nil, err
    }

    config := &Config{
        ImageStoreType: "",
        ImageStoreConfig: make(map[string]string),
        MetadataStoreType: "",
        MetadataStoreConfig: make(map[string]string),
    }

    err = json.Unmarshal(data, &config)
    if nil != err {
        return nil, err
    }

    return config, nil
}
