package main

import (
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// From goreleaser
var (
	name      string = "cf-exporter"
	version   string = "v0.0.0"
	date      string = "<date>"
	gitCommit string = "<commit>"
	gitClean  string = "<is-clean>"
	builtBy   string = "manual"
)

var (
	metadata *configs.Metadata
	config   *viper.Viper
	logger   *zap.Logger
)

func init() {
	config = configs.New(&configs.Metadata{
		Name:      name,
		Version:   version,
		Date:      date,
		GitCommit: gitCommit,
		GitClean:  gitClean,
		BuiltBy:   builtBy,
	})

	metadata = configs.GetMetadata()
	logger = loggers.SetDefault(loggers.New(config))
	err := config.ReadInConfig()
	if err != nil {
		logger.Warn(err.Error())
	}
}

func main() {
	defer logger.Sync()

	logger.Info("Start application", metadata.ToFields()...)

	// app := fiber.New()
	// app.Get("/", func(c fiber.Ctx) error {
	// 	return c.SendString("Hello, World!")
	// })

	// app.Listen(":3000")
}
