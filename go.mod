module github.com/sstarcher/helm-exporter

go 1.17

// Pulled from https://github.com/helm/helm/blob/master/go.mod
// To ensure correct dependency resolution
// replace github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309

require (
	cloud.google.com/go v0.97.0 // indirect
	github.com/Azure/go-autorest/autorest v0.11.21 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.16 // indirect
	github.com/BurntSushi/toml v0.4.1 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/Microsoft/hcsshim v0.8.22 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/bugsnag/bugsnag-go v1.5.0 // indirect
	github.com/bugsnag/panicwrap v1.2.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/containerd/containerd v1.5.7 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/docker/cli v20.10.9+incompatible // indirect
	github.com/docker/docker v20.10.9+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20210407135951-1de76d718b3f // indirect
	github.com/facebookgo/ensure v0.0.0-20160127193407-b4ab57deab51 // indirect
	github.com/facebookgo/flagenv v0.0.0-20160425205200-fcd59fca7456
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/subset v0.0.0-20150612182917-8dac2c3c4870 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/garyburd/redigo v1.6.0 // indirect
	github.com/go-errors/errors v1.4.1 // indirect
	github.com/go-logr/logr v1.1.0 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gofrs/uuid v3.2.0+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/jmoiron/sqlx v1.3.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kardianos/osext v0.0.0-20170510131534-ae77be60afb1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/knadh/koanf v1.2.4
	github.com/lib/pq v1.10.3 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.4.2 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20210501183033-44dafcb38ecc
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.31.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cast v1.4.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/gorelic v0.0.6 // indirect
	go.starlark.net v0.0.0-20210901212718-87f333178d59 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211005001312-d4b1ae081e3b // indirect
	golang.org/x/sys v0.0.0-20211004093028-2c5d950f24ef // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20211001223012-bfb93cce50d9 // indirect
	google.golang.org/grpc v1.41.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.7.0
	k8s.io/apiextensions-apiserver v0.22.2 // indirect
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/klog/v2 v2.20.0 // indirect
	k8s.io/kube-openapi v0.0.0-20210929172449-94abcedd1aa4 // indirect
	k8s.io/kubectl v0.22.2 // indirect
	k8s.io/utils v0.0.0-20210930125809-cb0fa318a74b // indirect
	sigs.k8s.io/kustomize/api v0.10.0 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
