package main

import (
	"flag"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sstarcher/helm-exporter/config"
	"github.com/sstarcher/helm-exporter/registries"

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
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	settings = cli.New()
	clients  = cmap.New()

	mutex = sync.RWMutex{}

	statsInfo      *prometheus.GaugeVec
	statsTimestamp *prometheus.GaugeVec

	namespaces = flag.String("namespaces", "", "namespaces to monitor.  Defaults to all")
	configFile = flag.String("config", "", "Configfile to load for helm overwrite registries.  Default is empty")

	intervalDuration = flag.String("interval-duration", "0", "Enable metrics gathering in background, each given duration. If not provided, the helm stats are computed synchronously.  Default is 0")

	infoMetric      = flag.Bool("info-metric", true, "Generate info metric.  Defaults to true")
	timestampMetric = flag.Bool("timestamp-metric", true, "Generate timestamps metric.  Defaults to true")

	fetchLatest = flag.Bool("latest-chart-version", true, "Attempt to fetch the latest chart version from registries. Defaults to true")

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

func initFlags() config.AppConfig {
	cliFlags := new(config.AppConfig)
	cliFlags.ConfigFile = *configFile
	return *cliFlags
}

func configureMetrics() (info *prometheus.GaugeVec, timestamp *prometheus.GaugeVec) {
	if *infoMetric == true {
		info = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "helm_chart_info",
			Help: "Information on helm releases",
		}, []string{
			"chart",
			"release",
			"version",
			"appVersion",
			"updated",
			"namespace",
			"latestVersion",
		})
	}

	if *timestampMetric == true {
		timestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "helm_chart_timestamp",
			Help: "Timestamps of helm releases",
		}, []string{
			"chart",
			"release",
			"version",
			"appVersion",
			"updated",
			"namespace",
			"latestVersion",
		})
	}

	return
}

func runStats(config config.Config, info *prometheus.GaugeVec, timestamp *prometheus.GaugeVec) {
	if info != nil {
		info.Reset()
	}
	if timestamp != nil {
		timestamp.Reset()
	}

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
			updated := item.Info.LastDeployed.Unix() * 1000
			namespace := item.Namespace
			status := statusCodeMap[item.Info.Status.String()]
			latestVersion := ""

			if *fetchLatest {
				latestVersion = getLatestChartVersionFromHelm(item.Chart.Name(), config.HelmRegistries)
			}

			if info != nil {
				info.WithLabelValues(chart, releaseName, version, appVersion, strconv.FormatInt(updated, 10), namespace, latestVersion).Set(status)
			}
			if timestamp != nil {
				timestamp.WithLabelValues(chart, releaseName, version, appVersion, strconv.FormatInt(updated, 10), namespace, latestVersion).Set(float64(updated))
			}
		}
	}
}

func getLatestChartVersionFromHelm(name string, helmRegistries registries.HelmRegistries) (version string) {
	version = helmRegistries.GetLatestVersionFromHelm(name)
	log.WithField("chart", name).Debugf("last chart repo version is  %v", version)
	return
}

func runStatsPeriodically(interval time.Duration, config config.Config) {
	for {
		info, timestamp := configureMetrics()
		runStats(config, info, timestamp)
		registerMetrics(prometheus.DefaultRegisterer, info, timestamp)
		time.Sleep(interval)
	}
}

func registerMetrics(register prometheus.Registerer, info, timestamp *prometheus.GaugeVec) {
	mutex.Lock()
	defer mutex.Unlock()

	if statsInfo != nil {
		register.Unregister(statsInfo)
	}
	register.MustRegister(info)
	statsInfo = info

	if statsTimestamp != nil {
		register.Unregister(statsTimestamp)
	}
	register.MustRegister(timestamp)
	statsTimestamp = timestamp
}

func newHelmStatsHandler(config config.Config, synchrone bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if synchrone {
			runStats(config, statsInfo, statsTimestamp)
		} else {
			mutex.RLock()
			defer mutex.RUnlock()
		}

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
	cliFlags := initFlags()
	config := config.LoadConfiguration(cliFlags.ConfigFile)

	runIntervalDuration, err := time.ParseDuration(*intervalDuration)
	if err != nil {
		log.Fatalf("invalid duration `%s`: %s", *intervalDuration, err)
	}

	if namespaces == nil || *namespaces == "" {
		go informer()
	} else {
		for _, namespace := range strings.Split(*namespaces, ",") {
			connect(namespace)
		}
	}

	if runIntervalDuration != 0 {
		go runStatsPeriodically(runIntervalDuration, config)
	} else {
		info, timestamp := configureMetrics()
		registerMetrics(prometheus.DefaultRegisterer, info, timestamp)
	}

	http.HandleFunc("/metrics", newHelmStatsHandler(config, runIntervalDuration == 0))
	http.HandleFunc("/healthz", healthz)
	log.Fatal(http.ListenAndServe(":9571", nil))
}
