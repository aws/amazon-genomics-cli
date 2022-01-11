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

func configToYaml(filePath string, configData Config) error {
	bytes, err := yaml.Marshal(configData)
	if err != nil {
		return err
	}
	return writeFile(filePath, bytes, 0644)
}

func configFromYaml(filePath string, configData Config) (Config, error) {
	bytes, err := readFile(filePath)
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(bytes, &configData); err != nil {
		return Config{}, err
	}
	return configData, nil
}
