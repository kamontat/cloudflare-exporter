package prom

import (
	"fmt"

	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/kamontat/cloudflare-exporter/loggers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func New(config *viper.Viper) *Prometheus {
	var reg = prometheus.NewRegistry()
	reg.MustRegister(collectors.NewBuildInfoCollector())
	reg.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	reg.MustRegister(collectors.NewGoCollector())

	whitelist := make(map[string]bool)
	includes := config.GetStringSlice(configs.CONF_METRICS_INCLUDE)
	if len(includes) > 0 {
		for _, v := range includes {
			whitelist[v] = true
		}
	}

	blacklist := make(map[string]bool)
	excludes := config.GetStringSlice(configs.CONF_METRICS_EXCLUDE)
	if len(excludes) > 0 {
		for _, v := range excludes {
			blacklist[v] = true
		}
	}

	return &Prometheus{
		registry:  reg,
		whitelist: whitelist,
		blacklist: blacklist,
		config:    config,
		logger:    loggers.Default(),
	}
}

func MustRegister[C prometheus.Collector](prom *Prometheus, system string, name string, getCollector func(name string, config *viper.Viper) C) (collector C) {
	var isRegis bool
	var checkName = prometheus.BuildFQName("", system, name)

	prom.logger.Debug("Checking metrics name", zap.String("metric-name", checkName))
	if len(prom.whitelist) < 1 && len(prom.blacklist) < 1 {
		isRegis = true
	} else if len(prom.whitelist) > 0 {
		prom.logger.Debug(
			"Found whitelist",
			zap.Strings("whitelist", prom.config.GetStringSlice(configs.CONF_METRICS_INCLUDE)),
		)
		// true, true    => true
		// true, false   => false (impossible)
		// false, true   => false (impossible)
		// false, false  => false
		exist, ok := prom.whitelist[checkName]
		isRegis = exist && ok
	} else if len(prom.blacklist) > 0 {
		prom.logger.Debug(
			"Found blacklist",
			zap.Strings("blacklist", prom.config.GetStringSlice(configs.CONF_METRICS_EXCLUDE)),
		)

		// true, true    => false
		// true, false   => true (impossible)
		// false, true   => true (impossible)
		// false, false  => true
		exist, ok := prom.blacklist[checkName]
		isRegis = !exist && !ok
	} else {
		prom.logger.Debug(
			"Found both whitelist and blacklist",
			zap.Strings("whitelist", prom.config.GetStringSlice(configs.CONF_METRICS_INCLUDE)),
			zap.Strings("blacklist", prom.config.GetStringSlice(configs.CONF_METRICS_EXCLUDE)),
		)

		white, wok := prom.whitelist[checkName]
		black, bok := prom.blacklist[checkName]
		isRegis = !bok && !black && wok && white
	}

	if isRegis {
		prom.logger.Debug(
			fmt.Sprintf("Adding metrics %s to prometheus", checkName),
			zap.String("namespace", NAMESPACE),
		)
		collector = getCollector(prometheus.BuildFQName(NAMESPACE, system, name), prom.config)
		prom.registry.MustRegister(collector)
		return
	}

	return
}
