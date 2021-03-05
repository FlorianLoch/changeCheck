package persistence

import (
	"path"

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
			log.Debug().Msg("Change to config file detected.")

			outChan <- v
		}
	}()

	return outChan, nil
}

func monitorFile(file string) (<-chan interface{}, error) {
	changeChan := make(chan interface{})

	dirName := path.Dir(file)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = watcher.Add(dirName)
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

				if event.Name != file {
					break
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					changeChan <- nil
				}

				if event.Op&fsnotify.Create == fsnotify.Create {
					changeChan <- nil
				}
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
