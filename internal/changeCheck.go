package internal

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	u "net/url"

	query "github.com/antchfx/htmlquery"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

var (
	ErrNoMatchingNodes = errors.New("no matching nodes found for xpath expression")
)

func CheckPage(url *u.URL, xpath string, lastResult []*string) (bool, []*string, error) {
	body, err := fetchHTML(url)
	if err != nil {
		return false, nil, err
	}
	defer body.Close()

	renderedNodes, err := parseAndFind(body, xpath)
	if err != nil {
		return false, nil, err
	}

	changeDetected := !compareNodes(lastResult, renderedNodes)

	return changeDetected, renderedNodes, nil
}

func fetchHTML(url *u.URL) (io.ReadCloser, error) {
	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func parseAndFind(html io.Reader, xpath string) ([]*string, error) {
	doc, err := query.Parse(html)
	if err != nil {
		return nil, err
	}

	nodes, err := query.QueryAll(doc, xpath)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 0 {
		return nil, ErrNoMatchingNodes
	}

	return renderNodes(nodes)
}

func renderNodes(nodes []*html.Node) ([]*string, error) {
	renderedNodes := make([]*string, 0, len(nodes))

	for _, node := range nodes {
		var buf bytes.Buffer
		err := html.Render(&buf, node)
		if err != nil {
			return nil, err
		}

		str := buf.String()
		renderedNodes = append(renderedNodes, &str)
	}

	return renderedNodes, nil
}

func compareNodes(a, b []*string) bool {
	return cmp.Equal(a, b)
}
