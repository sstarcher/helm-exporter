package registries

import (
	"net/http"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
	"github.com/sstarcher/helm-exporter/versioning"
)

// IndexEntries contains configured Helm indexes
type IndexEntries struct {
	Entries map[string][]IndexEntry `yaml:"entries"`
}

// IndexEntry the actual Helm index information
type IndexEntry struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

func (r HelmOverrideRegistry) getChartVersions(chart string) string {
	resp, err := http.Get(r.HelmRegistry.URL)
	if err != nil {
		log.WithError(err).WithField("chart", chart).WithField("registry", r.HelmRegistry.URL).Error("Failed to get chart info")
		return versioning.Failure
	}
	defer resp.Body.Close()

	index := IndexEntries{}
	err = yaml.NewDecoder(resp.Body).Decode(&index)
	if err != nil {
		log.WithError(err).WithField("chart", chart).WithField("registry", r.HelmRegistry.URL).Error("Failed to unmarshal chart info")
		return versioning.Failure
	}

	var versions []string
	entries := index.Entries[chart]
	if entries == nil {
		return versioning.Notfound
	}
	for _, entry := range entries {
		versions = append(versions, entry.Version)
	}

	return versioning.FindHighestVersionInList(versions, r.AllowAllReleases)
}
