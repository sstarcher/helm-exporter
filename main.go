package main

import (
	"flag"
	"net/http"
	"strconv"
	"strings"

	cmap "github.com/orcaman/concurrent-map"

	log "github.com/sirupsen/logrus"

	"os"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/facebookgo/flagenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	settings = cli.New()
	clients  = cmap.New()

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

	namespaces = flag.String("namespaces", "", "namespaces to monitor.  Defaults to all")

	statusCodeMap = map[string]float64{
		"unknown":          0,
		"deployed":         1,
		"uninstalled":      2,
		"superseded":       3,
		"failed":           -1,
		"uninstalling":     5,
		"pending-install":  6,
		"pending-upgrade":  7,
		"pending-rollback": 8,
	}

	prometheusHandler = promhttp.Handler()
)

func runStats() {
	stats.Reset()
	for _, client := range clients.Items() {
		list := action.NewList(client.(*action.Configuration))
		items, err := list.Run()
		if err != nil {
			log.Warnf("got error while listing %v", err)
			continue
		}

		for _, item := range items {
			chart := item.Chart.Name()
			releaseName := item.Name
			version := item.Chart.Metadata.Version
			appVersion := item.Chart.AppVersion()
			updated := strconv.FormatInt((item.Info.LastDeployed.Unix() * 1000), 10)
			namespace := item.Namespace
			status := statusCodeMap[item.Info.Status.String()]
			stats.WithLabelValues(chart, releaseName, version, appVersion, updated, namespace).Set(status)
		}
	}
}

func newHelmStatsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runStats()
		prometheusHandler.ServeHTTP(w, r)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {

}

func connect(namespace string) {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), log.Infof)
	if err != nil {
		log.Warnf("failed to connect to %s with %v", namespace, err)
	} else {
		log.Infof("Watching namespace %s", namespace)
		clients.Set(namespace, actionConfig)
	}
}

func informer() {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Infof)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := actionConfig.KubernetesClientSet()
	if err != nil {
		log.Fatal(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, 0)
	informer := factory.Core().V1().Namespaces().Informer()
	stopper := make(chan struct{})
	defer close(stopper)

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(v1.Object)
			connect(mObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Infof("Removing namespace %s", mObj.GetName())
			clients.Remove(mObj.GetName())
		},
	})

	informer.Run(stopper)
}

func main() {
	flagenv.Parse()
	flag.Parse()

	if namespaces == nil || *namespaces == "" {
		go informer()
	} else {
		for _, namespace := range strings.Split(*namespaces, ",") {
			connect(namespace)
		}
	}

	http.HandleFunc("/metrics", newHelmStatsHandler())
	http.HandleFunc("/healthz", healthz)
	http.ListenAndServe(":9571", nil)
}
