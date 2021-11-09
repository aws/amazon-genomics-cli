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
	bytes, err := readFile(filePath)
	if err != nil {
		return Config{}, err
	}
	var configData Config
	setDefault(&configData)
	if err := yaml.Unmarshal(bytes, &configData); err != nil {
		return Config{}, err
	}
	return configData, nil
}
func setDefault(configData *Config) Config {
	format := configData.Format.Format
	if format == "" {
		configData.Format.Format = defaultFormat
	}
	return *configData
}
