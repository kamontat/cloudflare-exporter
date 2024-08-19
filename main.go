package main

import (
	"os"
	"os/signal"

	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/kamontat/cloudflare-exporter/prom"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// From goreleaser
var (
	name      string = "cf-exporter"
	version   string = "v0.0.0"
	date      string = "<date>"
	gitCommit string = "<commit>"
	gitState  string = "dirty"
	builtBy   string = "manual"
)

var (
	metadata *configs.Metadata
	config   *viper.Viper
	logger   *zap.Logger
)

func init() {
	metadata = &configs.Metadata{
		Name:      name,
		Version:   version,
		Date:      date,
		GitCommit: gitCommit,
		GitState:  gitState,
		BuiltBy:   builtBy,
	}

	config = configs.New(metadata)
	logger = loggers.SetDefault(loggers.New(config))
}

func main() {
	defer logger.Sync()

	// Initiate cloudflare and prometheus
	prometheus := prom.New(config)
	client := cloudflare.New(config)

	app := NewApp(client, prometheus)
	app.Setup()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Start application
	go func() { app.Start() }()

	// Wait and shutdown application
	<-c
	logger.Info("Gracefully shutting down...")
	app.Shutdown()
}
