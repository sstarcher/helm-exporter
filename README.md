# Helm Exporter

[![](https://images.microbadger.com/badges/image/sstarcher/helm-exporter.svg)](http://microbadger.com/images/sstarcher/helm-exporter "Get your own image badge on microbadger.com")
[![Docker Registry](https://img.shields.io/docker/pulls/sstarcher/helm-exporter.svg)](https://registry.hub.docker.com/u/sstarcher/helm-exporter)&nbsp;

Exports helm release, chart, and version statistics in the prometheus format.

# Installation
## Prerequisites

- Kubernetes 1.19+
- Helm 3+

## Get Repo Info

```console
helm repo add sstarcher https://shanestarcher.com/helm-charts/
helm repo update
```

_See [helm repo](https://hub.helm.sh/charts/sstarcher/helm-exporter) for command documentation._

## Install Chart

```console
# Helm
$ helm install helm-exporter sstarcher/helm-exporter
```
* `helm install helm-exporter sstarcher/helm-exporter` will install and metrics should scrape automatically if prometheus is running
* If using Grafana you can use this Dashboard to have a list of what's running https://grafana.com/dashboards/9367

## Configuration for Latest versions

Two options exist for fetching the latest version information for a chart.
* Direct fetch from a chart repository.  This will download the index for the registry and use that information to fetch the chart.

```yaml
# Helm configuration
config:
  helmRegistries:
    overrideChartNames: {}
      mysql: stable/test
# If the helm charts are not stored on hub.helm.sh then a custom registry can be configured here.
# Currently only index.yaml registry is supported (helm supports other registries as well)
    override:
      - registry:
          url: "https://some.url" # Url to the index file
        charts: # Chart names
        - splunk
        - falco-eks-audit-bridge
```

## Configuration for password protected registries

If the registry needs authentication then you can use a Kubernetes secret to store the username and password.

```bash
kubectl create secret generic chartmuseum --from-literal=username=admin --from-literal=password=admin
```

And use following configuration:

```yaml
# Helm configuration
config:
  helmRegistries:
    overrideChartNames: {}
      mysql: stable/test
    # If the helm charts are not stored on hub.helm.sh then a custom registry can be configured here.
    # Currently only index.yaml registry is supported (helm supports other registries as well)
    override:
      - registry:
          url: "https://some.url" # Url to the index file
          secretRef:
              name: "chartmuseum" # Name of the secret containing the username and password
              userKey: "username" # Key of the username in the secret
              passKey: "password" # Key of the password in the secret
        charts: # Chart names
          - splunk
          - falco-eks-audit-bridge
```


* Query https://artifacthub.io for the chart matching your chart name and only using the specified registries.  If no registry name is specified and multiple charts match from helm hub no version will be found and it will log a warning.
```yaml
# Helm configuration
config:
  helmRegistries:
    registryNames:
    - bitnami
```

# Metrics
* http://host:9571/metrics

# Format
```
# HELP helm_chart_info Information on helm releases
# TYPE helm_chart_info gauge
helm_chart_info{chart="ark",release="ark",version="1.2.1",latestVersion="1.2.3",appVersion="1.2.3",updated="1553201431",namespace="test"} 1
helm_chart_info{chart="cluster-autoscaler",release="cluster-autoscaler",version="0.7.0",latestVersion=7.3.2,appVersion="",updated="1553201431",namespace="other"} 4
helm_chart_info{chart="dex",release="dex",version="0.1.0",latestVersion="3.4.0",appVersion="1.2.3",updated="1553201431",namespace="test"} 1

# HELP helm_chart_outdated Outdated helm versions of helm releases
# TYPE helm_chart_outdated gauge
helm_chart_outdated{chart="ark",latestVersion="1.2.3",namespace="test",release="ark",version="1.2.1"} 1
helm_chart_outdated{chart="cluster-autoscaler",latestVersion="7.3.2",namespace="other",release="cluster-autoscaler",version="0.7.0"} 1
helm_chart_outdated{chart="external-secrets",latestVersion="3.4.0",namespace="test",release="dex",version="0.1.0"} 1

# HELP helm_chart_timestamp Timestamps of helm releases
# TYPE helm_chart_timestamp gauge
helm_chart_timestamp{chart="ark",release="ark",version="1.2.1",latestVersion="1.2.3",appVersion="1.2.3",updated="1553201431",namespace="test"} 1.617197959e+12
helm_chart_timestamp{chart="cluster-autoscaler",release="cluster-autoscaler",version="0.7.0",latestVersion=7.3.2,appVersion="",updated="1553201431",namespace="other"} 1.617196128e+12
helm_chart_timestamp{chart="dex",release="dex",version="0.1.0",latestVersion="3.4.0",appVersion="1.2.3",updated="1553201431",namespace="test"} 1.62245881e+12

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

