package prom

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Prometheus struct {
	registry  *prometheus.Registry
	whitelist map[string]bool
	blacklist map[string]bool
	config    *viper.Viper
	logger    *zap.Logger
}

func (r *Prometheus) Handler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.InstrumentMetricHandler(
		r.registry,
		promhttp.HandlerFor(r.registry, promhttp.HandlerOpts{
			Registry:      r.registry,
			ErrorHandling: promhttp.ContinueOnError,
			ErrorLog:      &log{logger: r.logger},
		}),
	))
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
