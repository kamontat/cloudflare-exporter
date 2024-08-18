package main

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var (
	registry *prometheus.Registry
)

func init() {
	registry = prometheus.NewRegistry()
	registry.MustRegister(collectors.NewBuildInfoCollector())
	registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registry.MustRegister(collectors.NewGoCollector())
}

type log struct {
	logger *zap.Logger
}

func (l *log) Println(v ...interface{}) {
	if len(v) == 2 {
		var msg = fmt.Sprintf("%v", v[0])
		l.logger.Warn(msg, zap.Any("error", v[1]))
	} else if len(v) == 1 {
		var msg = fmt.Sprintf("%v", v[0])
		l.logger.Warn(msg)
	}

	l.logger.Warn(fmt.Sprint(v...))
}

func PromHttpHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			Registry:      registry,
			ErrorHandling: promhttp.ContinueOnError,
			ErrorLog:      &log{logger: loggers.Default()},
		}),
	))
}
