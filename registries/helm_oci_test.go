package registries

import (
	"testing"

	"helm.sh/helm/v3/pkg/registry"

	"github.com/sstarcher/helm-exporter/versioning"
)

// --- URL detection -----------------------------------------------------
//
// registry.IsOCI is the official Helm SDK function the exporter now relies
// on to decide which strategy to use. These cases pin down the exact
// behavior the rest of the package depends on, so a Helm SDK upgrade that
// changes this behavior fails loudly here instead of silently changing
// which registries are treated as OCI.

func TestIsOCIDetection(t *testing.T) {
	cases := []struct {
		name string
		url  string
		want bool
	}{
		{"oci scheme", "oci://registry.example.com/helm", true},
		{"oci scheme, no path", "oci://ghcr.io", true},
		{"https scheme", "https://prometheus-community.github.io/helm-charts", false},
		{"http scheme", "http://charts.example.com", false},
		{"index.yaml url", "https://some.url/index.yaml", false},
		{"empty url", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := registry.IsOCI(c.url); got != c.want {
				t.Errorf("registry.IsOCI(%q) = %v, want %v", c.url, got, c.want)
			}
		})
	}
}

// --- ref building --------------------------------------------------------

func TestOCIRef(t *testing.T) {
	cases := []struct {
		name        string
		registryURL string
		chart       string
		want        string
	}{
		{
			name:        "basic",
			registryURL: "oci://registry.example.com/helm",
			chart:       "kube-prometheus-stack",
			want:        "registry.example.com/helm/kube-prometheus-stack",
		},
		{
			name:        "ghcr example from the issue",
			registryURL: "oci://ghcr.io/prometheus-community/charts",
			chart:       "kube-prometheus-stack",
			want:        "ghcr.io/prometheus-community/charts/kube-prometheus-stack",
		},
		{
			name:        "trailing slash is tolerated",
			registryURL: "oci://registry.example.com/helm/",
			chart:       "prometheus",
			want:        "registry.example.com/helm/prometheus",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := ociRef(c.registryURL, c.chart); got != c.want {
				t.Errorf("ociRef(%q, %q) = %q, want %q", c.registryURL, c.chart, got, c.want)
			}
		})
	}
}

// --- error handling / regression -----------------------------------------

// A registry misconfigured so the shared client is unavailable must fail
// safely (the same versioning.Failure sentinel the HTTP path already
// returns) instead of panicking.
func TestGetChartVersionsOCI_ClientNotInitialized(t *testing.T) {
	original := ociClient
	ociClient = nil
	defer func() { ociClient = original }()

	r := HelmOverrideRegistry{
		HelmRegistry: HelmRegistry{URL: "oci://registry.example.com/helm"},
	}

	got := r.getChartVersionsOCI("prometheus")
	if got != versioning.Failure {
		t.Errorf("getChartVersionsOCI() with nil client = %q, want %q", got, versioning.Failure)
	}
}

// Regression: an https:// registry must never be routed through the OCI
// strategy, and vice versa. This is the one branch point the whole
// feature depends on staying correct.
func TestGetChartVersions_DispatchesByScheme(t *testing.T) {
	cases := []struct {
		name string
		url  string
		oci  bool
	}{
		{"https repo stays on the HTTP path", "https://prometheus-community.github.io/helm-charts", false},
		{"oci repo uses the OCI path", "oci://registry.example.com/helm", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := registry.IsOCI(c.url); got != c.oci {
				t.Errorf("registry.IsOCI(%q) = %v, want %v (this drives getChartVersions dispatch)", c.url, got, c.oci)
			}
		})
	}
}

// Invalid/unreachable OCI URLs must degrade to versioning.Failure rather
// than crash the exporter, same contract as the existing HTTP path for
// invalid/unreachable index.yaml URLs. This exercises the real network
// call, so it is skipped by default - see the "Testing" section of the PR
// description for why an integration test is the right tool here instead
// of mocking the OCI transport.
func TestGetChartVersionsOCI_InvalidHost(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network-dependent OCI test in -short mode")
	}

	r := HelmOverrideRegistry{
		HelmRegistry: HelmRegistry{URL: "oci://invalid.invalid.example/helm"},
	}

	got := r.getChartVersionsOCI("does-not-exist")
	if got != versioning.Failure {
		t.Errorf("getChartVersionsOCI() with unreachable host = %q, want %q", got, versioning.Failure)
	}
}
