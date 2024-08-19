//go:generate go run github.com/Khan/genqlient configs/genqlient.yaml

package cloudflare

import (
	"context"
	"fmt"
	"strings"

	"github.com/Khan/genqlient/graphql"
	cf "github.com/cloudflare/cloudflare-go"
	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/kamontat/cloudflare-exporter/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	API_URL = "https://api.cloudflare.com/client/v4"
)

func New(config *viper.Viper) *Client {
	logger := loggers.Default()
	client := &Client{
		API:      newAPI(config),
		GQL:      newGraphQL(config),
		Accounts: make(map[string]cf.Account),
		Zones:    make(map[string]cf.Zone),

		loaded: false,
		logger: logger,
		config: config,
	}

	err := client.setup()
	if err != nil {
		if config.GetString(configs.CONF_CF_ERROR_MODE) == configs.ERROR_MODE_LOG {
			logger.Error("Cannot set up cloudflare client", zap.Error(err))
		} else if config.GetString(configs.CONF_CF_ERROR_MODE) == configs.ERROR_MODE_STOP {
			utils.CheckError(err)
		}
	}

	return client
}

type Client struct {
	API      *cf.API
	GQL      graphql.Client
	Accounts map[string]cf.Account
	Zones    map[string]cf.Zone

	loaded bool
	logger *zap.Logger
	config *viper.Viper
}

func (c *Client) setup() error {
	c.logger.Debug("Set up Cloudflare client")
	if c.loaded {
		c.logger.Debug("Cloudflare data was loaded, never load again")
		return nil
	}

	var ctx = context.Background()
	resp := utils.CheckErrorWithData(c.API.ListZonesContext(ctx, cf.WithZoneFilters("", "", "")))
	zones := filterT(
		resp.Result,
		c.config.GetStringSlice(configs.CONF_CF_ZONE_INCLUDE),
		c.config.GetStringSlice(configs.CONF_CF_ZONE_EXCLUDE),
		func(z cf.Zone) string {
			return z.Name
		},
	)
	for _, zone := range zones {
		c.Zones[zone.ID] = zone
	}

	acc, info, err := c.API.Accounts(ctx, cf.AccountsListParams{
		PaginationOptions: cf.PaginationOptions{
			PerPage: 10,
		},
	})
	if err != nil {
		return err
	}

	for info.HasMorePages() {
		var accNext []cf.Account
		accNext, info, err = c.API.Accounts(ctx, cf.AccountsListParams{
			PaginationOptions: cf.PaginationOptions{
				Page:    info.Page,
				PerPage: info.PerPage,
			},
		})
		if err != nil {
			return err
		}

		acc = append(acc, accNext...)
	}

	accounts := filterT(
		acc,
		c.config.GetStringSlice(configs.CONF_CF_ACCOUNT_INCLUDE),
		c.config.GetStringSlice(configs.CONF_CF_ACCOUNT_EXCLUDE),
		func(a cf.Account) string {
			return a.Name
		},
	)
	for _, account := range accounts {
		c.Accounts[account.ID] = account
	}

	c.logger.Info(fmt.Sprintf("Set up cloudflare with %d accounts, %d zones", len(c.Accounts), len(c.Zones)))

	if c.config.GetBool(configs.CONF_DEBUG_MODE) {
		accts := make([]string, 0)
		for _, a := range c.Accounts {
			accts = append(accts, a.Name)
		}
		c.logger.Debug(fmt.Sprintf("Cloudflare accounts: %s", strings.Join(accts, ",")))

		zones := make([]string, 0)
		for _, z := range c.Zones {
			zones = append(zones, z.Name)
		}
		c.logger.Debug(fmt.Sprintf("Cloudflare zones: %s", strings.Join(zones, ",")))
	}

	return nil
}

func newAPI(config *viper.Viper) *cf.API {
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

func newGraphQL(config *viper.Viper) graphql.Client {
	client := NewHttpClient(config)
	return graphql.NewClient(API_URL+"/graphql", client)
}
