package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service"
	logx "github.com/mistifyio/mistify-logrus-ext"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var configFile, logLevel string

	flag.IntP("port", "a", 20000, "listen address")
	flag.StringVarP(&logLevel, "log-level", "l", "warning", "log level: debug/info/warning/error/critical/fatal")
	flag.StringVarP(&configFile, "config-file", "c", "", "config file")
	flag.Parse()

	if err := logx.DefaultSetup(logLevel); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "logx.DefaultSetup",
			"level": logLevel,
		}).Fatal("failed to set up logging")
	}

	if configFile == "" {
		log.Fatal("undefined config file")
	}

	viper.SetConfigFile(configFile)
	_ = viper.BindPFlag("port", flag.Lookup("port"))
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("error", err).Fatal("failed to load config")
	}

	context, err := imageservice.NewContext()
	if nil != err {
		log.Fatal("failed to create and initialize context")
	}

	log.WithFields(log.Fields{
		"port": viper.GetInt("port"),
	}).Info("running server")

	server := imageservice.Run(context, viper.GetInt("port"))
	// Block until the server is stopped
	<-server.StopChan()
}
