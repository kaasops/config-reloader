<p align="center">
  <a href="https://github.com/vlzemtsov/config-reloader/">
    <img src="https://img.shields.io/badge/license-Apache%20v2-green.svg" alt="license">
  </a>
</p>

# Kubernetes Config (ConfigMap and Secret) Reloader

This progect - based on https://github.com/jimmidyson/configmap-reload and https://github.com/prometheus-operator/prometheus-operator/pkgs/container/prometheus-config-reloader


**config-reloader** is a simple binary to trigger a reload when Kubernetes ConfigMaps or Secrets are updated.
It watches mounted volume dirs and notifies the target process changed files on dirs.
If changes exist - send webhook.

## Features
- Send webook if files in dirs changed (if configmap or secret have been changed)
- Control many dirs
- Unarchive .gz archive to file and update file, if .gz has been changed
- Init mode (stop after unarchive) 
- Prometheus metrics


It is available as a Docker image at https://hub.docker.com/r/vlzemtsov/config-reloader

### Usage

```
Usage of ./out/config-reloader:
  -volume-dir value
        the config map volume directory to watch for updates; may be used multiple times
  -web.listen-address string
    	  address to listen on for web interface and telemetry. (default ":9533")
  -web.telemetry-path string
    	  path under which to expose metrics. (default "/metrics")
  -webhook-method string
        the HTTP method url to use to send the webhook (default "POST")
  -webhook-status-code int
        the HTTP status code indicating successful triggering of reload (default 200)
  -webhook-url string
        the url to send a request to when the specified config map volume directory has been updated
  -webhook-retries integer
        the amount of times to retry the webhook reload request
```
