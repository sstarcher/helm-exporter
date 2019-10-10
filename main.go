package main

import (
	"flag"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/release"
	"k8s.io/helm/pkg/tlsutil"

	"github.com/facebookgo/flagenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	clients []*helm.Client
	mutex   sync.RWMutex

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

	tillers         = flag.String("tillers", "tiller-deploy.kube-system:44134", "tiller address HOST:PORT of Tillers, separated list tiller-deploy.kube-system44134,tiller-deploy.dev44134")
	tillerTLSEnable = flag.Bool("tiller-tls-enable", false, "enable TLS communication with tiller (default false)")
	tillerTLSKey    = flag.String("tiller-tls-key", "/etc/helm-exporter/tls.key", "path to private key file used to communicate with tiller")
	tillerTLSCert   = flag.String("tiller-tls-cert", "/etc/helm-exporter/tls.crt", "path to certificate key file used to communicate with tiller")
	tillerTLSVerify = flag.Bool("tiller-tls-verify", false, "enable verification of the remote tiller certificate (default false)")

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
	options := []helm.Option{helm.Host(tillerEndpoint)}
	if *tillerTLSEnable {
		tlsopts := tlsutil.Options{
			KeyFile:            *tillerTLSKey,
			CertFile:           *tillerTLSCert,
			InsecureSkipVerify: !(*tillerTLSVerify),
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

func newHelmStatsHandler() http.HandlerFunc {
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

func connect(tiller string) {
	for {
		client, err := newHelmClient(tiller)
		if err != nil {
			log.Warnf("failed to connect to %s with %v", tiller, err)
		} else {
			mutex.Lock()
			clients = append(clients, client)
			log.Infof("connected to %s", tiller)
			mutex.Unlock()
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func main() {
	flagenv.Parse()
	flag.Parse()

	for _, tiller := range strings.Split(*tillers, ",") {
		go connect(tiller)
	}

	http.HandleFunc("/metrics", newHelmStatsHandler())
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":9571", nil)
}
