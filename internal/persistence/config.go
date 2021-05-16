package persistence

import (
	"errors"
	u "net/url"
	"os"
	"path"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFileName = "change-check.config.yaml"
)

var (
	ErrInvalidURL = errors.New("URL invalid. Make sure scheme and host are set.")
)

type Config struct {
	Interval         int    // in seconds
	TelegramBotToken string `yaml:"telegram_bot_token"`
	TelegramChatID   int64  `yaml:"telegram_chat_id"`
	Pages            []*PageEntry
}

type PageEntry struct {
	Name     string
	RawURL   string `yaml:"url"`
	URL      *u.URL `yaml:"-"`
	XPath    string
	Debounce bool
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
	configFilePath, err := configFilePath()
	if err != nil {
		return nil, err
	}

	bytes, err := afero.ReadFile(fs, configFilePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func configFilePath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(cwd, ConfigFileName), nil
}

func parseConfig(data []byte) (*Config, error) {
	conf := &Config{}

	err := yaml.Unmarshal(data, conf)
	if err != nil {
		return nil, err
	}

	for _, pageEntry := range conf.Pages {
		url, err := u.Parse(pageEntry.RawURL)
		if err != nil {
			return nil, err
		}

		if url.Scheme == "" && url.Host == "" {
			return nil, ErrInvalidURL
		}

		if pageEntry.XPath == "" {
			pageEntry.XPath = "/"
		}

		pageEntry.URL = url

		if pageEntry.Name == "" {
			pageEntry.Name = url.Host + url.Path
		}
	}

	return conf, nil
}
