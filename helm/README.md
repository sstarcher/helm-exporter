# helm-exporter

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

Exporter for helm metrics

**Homepage:** <https://github.com/sstarcher/helm-exporter>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| sstarcher | shane.starcher@gmail.com | https://shanestarcher.com |

## Source Code

* <https://github.com/sstarcher/helm-exporter>

## Values

| Key                                                | Type   | Default                     | Description                                                    |
|----------------------------------------------------|--------|-----------------------------|----------------------------------------------------------------|
| affinity                                           | object | `{}`                        |                                                                |
| config.helmRegistries.overrideChartNames           | object | `{}`                        |                                                                |
| config.helmRegistries.override[0].allowAllReleases | bool   | `true`                      |                                                                |
| config.helmRegistries.override[0].charts           | list   | `[]`                        |                                                                |
| config.helmRegistries.override[0].registry.url     | string | `""`                        |                                                                |
| config.helmRegistries.registryNames                | list   | `[]`                        |                                                                |
| env                                                | list   | `[]`                        |                                                                |
| fullnameOverride                                   | string | `""`                        |                                                                |
| image.pullPolicy                                   | string | `"Always"`                  |                                                                |
| image.repository                                   | string | `"sstarcher/helm-exporter"` |                                                                |
| image.tag                                          | string | `""`                        |                                                                |
| imagePullSecrets                                   | list   | `[]`                        |                                                                |
| infoMetric                                         | bool   | `true`                      |                                                                |
| ingress.annotations                                | object | `{}`                        |                                                                |
| ingress.enabled                                    | bool   | `false`                     |                                                                |
| ingress.hosts[0].host                              | string | `"chart-example.local"`     |                                                                |
| ingress.hosts[0].paths                             | list   | `[]`                        |                                                                |
| ingress.tls                                        | list   | `[]`                        |                                                                |
| intervalDuration                                   | int    | `0`                         |                                                                |
| latestChartVersion                                 | bool   | `true`                      |                                                                |
| livenessProbe                                      | object | (see `values.yaml`)         |  Liveness probe configuration                                  |
| nameOverride                                       | string | `""`                        |                                                                |
| namespaces                                         | string | `""`                        |                                                                |
| nodeSelector                                       | object | `{}`                        |                                                                |
| podAnnotations                                     | object | `{}`                        |                                                                |
| podLabels                                          | object | `{}`                        |                                                                |
| podSecurityContext                                 | object | `{}`                        |                                                                |
| rbac.create                                        | bool   | `true`                      |                                                                |
| readinessProbe                                     | object | (see `values.yaml`)         |  Readiness probe configuration                                 |
| replicaCount                                       | int    | `1`                         |                                                                |
| resources                                          | object | `{}`                        |                                                                |
| securityContext                                    | object | `{}`                        |                                                                |
| service.annotations                                | object | `{}`                        |                                                                |
| service.port                                       | int    | `9571`                      |                                                                |
| service.type                                       | string | `"ClusterIP"`               |                                                                |
| serviceAccount.create                              | bool   | `true`                      |                                                                |
| serviceAccount.name                                | string | `nil`                       |                                                                |
| serviceMonitor.additionalLabels                    | object | `{}`                        |                                                                |
| serviceMonitor.create                              | bool   | `false`                     |                                                                |
| serviceMonitor.interval                            | string | `nil`                       |                                                                |
| serviceMonitor.namespace                           | string | `nil`                       |                                                                |
| serviceMonitor.scrapeTimeout                       | string | `nil`                       |                                                                |
| startupProbe                                       | object | (see `values.yaml`)         |  Startup probe configuration                                   |
| timestampMetric                                    | bool   | `true`                      |                                                                |
| tolerations                                        | list   | `[]`                        |                                                                |
| grafanaDashboard.enabled                           | bool   | `false`                     | Specifies whether a Grafana dashboard should be created        |
| grafanaDashboard.namespace                         | bool   | `nil`                       | Specifies then namespace where the dashboard should be created |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.5.0](https://github.com/norwoodj/helm-docs/releases/v1.5.0)
