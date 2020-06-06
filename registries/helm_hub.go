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

// Charts is data structure coming from hub.helm.sh
type Charts struct {
	Data []Chart `json:"data"`
}

// Chart contains attribute information for a chart coming from hub.helm.sh
type Chart struct {
	Attributes Attributes `json:"attributes"`
}

// Attributes contains version information for a chart coming from hub.helm.sh
type Attributes struct {
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
	url := fmt.Sprintf("https://hub.helm.sh/api/chartsvc/v1/charts/%s/versions", chart)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	chartsData := Charts{}
	err = json.NewDecoder(resp.Body).Decode(&chartsData)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, data := range chartsData.Data {
		versions = append(versions, data.Attributes.Version)
	}
	return versions, nil
}
