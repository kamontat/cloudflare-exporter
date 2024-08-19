package loggers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	fiberLogger "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func GocronLoggerAdapter(logger *zap.Logger) *gocronLoggerAdapter {
	return &gocronLoggerAdapter{wrapper: logger}
}

func FiberLoggerAdapter(logger *zap.Logger, config *viper.Viper) fiber.Handler {
	return fiberLogger.New(fiberLogger.Config{
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
	})
}

type gocronLoggerAdapter struct {
	wrapper *zap.Logger
}

func (l *gocronLoggerAdapter) Debug(msg string, args ...any) {
	l.wrapper.Debug(fmt.Sprintf(msg, args...))
}
func (l *gocronLoggerAdapter) Info(msg string, args ...any) {
	l.wrapper.Info(fmt.Sprintf(msg, args...))
}
func (l *gocronLoggerAdapter) Warn(msg string, args ...any) {
	l.wrapper.Warn(fmt.Sprintf(msg, args...))
}
func (l *gocronLoggerAdapter) Error(msg string, args ...any) {
	l.wrapper.Error(fmt.Sprintf(msg, args...))
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
