package internal

import (
	"sync"
	"time"
)

func Merge(chans ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup

	outChan := make(chan interface{})

	wg.Add(len(chans))

	for _, c := range chans {
		go func(c <-chan interface{}) {
			for v := range c {
				outChan <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		defer close(outChan)

		wg.Wait()
	}()

	return outChan
}

func Tick(interval time.Duration) <-chan interface{} {
	outChan := make(chan interface{})

	go func() {
		defer close(outChan)

		for range time.Tick(interval) {
			outChan <- nil
		}
	}()

	return outChan
}

func Shotgun() <-chan interface{} {
	c := make(chan interface{}, 1)
	c <- nil

	return c
}
