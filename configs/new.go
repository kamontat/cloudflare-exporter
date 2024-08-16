package configs

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func New(meta *Metadata) *viper.Viper {
	SetMetadata(meta)

	var v = viper.NewWithOptions(
		viper.WithLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelError,
		}))),
		viper.EnvKeyReplacer(strings.NewReplacer("-", "_")),
	)

	v.AutomaticEnv()
	v.SetEnvPrefix(SHORT_NAME)
	v.BindEnv(CONF_DEBUG_MODE, "DEBUG")   // special environment to enable debug
	v.BindEnv(CONF_SILENT_MODE, "SILENT") // special environment to enable silent mode

	v.SetConfigName(SETTING_CONFIG_FILE)
	v.SetConfigType(SETTING_CONFIG_EXT)
	v.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", meta.Name))
	v.AddConfigPath(".")

	return v
}
