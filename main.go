package main

import (
	"time"

	"net/http"

	"k8s.io/helm/pkg/helm"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	stats = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "helm_chart_info",
		Help: "Information on helm releases",
	}, []string{
		"chart",
		"release",
		"version",
	})

	client = NewClient()
)

func HelmStats() {
	items, err := client.ListReleases()
	if err == nil {
		for _, item := range items.GetReleases() {
			stats.WithLabelValues(item.GetChart().GetMetadata().GetName(), item.GetName(), item.GetChart().GetMetadata().GetVersion()).Set(1)
		}
	}
}

func NewClient() *helm.Client {
	client := helm.NewClient(helm.Host("tiller-deploy.kube-system:44134"))
	err := client.PingTiller()
	if err != nil {
		client = helm.NewClient(helm.Host("127.0.0.1:44134"))
		err := client.PingTiller()
		if err != nil {
			panic("unable to connect to 127.0.0.1:44134 and tiller-deploy.kube-system:44134")
		}
	}
	return client
}

func main() {
	go func() {
		for {
			HelmStats()
			time.Sleep(2 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
