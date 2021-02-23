package persistence

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const (
	dummyConfigStr = `
interval: 1000
telegram_bot_token: some secret token
telegram_chat_id: 12345
pages:
  - url: example.com/index.html
    xpath: //h1
`
)

var (
	dummyConfig = &Config{
		Interval:         1000,
		TelegramBotToken: "some secret token",
		TelegramChatID:   12345,
		Pages: []*PageEntry{{
			URL:   "example.com/index.html",
			XPath: "//h1",
		}},
	}
)

func TestLoadConfig(t *testing.T) {
	assert := assert.New(t)

	config, err := parseConfig([]byte(dummyConfigStr))
	assert.NoError(err)
	assert.Equal(dummyConfig, config)
}

func TestReadConfigFile(t *testing.T) {
	assert := assert.New(t)

	memFS := afero.NewMemMapFs()

	cwd, err := os.Getwd()
	assert.NoError(err)

	configFile := path.Join(cwd, configFileName)
	memFS.MkdirAll(cwd, 0755)
	err = afero.WriteFile(memFS, configFile, []byte(dummyConfigStr), 0644)
	assert.NoError(err)

	configBytes, err := readConfigFile(memFS)
	assert.NoError(err)
	assert.Equal(dummyConfigStr, string(configBytes))
}
