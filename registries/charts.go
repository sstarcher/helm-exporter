package registries

import (
	"regexp"

	log "github.com/sirupsen/logrus"
)

// HelmRegistries contains all the information regarding helm registries
type HelmRegistries struct {
	OverrideChartNames map[string]string      `koanf:"overrideChartNames"`
	OverrideRegistries []HelmOverrideRegistry `koanf:"override"`
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
	log.WithField("chart", chart).Debug("Fetching version for chart")

	for _, registry := range h.OverrideRegistries {
		for _, chartOverride := range registry.Charts {
			match, err := regexp.MatchString(chartOverride, chart)
			if err != nil {
				log.WithError(err).Fatal("Chart regexp not valid")
			}
			if match {
				chartName := h.OverrideChartNames[chart]
				if chartName == "" {
					chartName = chart
				}
				return registry.getChartVersions(chartName)
			}
		}
	}

	return h.useHelmHub(chart)
}
