package registries

import (
	"context"
	"github.com/sstarcher/helm-exporter/versioning"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
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

const indexYamlSuffix = "/index.yaml"

var clientSet kubernetes.Interface
var settings *cli.EnvSettings

func init() {
	settings = cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Warning(err)
	}

	var err error
	clientSet, err = actionConfig.KubernetesClientSet()
	if err != nil {
		log.Warning(err)
	}
}

func (r HelmOverrideRegistry) getChartVersions(chart string) string {

	// trim the index.yaml suffix from the chart url, just to avoid breaking changes.
	url := strings.TrimSuffix(r.HelmRegistry.URL, indexYamlSuffix)

	entry := &repo.Entry{
		Name: chart,
		URL:  url,
	}

	if clientSet == nil {
		log.Warning("kubernetes ClientSet is not initialized")
		return versioning.Failure
	}

	if r.HelmRegistry.SecretRef != nil {
		secrets, err := clientSet.CoreV1().Secrets(settings.Namespace()).Get(context.Background(), r.HelmRegistry.SecretRef.Name, v1.GetOptions{})
		if err != nil {
			log.Warning(err)
			return versioning.Failure
		}
		entry.Username = string(secrets.Data[r.HelmRegistry.SecretRef.UserKey])
		entry.Password = string(secrets.Data[r.HelmRegistry.SecretRef.PassKey])
	}

	provider := getter.All(settings)

	chartRepo, err := repo.NewChartRepository(entry, provider)
	if err != nil {
		log.Warning(err)
		return versioning.Failure
	}
	idx, err := chartRepo.DownloadIndexFile()
	if err != nil {
		log.Warning(err)
		return versioning.Failure
	}
	repoIndex, err := repo.LoadIndexFile(idx)
	if err != nil {
		log.Warning(err)
		return versioning.Failure
	}
	chartVersion, err := repoIndex.Get(chart, "")
	if err != nil {
		log.Warning(err)
		return versioning.Failure
	}
	return chartVersion.Version
}
