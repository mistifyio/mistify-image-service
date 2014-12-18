package imageservice

import (
    "io/ioutil"
    "encoding/json"
)

const (
    DEFAULT_INTERFACE = ""
    DEFAULT_PORT = "9001"
)

type (

    Config struct {
        // HTTP server
        Address string `json:"address"`
        Port string `json:"port"`

        // images store
        ImageStoreType string `json:"imageStoreType"`
        ImageStoreConfig map[string]string `json:"imageStoreConfig"`

        // metadata store
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
        Address: "0.0.0.0",
        Port: "9001",
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
