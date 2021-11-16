package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

var (
	writeFile = ioutil.WriteFile
	readFile  = ioutil.ReadFile
)

const (
	defaultFormat = "text"
)

func toYaml(filePath string, configData Config) error {
	bytes, err := yaml.Marshal(configData)
	if err != nil {
		return err
	}
	return writeFile(filePath, bytes, 0644)
}

func fromYaml(filePath string) (Config, error) {
	configData := Config{Format: Format{defaultFormat}}
	bytes, err := readFile(filePath)
	if err != nil {
		return configData, err
	}
	if err := yaml.Unmarshal(bytes, &configData); err != nil {
		return configData, err
	}
	return configData, nil
}
