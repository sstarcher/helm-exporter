package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/sstarcher/helm-exporter/registries"
)

// Config of the lcm application, normally loaded from the config file
type Config struct {
	HelmRegistries registries.HelmRegistries `koanf:"helmRegistries"`
}

// New returns a config object from the given file
func New(fileName string) Config {
	log.WithField("configFile", fileName).Debug("Loading config file")

	var lcmConfig Config
	k := koanf.New(".")

	// load defaults
	if len(fileName) > 0 {
		if err := k.Load(file.Provider(fileName), yaml.Parser()); err != nil {
			log.WithError(err).Fatal("Error loading config")
		}
		if err := k.Unmarshal("", &lcmConfig); err != nil {
			log.WithError(err).Fatal("Error unmarshaling config")
		}
	}

	return lcmConfig
}
