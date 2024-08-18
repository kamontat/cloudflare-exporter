package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/timeout"
	"github.com/kamontat/cloudflare-exporter/configs"
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
	app.Use(fiberLogger.New(fiberLogger.Config{
		DisableColors: config.GetBool(configs.CONF_OUTPUT_JSON),
		Next: func(c fiber.Ctx) bool {
			return c.Path() == config.GetString(configs.CONF_SERVER_HEALTH_PATH) ||
				c.Path() == config.GetString(configs.CONF_SERVER_LIVENESS_PATH) ||
				c.Path() == config.GetString(configs.CONF_SERVER_READINESS_PATH)
		},
		LoggerFunc: func(c fiber.Ctx, data *fiberLogger.Data, cfg fiberLogger.Config) error {
			var colors *fiber.Colors
			if !cfg.DisableColors {
				c := c.App().Config().ColorScheme
				colors = &c
			}

			err := ""
			if data.ChainErr != nil {
				err = data.ChainErr.Error()
			}

			logger.Info(fmt.Sprintf("%s %s: %s %s %13v %s",
				methodColor(c.Method(), colors),
				c.Path(),
				statusColor(c.Response().StatusCode(), colors),
				c.IP(),
				data.Stop.Sub(data.Start),
				err,
			),
				zap.String("RequestID", requestid.FromContext(c)),
			)
			return nil
		},
	}))

	return app
}

func rootPath(server *fiber.App) {
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

func metricPath(server *fiber.App) {
	var metricPath = config.GetString(configs.CONF_SERVER_METRIC_PATH)
	server.Get(metricPath, PromHttpHandler())
}

func healthPath(server *fiber.App) {
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
			token, err := api.VerifyAPIToken(c.Context())
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

func startServer(server *fiber.App, scheduler gocron.Scheduler) {
	addr := fmt.Sprintf(
		"%s:%d",
		config.GetString(configs.CONF_SERVER_ADDR),
		config.GetInt(configs.CONF_SERVER_PORT),
	)

	scheduler.Start()
	utils.CheckError(server.Listen(addr, fiber.ListenConfig{
		EnablePrintRoutes: config.GetBool(configs.CONF_DEBUG_MODE),
		OnShutdownError: func(err error) {
			logger.Error(err.Error())
		},
		OnShutdownSuccess: func() {
			logger.Info("Shutdown server successfully")
		},
	}))
}

func shutdownServer(server *fiber.App, scheduler gocron.Scheduler) {
	logger.Info("Gracefully shutting down...")
	logger.Info("Shutdown server", zap.Error(server.Shutdown()))
	logger.Info("Shutdown cronjob scheduler", zap.Error(scheduler.Shutdown()))
}

func parseDataSize(key string) int {
	return utils.CheckErrorWithData(units.ParseDataSize(config.GetString(key))).Byte()
}

func parseDuration(key string) time.Duration {
	return utils.CheckErrorWithData(time.ParseDuration(config.GetString(key)))
}

func methodColor(method string, colors *fiber.Colors) string {
	if colors == nil {
		return method
	}
	var color string
	switch method {
	case fiber.MethodGet:
		color = colors.Cyan
	case fiber.MethodPost:
		color = colors.Green
	case fiber.MethodPut:
		color = colors.Yellow
	case fiber.MethodDelete:
		color = colors.Red
	case fiber.MethodPatch:
		color = colors.White
	case fiber.MethodHead:
		color = colors.Magenta
	case fiber.MethodOptions:
		color = colors.Blue
	default:
		color = colors.Reset
	}
	return fmt.Sprintf("%s%s%s", color, method, colors.Reset)
}

func statusColor(code int, colors *fiber.Colors) string {
	if colors == nil {
		return strconv.Itoa(code)
	}
	var color string
	switch {
	case code >= fiber.StatusOK && code < fiber.StatusMultipleChoices:
		color = colors.Green
	case code >= fiber.StatusMultipleChoices && code < fiber.StatusBadRequest:
		color = colors.Blue
	case code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError:
		color = colors.Yellow
	default:
		color = colors.Red
	}
	return fmt.Sprintf("%s%d%s", color, code, colors.Reset)
}
