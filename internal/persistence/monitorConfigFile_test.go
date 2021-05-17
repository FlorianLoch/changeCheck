package persistence

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonitorFile(t *testing.T) {
	assert := assert.New(t)

	cwd, err := os.Getwd()
	if err != nil {
		assert.FailNow("Could not get current working directory", err)
	}

	file, err := ioutil.TempFile(cwd, "")
	if err != nil {
		assert.FailNow("Could not create temp file", err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	changeChan, err := monitorFile(file.Name())
	if err != nil {
		assert.FailNow("Could not start monitoring file", err)
	}

	// Peek the channel
	select {
	case <-changeChan:
		assert.FailNow("changeChan should not contain an element")
	default:
		// noop
	}

	randomByte := make([]byte, 1)
	_, err = rand.Read(randomByte)
	file.Write(randomByte)

	// Blocks until the change is detected
	<-changeChan
}
