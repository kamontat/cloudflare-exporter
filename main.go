package main

import (
	"os"
	"os/signal"
	"time"

	"github.com/Khan/genqlient/graphql"
	cf "github.com/cloudflare/cloudflare-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/kamontat/cloudflare-exporter/utils"
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
	api      *cf.API
	gql      graphql.Client
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
	api = cloudflare.NewAPI(config)
	gql = cloudflare.NewGraphQL(config)
}

func main() {
	defer logger.Sync()

	logger.Info("Start application", metadata.ToFields()...)

	scheduler := utils.CheckErrorWithData(gocron.NewScheduler())
	server := newServer()

	rootPath(server)
	healthPath(server)
	metricPath(server)

	scheduler.NewJob(gocron.DurationJob(3*time.Second), gocron.NewTask(func() {
		logger.Info("print every 3 seconds")
	}))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { startServer(server, scheduler) }()

	<-c
	shutdownServer(server, scheduler)
}
