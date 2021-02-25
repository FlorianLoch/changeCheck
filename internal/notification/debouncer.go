package notification

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Debouncer interface {
	Debounce(url, xpath string) bool
}

type WebDebouncer struct {
	entries map[string]string
	lock    sync.Mutex
}

func NewWebDebouncer() *WebDebouncer {
	return &WebDebouncer{
		entries: make(map[string]string),
	}
}

func (d *WebDebouncer) StartHTTPServer(interfacePort string) {
	http.ListenAndServe(interfacePort, http.HandlerFunc(d.handlerFunc))
}

func (d *WebDebouncer) handlerFunc(w http.ResponseWriter, r *http.Request) {
	hashed := strings.Trim(r.URL.Path, "/")

	d.lock.Lock()
	defer d.lock.Unlock()
	url, ok := d.entries[hashed]

	if !ok {
		http.Error(w, fmt.Sprintf("'%s' is not known.", hashed), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (d *WebDebouncer) Debounce(url, xpath string) bool {

}

func hash(url, xpath string) string {
	hashed := sha256.Sum256([]byte(url + ":" + xpath))
	return fmt.Sprintf("%x", hash)
}
