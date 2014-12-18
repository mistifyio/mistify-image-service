package main

import (
    flag "github.com/docker/docker/pkg/mflag"
    "os"
    "github.com/mistifyio/mistify-agent/log"
    "github.com/mistifyio/mistify-image-service"
    "strings"
)

const (
    DEFAULT_CONFIG_FILE = "config.json"
)

func main() {
    var help bool
    var configFile string

    flag.BoolVar(&help, []string{"h", "#help", "-help"}, false, "display help")
    flag.StringVar(&configFile, []string{"c", "#config-file", "-config-file"}, DEFAULT_CONFIG_FILE, "configuration file")

    flag.Parse()

    if help {
        flag.PrintDefaults()
        os.Exit(0)
    }

    config, err := imageservice.ConfigFromFile(configFile)
    if nil != err {
        log.Fatal(err)
    }

    context, err := imageservice.NewContext(config)
    if nil != err {
        log.Fatal(err)
    }

    err = imageservice.Run(context, strings.Join([]string{config.Address, config.Port}, ":"))
    if nil != err {
        log.Fatal(err)
    }

}
