package internal

import (
	"testing"
	"time"
)

func TestMerge(t *testing.T) {
	chan1 := make(chan interface{})
	chan2 := make(chan interface{})
	merged := Merge(chan1, chan2)

	go func() {
		chan1 <- nil
		chan2 <- nil
		chan2 <- nil
		chan1 <- nil
	}()

	// Blocks if we cannot get all entries
	for range [4]int{} {
		<-merged
	}
}

func TestTick(t *testing.T) {
	tickChan := Tick(time.Duration(1))

	<-tickChan
	<-tickChan
}
