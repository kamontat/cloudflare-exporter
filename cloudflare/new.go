//go:generate go run github.com/Khan/genqlient

package cloudflare

import (
	"github.com/Khan/genqlient/graphql"
	cf "github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/utils"
	"github.com/spf13/viper"
)

const (
	API_URL = "https://api.cloudflare.com/client/v4"
)

func NewAPI(config *viper.Viper) *cf.API {
	debug := config.GetBool(configs.CONF_DEBUG_MODE)
	token := config.GetString(configs.CONF_CF_API_TOKEN)
	if token != "" {
		return utils.CheckErrorWithData(cf.NewWithAPIToken(
			token,
			cf.BaseURL(API_URL),
			cf.Debug(debug),
		))
	} else {
		return utils.CheckErrorWithData(cf.New(
			config.GetString(configs.CONF_CF_API_KEY),
			config.GetString(configs.CONF_CF_API_EMAIL),
			cf.BaseURL(API_URL),
			cf.Debug(debug),
		))
	}
}

func NewGraphQL(config *viper.Viper) graphql.Client {
	client := NewHttpClient(config)
	return graphql.NewClient(API_URL, client)
}
