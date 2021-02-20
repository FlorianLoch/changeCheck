package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dummyHTML1  = "<html><body><h1>Hello</h1></body></html>"
	dummyHTML2  = "<html><body><h1>World</h1></body></html>"
	sampleXPath = "//h1"
)

func TestCheckingPageForChanges(t *testing.T) {
	assert := assert.New(t)

	nodesA, err := parseAndFind(strings.NewReader(dummyHTML1), sampleXPath)
	assert.NoError(err)

	nodesB, err := parseAndFind(strings.NewReader(dummyHTML2), sampleXPath)
	assert.NoError(err)

	changeDetected := !compareNodes(nodesA, nodesB)
	assert.True(changeDetected)

	nodesC, err := parseAndFind(strings.NewReader(dummyHTML1), sampleXPath)
	assert.NoError(err)

	changeDetected = !compareNodes(nodesA, nodesC)
	assert.False(changeDetected)
}
