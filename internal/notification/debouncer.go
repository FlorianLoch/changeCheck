package notification

import (
	"encoding/base64"
	"net/http"
	u "net/url"
	"strings"
	"sync"

	"github.com/rs/zerolog/log"
)

type Debouncer interface {
	ShallNotify(url *u.URL) (bool, string)
}

type DummyDebouncer struct{}

func (d *DummyDebouncer) ShallNotify(url *u.URL) (bool, string) {
	return true, ""
}

type WebDebouncer struct {
	entries    map[string]struct{}
	lock       sync.Mutex
	appBaseURL *u.URL
}

func NewWebDebouncer(rawAppBaseURL string) (*WebDebouncer, error) {
	appURL, err := u.Parse(rawAppBaseURL)
	if err != nil {
		return nil, err
	}

	return &WebDebouncer{
		entries:    make(map[string]struct{}),
		appBaseURL: appURL,
	}, nil
}

func (d *WebDebouncer) StartHTTPServer(interfacePort string) {
	go func() {
		err := http.ListenAndServe(interfacePort, http.HandlerFunc(d.handlerFunc))
		if err != nil {
			log.Fatal().Err(err).Msg("Webserver of WebDebouncer shut down.")
		}
	}()
}

func (d *WebDebouncer) handlerFunc(w http.ResponseWriter, r *http.Request) {
	encoded := strings.Trim(r.URL.Path, "/")
	decoded, err := decode(encoded)
	if err != nil {
		http.Error(w, "Value given cannot be decoded as base64.", http.StatusBadRequest)
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	delete(d.entries, decoded)

	http.Redirect(w, r, decoded, http.StatusTemporaryRedirect)
}

func (d *WebDebouncer) ShallNotify(url *u.URL) (bool, string) {
	stringifiedURL := url.String()

	d.lock.Lock()
	defer d.lock.Unlock()

	_, ok := d.entries[stringifiedURL]
	if ok {
		return false, ""
	} else {
		d.entries[stringifiedURL] = struct{}{}

		tmp := *d.appBaseURL
		tmp.Path = encode(stringifiedURL)

		return true, tmp.String()
	}
}

func encode(url string) string {
	return base64.StdEncoding.EncodeToString([]byte(url))
}

func decode(encodedURL string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedURL)
	return string(decodedBytes), err
}
