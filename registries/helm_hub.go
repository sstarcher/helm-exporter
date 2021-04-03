package registries

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/sstarcher/helm-exporter/versioning"
	log "github.com/sirupsen/logrus"
	"errors"
)

var errMultipleCharts = errors.New("multiple charts found")

// Chart contains attribute information for a chart
type Chart struct {
	AvailableVersions []AvailableVersions `json:"available_versions"`
}

// AvailableVersions contains list of versions for a chart
type AvailableVersions struct {
	Version string `json:"version"`
}

// SearchResultData contains search results from hub.helm.sh
type SearchResultData struct {
	Data []ChartSearchResult `json:"data"`
}

// ChartSearchResult contains chart search results from hub.helm.sh
type ChartSearchResult struct {
	ID string `json:"id"`
}

func (h HelmRegistries) useHelmHub(chart string) string {
	chartName := h.OverrideChartNames[chart]
	if chartName == "" {
		var err error
		chartName, err = findChart(chart)
		if err != nil {
			if err == errMultipleCharts {
				log.WithError(err).WithField("chart", chart).Error("Failed to search chart info, found multiple charts.")
				return versioning.Multiple
			} else {
			log.WithError(err).WithField("chart", chart).Error("Failed to search chart info")
			return versioning.Failure
			}
		}
	}

	versions, err := getChartVersions(chartName)
	if err != nil {
		log.WithError(err).WithField("chart", chart).Error("Failed to fetch chart info")
		return versioning.Failure
	}

	return versioning.FindHighestVersionInList(versions, false)
}

func findChart(chart string) (string, error) {
	url := fmt.Sprintf("https://hub.helm.sh/api/chartsvc/v1/charts/search?q=%s", chart)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	searchData := SearchResultData{}
	err = json.NewDecoder(resp.Body).Decode(&searchData)
	if err != nil {
		return "", err
	}

	if len(searchData.Data) == 0 {
		return "", fmt.Errorf("Could not find the chart")
	} else if len(searchData.Data) == 1 {
		return searchData.Data[0].ID, nil
	}
	return "", errMultipleCharts
}

func getChartVersions(chart string) ([]string, error) {
	url := fmt.Sprintf("https://artifacthub.io/api/v1/packages/helm/%s", chart)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	chartData := Chart{}
	err = json.NewDecoder(resp.Body).Decode(&chartData)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, data := range chartData.AvailableVersions {
		versions = append(versions, data.Version)
	}
	return versions, nil
}
