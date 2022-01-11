package config

import (
	"fmt"
	"hash/fnv"
	"math/big"
	"path/filepath"
	"regexp"
	"strings"
	"github.com/aws/amazon-genomics-cli/internal/pkg/osutils"
)

const (
	configDirName  = ".agc"
	configFileName = "config.yaml"
)

var defaultConfig = Config{Format: Format{defaultFormat}}

type Client struct {
	configFilePath string
}

func NewConfigClient() (*Client, error) {
	homeDir, err := osutils.DetermineHomeDir()
	if err != nil {
		return nil, err
	}

	configDirPath := filepath.Join(homeDir, configDirName)

	if err := osutils.EnsureDirExistence(configDirPath); err != nil {
		return nil, err
	}

	configFilePath := filepath.Join(configDirPath, configFileName)

	return &Client{configFilePath: configFilePath}, nil
}

func hash(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("Cannot write a hash for string %q: %v", s, err))
	}

	hashValue := h.Sum32()
	hash62 := big.NewInt(int64(hashValue)).Text(62)

	return hash62
}

func sanitizeUserName(userName string) string {
	const maxUserNameLength = 10

	reg := regexp.MustCompile("[^A-Za-z0-9]+")
	sanitizedUserName := reg.ReplaceAllString(userName, "")
	if len(sanitizedUserName) > maxUserNameLength {
		sanitizedUserName = sanitizedUserName[:maxUserNameLength]
	}

	return sanitizedUserName
}

func userIdFromEmailAddress(emailAddress string) string {
	emailAddress = strings.ToLower(emailAddress)
	userName := emailAddress[:strings.IndexByte(emailAddress, '@')]
	sanitizedUserName := sanitizeUserName(userName)
	return sanitizedUserName + hash(emailAddress)
}
func (c Client) Read() (Config, error) {
	return c.loadFromFile()
}

func (c Client) loadFromFile() (Config, error) {
	configData, err := configFromYaml(c.configFilePath, defaultConfig)
	if err != nil {
		return defaultConfig, err
	}
	configData.User.Id = userIdFromEmailAddress(configData.User.Email)
	return configData, nil
}

func (c Client) storeToFile(config Config) error {
	return configToYaml(c.configFilePath, config)
}

func (c Client) GetUserEmailAddress() (string, error) {
	configData, err := c.loadFromFile()

	if err != nil {
		return "", fmt.Errorf("can not read email address from config file. Please initialize your AGC setup by running `agc configure email`")
	}
	return configData.User.Email, nil
}

func (c Client) SetUserEmailAddress(userEmailAddress string) error {
	configData, _ := c.loadFromFile()
	configData.User.Email = userEmailAddress
	return c.storeToFile(configData)
}

func (c Client) GetUserId() (string, error) {
	userEmailAddress, err := c.GetUserEmailAddress()
	if err != nil {
		return "", err
	}
	userId := userIdFromEmailAddress(userEmailAddress)
	return userId, nil
}

func (c Client) SetFormat(format string) error {
	configData, err := c.loadFromFile()
	if err != nil {
		return err
	}
	configData.Format.Name = format
	return c.storeToFile(configData)
}

func (c Client) GetFormat() (string, error) {
	configData, err := c.loadFromFile()
	if err != nil {
		return "", err
	}
	return configData.Format.Name, nil
}
