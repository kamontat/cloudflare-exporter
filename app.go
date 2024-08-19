package main

import (
	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/prom"
	"go.uber.org/zap"
)

func NewApp(client *cloudflare.Client, prometheus *prom.Prometheus) *app {
	fields := metadata.ToFields()
	fields = append(fields, zap.Bool("production", config.GetBool(configs.CONF_PRODUCTION)))
	logger.Info("Start application", fields...)

	metrics := newMetricSet(client, prometheus, config)
	return &app{
		server:     newServer(),
		scheduler:  newScheduler(metrics),
		client:     client,
		prometheus: prometheus,
		metrics:    metrics,
	}
}

type app struct {
	server    *fiber.App
	scheduler gocron.Scheduler

	client     *cloudflare.Client
	prometheus *prom.Prometheus
	metrics    *MetricSet
}

func (a *app) Setup() {
	setupRootPath(a.server)
	setupHealthPath(a.server, a.client)
	setupMetricPath(a.server, a.prometheus)

	setupJob(a.scheduler, JOB_NAME_DEFAULT,
		gocron.DurationJob(config.GetDuration(configs.CONF_CF_INTERVAL)),
		a.metrics.Fetch,
	)
}

func (a *app) Start() {
	startScheduler(a.scheduler)
	startServer(a.server)
}

func (a *app) Shutdown() {
	shutdownScheduler(a.scheduler)
	shutdownServer(a.server)
}
