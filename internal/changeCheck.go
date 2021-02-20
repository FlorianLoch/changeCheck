package internal

import (
	"io"
	"net/http"

	query "github.com/antchfx/htmlquery"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/net/html"
)

func checkPage(url string, xpath string, lastResult []*html.Node) (bool, []*html.Node, error) {
	body, err := fetchHTML(url)
	if err != nil {
		return false, nil, err
	}
	defer body.Close()

	nodes, err := parseAndFind(body, xpath)
	if err != nil {
		return false, nil, err
	}

	changeDetected := !compareNodes(lastResult, nodes)

	return changeDetected, nodes, nil
}

func fetchHTML(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

func parseAndFind(html io.Reader, xpath string) ([]*html.Node, error) {
	doc, err := query.Parse(html)
	if err != nil {
		return nil, err
	}

	nodes, err := query.QueryAll(doc, xpath)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func compareNodes(a, b []*html.Node) bool {
	return cmp.Equal(a, b)
}
