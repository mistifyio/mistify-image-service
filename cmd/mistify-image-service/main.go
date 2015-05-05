package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/mistifyio/mistify-image-service"
	logx "github.com/mistifyio/mistify-logrus-ext"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	var configFile string

	viper.SetDefault("port", uint(20000))
	viper.SetDefault("log-level", "warning")

	flag.IntP("port", "a", 0, "listen address")
	flag.StringP("log-level", "l", "", "log level: debug/info/warning/error/critical/fatal")
	flag.StringVarP(&configFile, "config-file", "c", "", "config file")
	flag.Parse()

	_ = viper.BindPFlag("port", flag.Lookup("port"))
	_ = viper.BindPFlag("logLevel", flag.Lookup("log-level"))

	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.WithField("error", err).Fatal(err)
	}

	if err := logx.DefaultSetup(viper.GetString("logLevel")); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"func":  "logx.DefaultSetup",
			"level": viper.GetString("logLevel"),
		}).Fatal("failed to set up logging")
	}

	context, err := imageservice.NewContext()
	if nil != err {
		log.Fatal("failed to create and initialize context")
	}

	if err := imageservice.Run(context, viper.GetInt("port")); err != nil {
		log.WithField("error", err).Fatal("failed to run server")
	}
}
