package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/afero"

	"github.com/florianloch/change-check/internal"
	"github.com/florianloch/change-check/internal/notification"
	"github.com/florianloch/change-check/internal/persistence"
)

const (
	envAppBaseURL = "APP_BASE_URL"
	envInterface  = "INTERFACE"
	envPort       = "PORT"
)

var (
	// Build flags set by Makefile
	gitVersion    string
	gitAuthorDate string
	buildDate     string
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Info().Str("gitCommit", gitVersion).Str("gitDate", gitAuthorDate).Str("builtAt", buildDate).Msg("")

	config, err := persistence.LoadConfig()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed loading config. Make sure '%s' exists.", persistence.ConfigFileName)
	}

	log.Info().Msgf("Going to monitor %d page(s).", len(config.Pages))

	p := persistence.NewFSPersistor(afero.NewOsFs())

	appBaseURL := os.Getenv(envAppBaseURL)
	port := os.Getenv(envPort)
	interfaze := os.Getenv(envInterface)
	if interfaze == "" {
		interfaze = "0.0.0.0"
	}

	var d notification.Debouncer
	if appBaseURL != "" && port != "" {
		webDebouncer, err := notification.NewWebDebouncer(appBaseURL)
		if err != nil {
			log.Fatal().Err(err).Str("appBaseURL", appBaseURL).Msg("Could not initiliaze WebDebouncer.")
		}

		webDebouncer.StartHTTPServer(interfaze + ":" + port)

		log.Info().Msgf("Using WebDebouncer. Reachable at '%s'.", appBaseURL)

		d = webDebouncer
	} else {
		d = &notification.DummyDebouncer{}

		log.Info().Msgf("Using DummyDebouncer. Set '%s' and '%s' in order to use the WebDebouncer.", envPort, envAppBaseURL)
	}

	n, err := notification.NewTelegramNotifier(config.TelegramBotToken, config.TelegramChatID, d)
	if err != nil {
		log.Fatal().Err(err).Str("botToken", config.TelegramBotToken).Msg("Could not initialize Telegram client.")
	}

	monitor(config, p, n)
}

func monitor(config *persistence.Config, p persistence.Persistor, n notification.Nofifier) {
	for true {
		for _, page := range config.Pages {
			changed, err := checkPage(page, p)

			if err != nil {
				log.Error().Err(err).Stringer("url", page.URL).Str("xpath", page.XPath).Msg("Failed to check page.")
				continue
			}

			if changed {
				log.Info().Str("xpath", page.XPath).Msgf("'%s' changed.\n", page.URL)

				err = n.Notify(page.URL, page.Debounce)
				if err != nil {
					log.Error().Err(err).Stringer("url", page.URL).Str("xpath", page.XPath).Msg("Failed to notify.")
				}
			} else {
				log.Info().Str("xpath", page.XPath).Msgf("'%s' NOT changed.\n", page.URL)
			}
		}

		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func checkPage(page *persistence.PageEntry, p persistence.Persistor) (bool, error) {
	oldNodes, err := p.Load(page.URL, page.XPath)
	if err != nil {
		return false, err
	}

	changed, newNodes, err := internal.CheckPage(page.URL, page.XPath, oldNodes)
	if err != nil {
		return false, err
	}

	if len(newNodes) == 0 {
		log.Warn().Stringer("url", page.URL).Str("xpath", page.XPath).Msg("No matching nodes found!")
	}

	if changed {
		err = p.Store(page.URL, page.XPath, newNodes)
	}

	// oldNodes == nil when this target has not been checked yet
	return changed && oldNodes != nil, err
}
