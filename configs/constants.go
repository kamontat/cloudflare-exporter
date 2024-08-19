package configs

const (
	CONF_PRODUCTION            = "production"
	CONF_DEBUG_MODE            = "debug"
	CONF_SILENT_MODE           = "silent"
	CONF_OUTPUT_JSON           = "json"
	CONF_SERVER_ADDR           = "server.addr"
	CONF_SERVER_PORT           = "server.port"
	CONF_SERVER_METRIC_PATH    = "server.metric-path"
	CONF_SERVER_HEALTH_PATH    = "server.health-path"
	CONF_SERVER_READINESS_PATH = "server.readiness-path"
	CONF_SERVER_LIVENESS_PATH  = "server.liveness-path"
	CONF_SERVER_CONCURRENCY    = "server.concurrency"
	CONF_SERVER_BODY_LIMIT     = "server.body-limit"
	CONF_SERVER_READ_BUFFER    = "server.read-buffer"
	CONF_SERVER_WRITE_BUFFER   = "server.write-buffer"
	CONF_SERVER_READ_TIMEOUT   = "server.read-timeout"
	CONF_SERVER_WRITE_TIMEOUT  = "server.write-timeout"
	CONF_SERVER_IDLE_TIMEOUT   = "server.idle-timeout"
	CONF_CF_API_TOKEN          = "cf.api-token"
	CONF_CF_API_EMAIL          = "cf.api-email"
	CONF_CF_API_KEY            = "cf.api-key"
	CONF_CF_ERROR_MODE         = "cf.error-mode"
	CONF_CF_INTERVAL           = "cf.interval"
	CONF_CF_TIMEOUT            = "cf.timeout"
	CONF_CF_ACCOUNT_INCLUDE    = "cf.account.include"
	CONF_CF_ACCOUNT_EXCLUDE    = "cf.account.exclude"
	CONF_CF_ZONE_INCLUDE       = "cf.zone.include"
	CONF_CF_ZONE_EXCLUDE       = "cf.zone.exclude"
	CONF_METRICS_INCLUDE       = "metrics.include"
	CONF_METRICS_EXCLUDE       = "metrics.exclude"

	SETTING_CONFIG_FILE = "config"
	SETTING_CONFIG_EXT  = "yaml"
)

const (
	ERROR_MODE_LOG  = "log"
	ERROR_MODE_STOP = "stop"
)
