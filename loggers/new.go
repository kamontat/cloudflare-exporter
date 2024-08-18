package loggers

import (
	"sync/atomic"

	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger atomic.Pointer[zap.Logger]

func init() {
	defaultLogger.Store(zap.NewNop())
}

func Default() *zap.Logger {
	return defaultLogger.Load()
}

func SetDefault(l *zap.Logger) *zap.Logger {
	defaultLogger.Store(l)
	return l
}

func New(config *viper.Viper) *zap.Logger {
	var fields = make(map[string]interface{})

	var encode = "console"
	var encodeLvl = zapcore.CapitalColorLevelEncoder
	if config.GetBool(configs.CONF_OUTPUT_JSON) {
		encode = "json"
		encodeLvl = zapcore.CapitalLevelEncoder
	}

	var defaultEncoder = &zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "fn",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLvl,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if config.GetBool(configs.CONF_DEBUG_MODE) {
		return utils.CheckErrorWithData((&zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: true,
			Encoding:          encode,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
			EncoderConfig:     *defaultEncoder,
			InitialFields:     fields,
		}).Build())
	} else if config.GetBool(configs.CONF_SILENT_MODE) {
		return zap.NewNop()
	} else {
		return utils.CheckErrorWithData((&zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     true,
			DisableStacktrace: true,
			Encoding:          encode,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig:    *defaultEncoder,
			InitialFields:    fields,
		}).Build())
	}
}
