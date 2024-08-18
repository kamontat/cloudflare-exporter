package configs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kamontat/cloudflare-exporter/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func New(meta *Metadata) *viper.Viper {
	SetMetadata(meta)

	var level = slog.LevelDebug
	if meta.BuiltBy == "goreleaser" {
		level = slog.LevelWarn
	}

	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}))

	if _, err := os.Stat(".env"); err == nil {
		logger.Info("Load .env file")
		utils.CheckError(godotenv.Load())
	}

	var v = viper.NewWithOptions(
		viper.WithLogger(logger),
		viper.EnvKeyReplacer(strings.NewReplacer("-", "_")),
	)

	v.SetDefault(CONF_DEBUG_MODE, false)
	v.SetDefault(CONF_SILENT_MODE, false)
	v.SetDefault(CONF_OUTPUT_JSON, false)
	v.SetDefault(CONF_SERVER_ADDR, "")
	v.SetDefault(CONF_SERVER_PORT, 3000)
	v.SetDefault(CONF_SERVER_METRIC_PATH, "/metrics")
	v.SetDefault(CONF_SERVER_HEALTH_PATH, "/health")
	v.SetDefault(CONF_SERVER_LIVENESS_PATH, "/health/liveness")
	v.SetDefault(CONF_SERVER_READINESS_PATH, "/health/readiness")

	v.SetDefault(CONF_SERVER_CONCURRENCY, 100)
	v.SetDefault(CONF_SERVER_BODY_LIMIT, "10MB")
	v.SetDefault(CONF_SERVER_READ_BUFFER, "10KB")
	v.SetDefault(CONF_SERVER_WRITE_BUFFER, "10KB")
	v.SetDefault(CONF_SERVER_READ_TIMEOUT, "1s")
	v.SetDefault(CONF_SERVER_WRITE_TIMEOUT, "10s")
	v.SetDefault(CONF_SERVER_IDLE_TIMEOUT, "5m")
	v.SetDefault(CONF_CF_API_TOKEN, "")
	v.SetDefault(CONF_CF_API_EMAIL, "")
	v.SetDefault(CONF_CF_API_KEY, "")
	v.SetDefault(CONF_CF_TIMEOUT, "5s")

	v.AutomaticEnv()

	v.BindPFlags(setupFlags())

	v.SetConfigName(SETTING_CONFIG_FILE)
	v.AddConfigPath(fmt.Sprintf("/etc/%s", meta.Name))
	v.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", meta.Name))
	v.AddConfigPath(".")

	err := v.ReadInConfig()
	if err != nil {
		logger.Warn(err.Error())
	}
	return v
}

func setupFlags() *pflag.FlagSet {
	var flagset = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	flagset.BoolP(CONF_DEBUG_MODE, "D", false, "Enabled debug mode")
	flagset.BoolP(CONF_SILENT_MODE, "S", false, "Enabled silent mode")
	flagset.BoolP(CONF_OUTPUT_JSON, "J", false, "Print output as json format")

	flagset.StringP(CONF_SERVER_ADDR, "a", "", "Server address")
	flagset.IntP(CONF_SERVER_PORT, "p", 0, "Server port number")
	flagset.StringP(CONF_SERVER_METRIC_PATH, "m", "", "The Prometheus metrics path")
	flagset.StringP(CONF_SERVER_HEALTH_PATH, "h", "", "The healthcheck path")
	flagset.StringP(CONF_SERVER_LIVENESS_PATH, "l", "", "The Kubernetes liveness health path")
	flagset.StringP(CONF_SERVER_READINESS_PATH, "n", "", "The Kubernetes readiness health path")

	flagset.IntP(CONF_SERVER_CONCURRENCY, "c", 0, "Maximum number of concurrent connections")
	flagset.StringP(CONF_SERVER_BODY_LIMIT, "b", "", "Max body size that the server accepts")
	flagset.StringP(CONF_SERVER_READ_BUFFER, "r", "", "Per-connection buffer size for requests' reading")
	flagset.StringP(CONF_SERVER_WRITE_BUFFER, "w", "", "Per-connection buffer size for responses' writing")
	flagset.StringP(CONF_SERVER_READ_TIMEOUT, "e", "", "The amount of time allowed to read the full request including body")
	flagset.StringP(CONF_SERVER_WRITE_TIMEOUT, "i", "", "The maximum duration before timing out writes of the response")
	flagset.StringP(CONF_SERVER_IDLE_TIMEOUT, "d", "", "The maximum amount of time to wait for the next request when keep-alive is enabled")

	flagset.StringP(CONF_CF_API_TOKEN, "T", "", "https://developers.cloudflare.com/fundamentals/api/get-started/create-token")
	flagset.StringP(CONF_CF_API_EMAIL, "E", "", "https://developers.cloudflare.com/fundamentals/api/get-started/keys/")
	flagset.StringP(CONF_CF_API_KEY, "K", "", "https://developers.cloudflare.com/fundamentals/api/get-started/keys/")
	flagset.StringP(CONF_CF_TIMEOUT, "t", "", "How long should we wait for cloudflare to response")

	utils.CheckError(flagset.Parse(os.Args[1:]))
	return flagset
}
