# Cloudflare Exporter

Cloudflare metric exporter for Prometheus

## Usage

You can either use binary file in [Release](https://github.com/kamontat/cloudflare-exporter/releases/latest) or [Docker image](https://github.com/users/kamontat/packages/container/package/cloudflare-exporter).
You can setting using either config files, commandline options, or environment variables.

Command will find config file from below locations:
- **/etc/cf-exporter/config.yaml**
- **$HOME/.config/cf-exporter/config.yaml**
- **$PWD/config.yaml**

### Debug mode

> `debug: true` in config.yaml, `--debug`, or `DEBUG=true` variable.

Enabled debug mode.

### Silent mode

> `silent: true` in config.yaml, `--silent`, or `SILENT=true` variable.

Enabled silent mode (ignored when debug mode is enabled).

### JSON mode

> `json: true` in config.yaml, `--json`, or `CFE_JSON=true` variable.

Use JSON format to print log message.
