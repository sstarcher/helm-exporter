package registries

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrMultipleCharts multiple charts found
	ErrMultipleCharts = errors.New("multiple charts found")
	// ErrNoChartsFound no charts found for the name
	ErrNoChartsFound = fmt.Errorf("Could not find the chart")

	hubCache = artifacthubDump{}
	mu       sync.Mutex
)

const (
	baseURL   = "https://artifacthub.io"
	userAgent = "helm-exporter/1.0"
)

type artifacthubDump []hubChart

type hubChart struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	Repository struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repository"`
}

func httpGet(endpoint string, data interface{}) error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseURL, endpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(data)
}

func init() {
	err := update()
	if err != nil {
		log.Warnf("failed to update artifacthub cache due to %v ", err)
	}

	go func() {
		for {
			time.Sleep(time.Hour)
			log.Info("updating artifacthub dump")
			err := update()
			if err != nil {
				log.Warnf("failed to update artifacthub cache due to %v ", err)
			}
		}
	}()
}

func update() error {
	data := &artifacthubDump{}
	err := httpGet("api/v1/helm-exporter", data)
	if err != nil {
		return err
	}

	mu.Lock()
	hubCache = *data
	mu.Unlock()
	return nil
}
