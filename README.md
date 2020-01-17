# Helm Exporter

[![](https://images.microbadger.com/badges/image/sstarcher/helm-exporter.svg)](http://microbadger.com/images/sstarcher/helm-exporter "Get your own image badge on microbadger.com")
[![Docker Registry](https://img.shields.io/docker/pulls/sstarcher/helm-exporter.svg)](https://registry.hub.docker.com/u/sstarcher/helm-exporter)&nbsp;

Exports helm release, chart, and version statistics in the prometheus format.

# Installation
* A helm chart is available in this [repository](https://hub.helm.sh/charts/sstarcher/helm-exporter).
* `helm install helm-exporter sstarcher/helm-exporter` will install and metrics should scrape automatically if prometheus is running
* If using Grafana you can use this Dashboard to have a list of what's running https://grafana.com/dashboards/9367

# Metrics
* http://host:9571/metrics

# Format
```
helm_chart_info{chart="ark",release="ark",version="1.2.1",appVersion="1.2.3",updated="1553201431",namespace="test"} 1
helm_chart_info{chart="cluster-autoscaler",release="cluster-autoscaler",version="0.7.0",appVersion="",updated="1553201431",namespace="other"} 4
helm_chart_info{chart="dex",release="dex",version="0.1.0",appVersion="1.2.3",updated="1553201431",namespace="test"} 1
```

The metric value is the helm status code.  These status codes indexes do not map up directly to helm.  This is so I can make the bad cases negative values.
* -1 FAILED
* 0 UNKNOWN
* 1 DEPLOYED
* 2 DELETED
* 3 SUPERSEDED
* 5 DELETING
* 6 PENDING_INSTALL
* 7 PENDING_UPGRADE
* 8 PENDING_ROLLBACK

# Prior Art
* https://github.com/Kubedex/exporter

# Todo
* /healthz endpoint method

