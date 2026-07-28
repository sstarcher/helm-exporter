package registries

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/registry"

	"github.com/sstarcher/helm-exporter/versioning"
)

// ociClient is a single, shared OCI registry client. It is configured the
// same way the Helm CLI configures itself (see registry.NewClient in
// cmd/helm/root.go upstream): it reads whatever credentials are already on
// disk for `helm registry login`, honouring HELM_REGISTRY_CONFIG. No new
// authentication mechanism is introduced - private OCI registries work
// exactly like they do for `helm pull oci://...`.
var ociClient *registry.Client

func init() {
	client, err := registry.NewClient(
		registry.ClientOptCredentialsFile(settings.RegistryConfig),
	)
	if err != nil {
		log.Warning(err)
		return
	}
	ociClient = client
}

// ociRef builds the OCI reference (host/path/chart, no scheme) that the
// registry client expects, from a `oci://host/path` registry URL and a
// chart name. Split out as its own function so it can be unit tested
// without a network round trip.
func ociRef(registryURL, chart string) string {
	trimmed := strings.TrimSuffix(registryURL, "/")
	trimmed = strings.TrimPrefix(trimmed, "oci://")
	return trimmed + "/" + chart
}

// getChartVersionsOCI fetches the latest version of a chart published to
// an OCI registry. OCI registries have no index.yaml; instead every chart
// version is pushed as its own tag, so the available versions are the
// tags on the repository. That mirrors how `helm pull oci://... --version
// <tag>` resolves versions, and it lets us reuse the existing
// FindHighestVersionInList helper (the same one the rest of the exporter
// uses) instead of adding a second way to pick "the latest" version.
func (r HelmOverrideRegistry) getChartVersionsOCI(chart string) string {
	if ociClient == nil {
		log.Warning("OCI registry client is not initialized")
		return versioning.Failure
	}

	ref := ociRef(r.HelmRegistry.URL, chart)

	tags, err := ociClient.Tags(ref)
	if err != nil {
		log.Warning(err)
		return versioning.Failure
	}

	return versioning.FindHighestVersionInList(tags, r.AllowAllReleases)
}
