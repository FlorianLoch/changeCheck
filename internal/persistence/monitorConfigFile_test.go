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

	file, err := ioutil.TempFile("", "")
	assert.NoError(err)
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	changeChan, err := monitorFile(file.Name())
	assert.NoError(err)

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
