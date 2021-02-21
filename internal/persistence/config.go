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

type config struct {
	Interval         int    // in milliseconds
	TelegramToken    string `yaml:"telegram_token"` // if set to "ENV" the value will be fetched from environment
	TelegramReceiver string `yaml:"telegram_receiver"`
	Pages            []*PageEntry
}

type PageEntry struct {
	URL   string
	XPath string
}

func LoadConfig(fs afero.Fs) ([]byte, error) {
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

func ParseConfig(data []byte) (*config, error) {
	conf := config{}

	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
