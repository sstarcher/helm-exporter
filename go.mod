module github.com/sstarcher/helm-exporter

go 1.13

// Pulled from https://github.com/helm/helm/blob/master/go.mod
// To ensure correct dependency resolution
replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

require (
	github.com/Azure/go-autorest v12.2.0+incompatible // indirect
	github.com/arminc/k8s-platform-lcm v0.0.0-20200523195516-127c372270b4 // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/flagenv v0.0.0-20160425205200-fcd59fca7456
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/knadh/koanf v0.10.0
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/onsi/ginkgo v1.12.0 // indirect
	github.com/onsi/gomega v1.9.0 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6
	github.com/prometheus/client_golang v1.2.1
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/helm/v3 v3.2.1
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3

)
