# helm-exporter

Installs [helm-exporter](https://github.com/sstarcher/helm-exporter) to export helm stats to prometheus.

## TL;DR;

```console
$ helm install sstarcher/helm-exporter
```

## Introduction

This chart bootstraps a [helm-exporter](https://github.com/sstarcher/helm-exporter) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

The chart comes with a ServiceMonitor for use with the [Prometheus Operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator).

## Installing the Chart

To install the chart with the release name `my-release`:

```console
$ helm install sstarcher/helm-exporter --name my-release
```

The command deploys helm-exporter on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the helm-exporter chart that are in addition to values in a default helm chart.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Assign custom affinity rules for helm-exporter [https://kubernetes.io/docs/concepts/configuration/assign-pod-node/](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/) |
| config.helmRegistries.overrideChartNames | object | `{}` | Provide a name to substitute for the full names of resources e.g. `mysql: stable/mysql` |
| config.helmRegistries.override[0].allowAllReleases | bool | `true` | This allows all semver versions, like release candidates or custom suffixes. |
| config.helmRegistries.override[0].charts | list | `[]` | Chart names for the override (chart) registry/repo url |
| config.helmRegistries.override[0].registry.url | string | `""` |  Url to the index file for a custom helm repo |
| fullnameOverride | string | `""` | Provide a name to substitute for the full names of resources |
| image.pullPolicy | string | `"Always"` | Image pull policy for the webhook integration jobs |
| image.repository | string | `"sstarcher/helm-exporter"` | Repository to use for the webhook integration jobs |
| imagePullSecrets | list | `[]` | Reference to one or more secrets to be used when pulling images |
| infoMetric | bool | `true` | Specifies whether to generate the info metric. |
| ingress.annotations | object | `{}` |  Annotations for the helm-exporter |
| ingress.enabled | bool | `false` | If true, helm-exporter Ingress will be created |
| ingress.hosts[0].host | string | `"chart-example.local"` | Ingress hostname |
| ingress.hosts[0].paths | list | `[]` | Ingress paths |
| ingress.tls | list | `[]` | Ingress TLS configuration (YAML) |
| latestChartVersion | bool | `true` | Specifies whether to fetch the latest chart versions from repositories. |
| nameOverride | string | `""` | Provide a name in place of helm-exporter |
| namespaces | string | `""` | Specifies which namespaces to query for helm 3 metrics.  Defaults to all |
| nodeSelector | object | `{}` | helm-exporter node selector [https://kubernetes.io/docs/user-guide/node-selection/](https://kubernetes.io/docs/user-guide/node-selection/ ) |
| podAnnotations | object | `{}` | Annotations to add to the pod |
| podSecurityContext | object | `{}` | SecurityContext for helm-exporter pod |
| rbac.create | bool | `true` | Create RBAC resources |
| replicaCount | int | `1` | Number of instances to deploy. |
| resources | object | `{}` | Define resources requests and limits for single Pods. |
| securityContext | object | `{}` | SecurityContext for a container |
| service.port | int | `9571` | Port for Service to listen on. |
| service.type | string | `"ClusterIP"` | Service type |
| serviceAccount.create | bool | `true` | Create a default serviceaccount to use |
| serviceAccount.name | string | `default` | Name for prometheus serviceaccount |
| serviceMonitor.create | bool | `false` | Set to true if using the Prometheus Operator |
| serviceMonitor.interval | string | `nil` | Interval at which metrics should be scraped |
| serviceMonitor.namespace | string | `nil` | The namespace where the Prometheus Operator is deployed |
| serviceMonitor.additionalLabels |object | `{}` | Additional labels to add to the ServiceMonitor	|
| serviceMonitor.scrapeTimeout | string | `nil` | Scrape Timeout when the metrics endpoint is scraped |
| timestampMetric | bool | `true` | Specifies whether to generate the timestamps metric. |
| tolerations | list | `[]` | Tolerations for use with node taints [https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)|


```console
$ helm install my-release sstarcher/helm-exporter
```

Alternatively, a YAML file that specifies the values for the above parameters can be provided while installing the chart. For example,

```console
$ helm install my-release sstarcher/helm-exporter -f values.yaml
```

> **Tip**: You can use the default [values.yaml](values.yaml)
