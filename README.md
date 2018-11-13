# Helm Exporter

[![](https://images.microbadger.com/badges/image/sstarcher/helm_exporter.svg)](http://microbadger.com/images/sstarcher/helm_exporter "Get your own image badge on microbadger.com")
[![Docker Registry](https://img.shields.io/docker/pulls/sstarcher/helm_exporter.svg)](https://registry.hub.docker.com/u/sstarcher/helm_exporter)&nbsp;

Exports helm release, chart, and version staistics in the prometheus format.

# Installation
* A helm chart is available in this [repository](./helm/helm_exporter).
* `helm install -f helm/helm_exporter` will install and metrics should scrape automatically if prometheus is running

# Metrics
* http://host:9100/metrics

# Format
```
helm_chart_info{chart="ark",release="ark",version="1.2.1"} 1
helm_chart_info{chart="cluster-autoscaler",release="cluster-autoscaler",version="0.7.0"} 1
helm_chart_info{chart="dex",release="dex",version="0.1.0"} 1
```

# Prior Art
* https://github.com/Kubedex/exporter
