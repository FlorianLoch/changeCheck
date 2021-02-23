package persistence

import (
	"os"
	"path"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	configFileName = "changeCheck.config.yaml"
)

type Config struct {
	Interval         int    // in seconds
	TelegramBotToken string `yaml:"telegram_bot_token"` // if set to "ENV" the value will be fetched from environment
	TelegramChatID   int64  `yaml:"telegram_chat_id"`
	Pages            []*PageEntry
}

type PageEntry struct {
	URL   string
	XPath string
}

func LoadConfig() (*Config, error) {
	fs := afero.NewOsFs()
	configBytes, err := readConfigFile(fs)
	if err != nil {
		return nil, err
	}

	return parseConfig(configBytes)
}

func readConfigFile(fs afero.Fs) ([]byte, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	configFilePath := path.Join(cwd, configFileName)

	bytes, err := afero.ReadFile(fs, configFilePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func parseConfig(data []byte) (*Config, error) {
	conf := Config{}

	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
