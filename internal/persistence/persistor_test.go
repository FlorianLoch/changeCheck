package persistence

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestFsPersistor(t *testing.T) {
	assert := assert.New(t)

	memFS := afero.NewMemMapFs()

	node1, err := html.Parse(strings.NewReader("<html><body><h1>Hello</h1></body></html>"))
	assert.NoError(err)
	node2, err := html.Parse(strings.NewReader("<html><body><h1>Hello</h1></body></html>"))
	assert.NoError(err)

	dummyNodes := []*html.Node{node1, node2}

	p := NewFSPersistor(memFS)

	// First call should return empty node slice
	nodes, err := p.Load("example.com", "//h1")
	assert.NoError(err)
	assert.Empty(nodes)

	err = p.Store("example.com", "//h1", dummyNodes)
	assert.NoError(err)

	// First call should return empty node slice
	nodes, err = p.Load("example.com", "//h1")
	assert.NoError(err)
	assert.Equal(dummyNodes, nodes)
}
