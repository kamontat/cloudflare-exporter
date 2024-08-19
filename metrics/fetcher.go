package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func New(ctx context.Context, client *cloudflare.Client, config *viper.Viper, job gocron.Job) *fetcher {
	logger := loggers.Default()
	lastTime, err1 := job.LastRun()
	nextTime, err2 := job.NextRun()

	var now = time.Now()
	var last = now.Add(-config.GetDuration(configs.CONF_CF_INTERVAL))
	if err1 == nil && err2 == nil {
		diff := nextTime.Sub(lastTime)
		now = lastTime
		last = now.Add(-diff)
	}

	return &fetcher{
		context: ctx,
		client:  client,
		logger:  logger,
		config:  config,

		wg:   sync.WaitGroup{},
		last: last,
		now:  now,
	}
}

type fetcher struct {
	context context.Context
	client  *cloudflare.Client
	logger  *zap.Logger
	config  *viper.Viper
	wg      sync.WaitGroup
	last    time.Time
	now     time.Time
}

func (f *fetcher) Wait() {
	f.wg.Wait()
}
