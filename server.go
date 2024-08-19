package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/timeout"
	"github.com/kamontat/cloudflare-exporter/cloudflare"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/kamontat/cloudflare-exporter/prom"
	"github.com/kamontat/cloudflare-exporter/units"
	"github.com/kamontat/cloudflare-exporter/utils"
	"go.uber.org/zap"
)

func newServer() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         metadata.Name,
		ServerHeader:    fmt.Sprintf("%s %s", metadata.Name, metadata.Version),
		Concurrency:     config.GetInt(configs.CONF_SERVER_CONCURRENCY),
		BodyLimit:       parseDataSize(configs.CONF_SERVER_BODY_LIMIT),
		ReadBufferSize:  parseDataSize(configs.CONF_SERVER_READ_BUFFER),
		WriteBufferSize: parseDataSize(configs.CONF_SERVER_WRITE_BUFFER),
		ReadTimeout:     parseDuration(configs.CONF_SERVER_READ_TIMEOUT),
		WriteTimeout:    parseDuration(configs.CONF_SERVER_WRITE_TIMEOUT),
		IdleTimeout:     parseDuration(configs.CONF_SERVER_IDLE_TIMEOUT),
		ColorScheme:     fiber.Colors{},
	})

	app.Use(requestid.New(requestid.ConfigDefault))
	app.Use(compress.New(compress.ConfigDefault))
	app.Use(helmet.New(helmet.ConfigDefault))
	app.Use(loggers.FiberLoggerAdapter(logger, config))

	return app
}

func setupRootPath(server *fiber.App) {
	var metricPath = config.GetString(configs.CONF_SERVER_METRIC_PATH)
	var healthPath = config.GetString(configs.CONF_SERVER_HEALTH_PATH)
	var livenessPath = config.GetString(configs.CONF_SERVER_LIVENESS_PATH)
	var readinessPath = config.GetString(configs.CONF_SERVER_READINESS_PATH)
	server.Get("/", func(c fiber.Ctx) error {
		return c.SendString(
			fmt.Sprintf("Metrics is %s, Health: %s (liveness=%s, readiness=%s)",
				metricPath,
				healthPath,
				livenessPath,
				readinessPath,
			),
		)
	})
}

func setupMetricPath(server *fiber.App, prometheus *prom.Prometheus) {
	var metricPath = config.GetString(configs.CONF_SERVER_METRIC_PATH)
	server.Get(metricPath, prometheus.Handler())
}

func setupHealthPath(server *fiber.App, client *cloudflare.Client) {
	var healthPath = config.GetString(configs.CONF_SERVER_HEALTH_PATH)
	server.Get(healthPath, healthcheck.NewHealthChecker())

	var livenessPath = config.GetString(configs.CONF_SERVER_LIVENESS_PATH)
	server.Get(livenessPath, healthcheck.NewHealthChecker(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			return true
		},
	}))

	var readinessPath = config.GetString(configs.CONF_SERVER_READINESS_PATH)
	server.Get(readinessPath, timeout.New(healthcheck.NewHealthChecker(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			token, err := client.API.VerifyAPIToken(c.Context())
			if err != nil {
				logger.Warn("Cannot verify cloudflare token", zap.Error(err))
				return false
			}
			if token.Status != "active" {
				logger.Warn("Cloudflare token is not active", zap.String("status", token.Status))
				return false
			}
			return true
		},
	}), config.GetDuration(configs.CONF_CF_TIMEOUT)))
}

func startServer(server *fiber.App) {
	addr := fmt.Sprintf(
		"%s:%d",
		config.GetString(configs.CONF_SERVER_ADDR),
		config.GetInt(configs.CONF_SERVER_PORT),
	)

	server.Hooks().OnListen(func(ld fiber.ListenData) error {
		schema := "http"
		if ld.TLS {
			schema = "https"
		}

		logger.Info(fmt.Sprintf("Listening server at %s://%s:%s", schema, ld.Host, ld.Port))
		return nil
	})

	utils.CheckError(server.Listen(addr, fiber.ListenConfig{
		EnablePrintRoutes:     config.GetBool(configs.CONF_DEBUG_MODE),
		DisableStartupMessage: config.GetBool(configs.CONF_PRODUCTION),
		OnShutdownError: func(err error) {
			logger.Error(err.Error())
		},
		OnShutdownSuccess: func() {
			logger.Info("Shutdown server successfully")
		},
	}))
}

func shutdownServer(server *fiber.App) {
	logger.Info("Shutdown server", zap.Error(server.Shutdown()))
}

func parseDataSize(key string) int {
	return utils.CheckErrorWithData(units.ParseDataSize(config.GetString(key))).Byte()
}

func parseDuration(key string) time.Duration {
	return utils.CheckErrorWithData(time.ParseDuration(config.GetString(key)))
}
