package persistence

import (
	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

func MonitorConfigFile() (<-chan interface{}, error) {
	configFilePath, err := configFilePath()
	if err != nil {
		return nil, err
	}

	inChan, err := monitorFile(configFilePath)
	if err != nil {
		return nil, err
	}

	outChan := make(chan interface{})

	go func() {
		defer close(outChan)

		for v := range inChan {
			log.Debug().Msg("Change to config file noticed.")

			outChan <- v
		}
	}()

	return outChan, nil
}

func monitorFile(file string) (<-chan interface{}, error) {
	changeChan := make(chan interface{})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(file)
	if err != nil {
		watcher.Close()
		return nil, err
	}

	go func() {
		defer func() {
			log.Debug().Msg("FileWatcher stopped.")
			close(changeChan)
			watcher.Close()
		}()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					changeChan <- nil
				}

				// TODO: Vim unlinks and recreates the file... this is nasty and does not work yet.
				// A work around might be watching the directory and then filter by the file's name.
			case err, ok := <-watcher.Errors:
				if !ok {
					log.Error().Err(err).Msg("FileWatcher stopped due to an error.")
					return
				}
			}
		}
	}()

	return changeChan, nil
}
