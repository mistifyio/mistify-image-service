package main

import (
	flag "github.com/docker/docker/pkg/mflag"
	"os"
    
    "github.com/mistifyio/mistify-image-service"
)

const (
    DEFAULT_CONFIG_FILE = "config.json"
)

func main() {
	var help bool

	flag.BoolVar(&help, []string{"h", "#help", "-help"}, false, "display help")
    flag.StringVar(&configFile, []string{"c", "#config-file", "-config-file"}, DEFAULT_CONFIG_FILE, "configuration file")

	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}
    
    context := imageservice.NewContext()
    
    imageservice.Run(context, "127.0.0.1")
}
