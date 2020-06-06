package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/sstarcher/helm-exporter/registries"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Config of the lcm application, normally loaded from the config file
type Config struct {
	HelmRegistries registries.HelmRegistries `koanf:"helmRegistries"`
}

// AppConfig is the config for the app which can be set trough cli and config
type AppConfig struct {
	ConfigFile string
}

// LoadConfiguration loads the configuration from file
func LoadConfiguration(configFile string) Config {
	log.WithField("configFile", configFile).Debug("Loading config file")

	var lcmConfig Config
	k := koanf.New(".")

	// load defaults
	if len(configFile) > 0 {
		if err := k.Load(file.Provider(configFile), yaml.Parser()); err != nil {
			log.WithError(err).Fatal("Error loading config")
		}
		if err := k.Unmarshal("", &lcmConfig); err != nil {
			log.WithError(err).Fatal("Error unmarshaling config")
		}
	}

	return lcmConfig
}
