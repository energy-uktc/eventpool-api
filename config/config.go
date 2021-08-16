package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

//Config is the application configuration file
type applicationConfig struct {
	Database string `mapstructure:"database"`
	Hostname string `mapstructure:"hostname"`
}

var Properties applicationConfig

//ReadConfig read the applications config file
func init() {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	configFile := "configs/config.yaml"
	if os.Getenv("ENVIRONMENT") != "" {
		configFile = fmt.Sprintf("configs/config-%s.yaml", os.Getenv("ENVIRONMENT"))
	}

	err := config.LoadFiles(configFile)
	if err != nil {
		log.Fatal(err)
	}

	err = config.BindStruct("", &Properties)
	if err != nil {
		log.Fatal("Configuration could not be loaded")
	}

}

//Load part of the configuration in struct
func Load(key string, destination interface{}) {
	err := config.BindStruct(key, &destination)
	if err != nil {
		log.Fatalf("Configuration could not be loaded for key %q", key)
	}
}

//Load part of the configuration in struct
func GetPath(key string) interface{} {
	return config.Get(key, true)
}
