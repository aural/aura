package aural

import (
	"log"

	prefer "github.com/LimpidTech/prefer.go"
)

type Configuration struct {
	Address string `yaml:"address";json:"address"`
}

func GetConfiguration() Configuration {
	configuration := Configuration{
		Address: "tcp://127.0.0.1:9090",
	}

	configurator, err := prefer.Load("aural", &configuration)
	if err != nil {
		return configuration
	}

	log.Println("Using configuration file:", configurator.Identifier)
	return configuration
}
