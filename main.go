package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/tlsutil"

	"github.com/facebookgo/flagenv"
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

	localTiller     = "127.0.0.1:44134"
	tillerNamespace = flag.String("tiller-namespaces", "kube-system", "namespaces of Tillers , separated list kube-system,dev")
	tillerTLSEnable = flag.Bool("tiller-tls-enable", false, "enable TLS communication with tiller (default false)")
	tillerTLSKey    = flag.String("tiller-tls-key-path", "/etc/helm-exporter/tls.key", "path to private key file used to communicate with tiller")
	tillerTLSCert   = flag.String("tiller-tls-cert-path", "/etc/helm-exporter/tls.crt", "path to certificate key file used to communicate with tiller")

	statusCodes = []release.Status_Code{
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

// newHelmClient creates a Helm client to the given Tiller. Tries to
// ping Tiller and returns an error if this fails.
func newHelmClient(tillerEndpoint string) (*helm.Client, error) {
	log.Printf("Attempting to connect to %s", tillerEndpoint)

	options := []helm.Option{helm.Host(tillerEndpoint)}
	if *tillerTLSEnable {
		tlsopts := tlsutil.Options{
			KeyFile:            *tillerTLSKey,
			CertFile:           *tillerTLSCert,
			InsecureSkipVerify: true,
		}
		tlscfg, err := tlsutil.ClientConfig(tlsopts)
		if err != nil {
			return nil, err
		}
		options = append(options, helm.WithTLS(tlscfg))
	}

	client := helm.NewClient(options...)
	err := client.PingTiller()
	return client, err
}

// filterList returns a list scrubbed of old releases.
// Taken from https://github.com/helm/helm/blob/master/cmd/helm/list.go#L197
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

func newHelmStatsHandler(clients []*helm.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats.Reset()
		for _, client := range clients {
			items, err := client.ListReleases(helm.ReleaseListStatuses(statusCodes))
			if err == nil {
				for _, item := range filterList(items.GetReleases()) {
					chart := item.GetChart().GetMetadata().GetName()
					status := item.GetInfo().GetStatus().GetCode()
					releaseName := item.GetName()
					version := item.GetChart().GetMetadata().GetVersion()
					appVersion := item.GetChart().GetMetadata().GetAppVersion()
					updated := strconv.FormatInt((item.GetInfo().GetLastDeployed().Seconds * 1000), 10)
					namespace := item.GetNamespace()
					if status == release.Status_FAILED {
						status = -1
					}
					stats.WithLabelValues(chart, releaseName, version, appVersion, updated, namespace).Set(float64(status))
				}
			}
		}
		prometheusHandler.ServeHTTP(w, r)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {

}

func main() {
	flagenv.Parse()
	flag.Parse()

	var clients []*helm.Client
	tNamespaces := strings.Split(*tillerNamespace, ",")

	for _, np := range tNamespaces {
		client, err := newHelmClient(fmt.Sprintf("tiller-deploy.%s:44134", np))
		if err != nil {
			log.Printf("Failed to connect: %v", err)

			client, err = newHelmClient(localTiller)
			if err != nil {
				log.Printf("Failed to connect: %v", err)
				log.Fatalln("Giving up.")
			}
		}
		clients = append(clients, client)
	}

	http.HandleFunc("/metrics", newHelmStatsHandler(clients))
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":9571", nil)
}
