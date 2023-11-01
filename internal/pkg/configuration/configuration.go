// TODO
package configuration

import (
	"encoding/json"
	"flag"
	"os"
	"sync"
)

type Configuration struct {
	Mode string `json:"mode" env:"REMEMBRANCE_MODE" cmd:"mode"`
}

const defaultConfigurationFilePath = "./config.json"

// TODO make a pointer and protect with a mutex
var (
	currentConfiguration *Configuration

	loaded bool = false

	guard sync.RWMutex
)

func GetConfiguration(forceReload bool, configurationFilePath string) (Configuration, error) {
	guard.RLock()

	if !loaded || forceReload {
		guard.RUnlock()
		guard.Lock()
		defer guard.Unlock()

		if !loaded || forceReload {
			temp := Configuration{}
			err := loadConfiguration(&temp, configurationFilePath)
			if err != nil {
				return *currentConfiguration, err
			}

			currentConfiguration = &temp
			loaded = true
		}
	}

	return *currentConfiguration, nil
}

func applyEnvironmentVariablesValues(configuration *Configuration) {
	mode, present := os.LookupEnv("REMEMBRANCE_MODE")
	if present {
		configuration.Mode = mode
	}
}

func applyConfigurationFileValues(configuration *Configuration, configurationFilePath string) error {
	var actualConfigurationFilePath string
	if configurationFilePath == "" {
		actualConfigurationFilePath = defaultConfigurationFilePath

	} else {
		actualConfigurationFilePath = configurationFilePath
	}

	file, err := os.ReadFile(actualConfigurationFilePath)
	if err != nil {
		return err
	}

	var rawConfiguration interface{}
	marshallError := json.Unmarshal(file, &rawConfiguration)
	if marshallError != nil {
		return marshallError
	}

	jsonConfiguration := rawConfiguration.(map[string]interface{})

	mode, found := jsonConfiguration["mode"]
	if found {
		switch mode := mode.(type) {
		case string:
			if mode != "" {
				configuration.Mode = mode
			}
		}
	}

	return nil
}

func applyDefaultValues(configuration *Configuration) {
	configuration.Mode = "default"
}

func applyCommandLineValues(configuration *Configuration) {
	const blankModeValue = "blank"

	mode := flag.String("mode", blankModeValue, "The mode this node has to be started")

	if !flag.Parsed() {
		flag.Parse()
	}

	if (*mode) != blankModeValue {
		configuration.Mode = *mode
	}
}

func loadConfiguration(configuration *Configuration, configurationFilePath string) error {
	applyDefaultValues(configuration)

	err := applyConfigurationFileValues(configuration, configurationFilePath)
	if err != nil {
		return err
	}

	applyEnvironmentVariablesValues(configuration)

	applyCommandLineValues(configuration)

	return nil
}
