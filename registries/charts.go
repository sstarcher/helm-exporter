package registries

import (
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/sstarcher/helm-exporter/versioning"
)

// HelmRegistries contains all the information regarding helm registries
type HelmRegistries struct {
	OverrideChartNames map[string]string      `koanf:"overrideChartNames"`
	OverrideRegistries []HelmOverrideRegistry `koanf:"override"`
	RegistryNames      []string               `koanf:"registryNames"`
}

// HelmOverrideRegistry contains information about which registry to use to fetch helm versions
type HelmOverrideRegistry struct {
	HelmRegistry     HelmRegistry `koanf:"registry"`
	Charts           []string     `koanf:"charts"`
	AllowAllReleases bool         `koanf:"allowAllReleases"`
}

// HelmRegistry contains information about the helm registry
type HelmRegistry struct {
	URL string `koanf:"url"`
}

// GetLatestVersionFromHelm fetches the latest version of the helm chart
func (h HelmRegistries) GetLatestVersionFromHelm(chart string) string {
	if val, ok := h.OverrideChartNames[chart]; ok {
		chart = val
	}

	log.WithField("chart", chart).Debug("Fetching version for chart")

	for _, registry := range h.OverrideRegistries {
		for _, chartOverride := range registry.Charts {
			match, err := regexp.MatchString(chartOverride, chart)
			if err != nil {
				log.WithError(err).Fatal("Chart regexp not valid")
			}
			if match {
				return registry.getChartVersions(chart)
			}
		}
	}

	return h.fromArtifactHub(chart)
}

func (h HelmRegistries) fromArtifactHub(chart string) string {
	logger := log.WithField("chart", chart)

	charts := []hubChart{}
	for _, val := range hubCache {
		if val.Name == chart {
			if len(h.RegistryNames) == 0 {
				charts = append(charts, val)
			} else {
				for _, reg := range h.RegistryNames {
					if val.Repository.Name == reg {
						charts = append(charts, val)
					}
				}
			}
		}
	}

	if len(charts) == 0 {
		logger.Errorf("unable to find any charts for %s", chart)
		return versioning.Failure
	} else if len(charts) > 1 {
		regs := []string{}
		for _, val := range charts {
			regs = append(regs, val.Repository.Name)
		}
		logger.Errorf("Failed to search chart info, found multiple registries that contain this %s [%s].", chart, strings.Join(regs, ", "))
		return versioning.Multiple
	}

	return charts[0].Version
}
