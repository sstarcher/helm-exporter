package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"k8s.io/helm/pkg/helm"

	"github.com/facebookgo/flagenv"

	"k8s.io/helm/pkg/proto/hapi/release"

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
		"appVersion",
		"updated",
		"namespace",
	})

	client = NewClient()

	tillerNamespace = flag.String("tiller-namespace", "kube-system", "namespace of Tiller (default \"kube-system\")")

	inClusterTiller = fmt.Sprintf("tiller-deploy.%s:44134", *tillerNamespace)
	localTiller     = "127.0.0.1:44134"
	statusCodes     = []release.Status_Code{
		release.Status_UNKNOWN,
		release.Status_DEPLOYED,
		release.Status_DELETED,
		release.Status_DELETING,
		release.Status_FAILED,
		release.Status_PENDING_INSTALL,
		release.Status_PENDING_UPGRADE,
		release.Status_PENDING_ROLLBACK,
	}

	prometheusHandler = promhttp.Handler()
)

// NewClient is the connection to tiller
func NewClient() *helm.Client {
	fmt.Printf("attempting to connect to %s\n", inClusterTiller)
	client := helm.NewClient(helm.Host(inClusterTiller))
	err := client.PingTiller()
	if err != nil {
		fmt.Printf("attempting to connect to %s\n", localTiller)
		client = helm.NewClient(helm.Host(localTiller))
		err := client.PingTiller()
		if err != nil {
			panic(fmt.Sprintf("unable to connect to %s and %s\n", inClusterTiller, localTiller))
		}
		fmt.Printf("connected to %s\n", localTiller)
		return client
	}
	fmt.Printf("connected to %s\n", inClusterTiller)
	return client
}

// Taken from https://github.com/helm/helm/blob/master/cmd/helm/list.go#L197
// filterList returns a list scrubbed of old releases.
func filterList(rels []*release.Release) []*release.Release {
	idx := map[string]int32{}

	for _, r := range rels {
		name, version := r.GetName(), r.GetVersion()
		if max, ok := idx[name]; ok {
			// check if we have a greater version already
			if max > version {
				continue
			}
		}
		idx[name] = version
	}

	uniq := make([]*release.Release, 0, len(idx))
	for _, r := range rels {
		if idx[r.GetName()] == r.GetVersion() {
			uniq = append(uniq, r)
		}
	}
	return uniq
}

func helmStats(w http.ResponseWriter, r *http.Request) {
	items, err := client.ListReleases(helm.ReleaseListStatuses(statusCodes))
	if err == nil {
		stats.Reset()
		for _, item := range filterList(items.GetReleases()) {
			chart := item.GetChart().GetMetadata().GetName()
			status := item.GetInfo().GetStatus().GetCode()
			releaseName := item.GetName()
			version := item.GetChart().GetMetadata().GetVersion()
			appVersion := item.GetChart().GetMetadata().GetAppVersion()
			updated := strconv.FormatInt(item.GetInfo().GetLastDeployed().Seconds, 10)
			namespace := item.GetNamespace()
			if status == release.Status_FAILED {
				status = -1
			}
			stats.WithLabelValues(chart, releaseName, version, appVersion, updated, namespace).Set(float64(status))
		}
	}
	prometheusHandler.ServeHTTP(w, r)
}

func main() {
	flagenv.Parse()
	flag.Parse()

	http.HandleFunc("/metrics", helmStats)
	http.ListenAndServe(":9571", nil)
}
