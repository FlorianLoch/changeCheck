package persistence

import (
	u "net/url"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFsPersistor(t *testing.T) {
	assert := assert.New(t)

	memFS := afero.NewMemMapFs()

	node1 := "<html><body><h1>Hello</h1></body></html>"
	node2 := "<html><body><h1>Hello</h1></body></html>"

	dummyNodes := []*string{&node1, &node2}

	p := NewFSPersistor(memFS)

	dummyURL, _ := u.Parse("example.com")

	// First call should return nil, indicating this target has not been checked before
	nodes, err := p.Load(dummyURL, "//h1")
	assert.NoError(err)
	assert.Nil(nodes)

	err = p.Store(dummyURL, "//h1", dummyNodes)
	assert.NoError(err)

	nodes, err = p.Load(dummyURL, "//h1")
	assert.NoError(err)
	assert.Equal(dummyNodes, nodes)

	// Shoudl overwrite the existing file
	dummyNodes = []*string{&node1}
	err = p.Store(dummyURL, "//h1", dummyNodes)
	assert.NoError(err)

	nodes, err = p.Load(dummyURL, "//h1")
	assert.NoError(err)
	assert.Equal(dummyNodes, nodes)
}
