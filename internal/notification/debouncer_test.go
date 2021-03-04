package notification

import (
	"net/http"
	u "net/url"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
)

const (
	appBaseURL = "http://change-check.fdlo.ch"
)

func TestDebouncing(t *testing.T) {
	assert := assert.New(t)
	debouncer, e := setup(t)

	dummyURL, _ := u.Parse("http://example.com/testPage")
	dummyURL2, _ := u.Parse("http://exmaple.com/someUnknownPage")

	shallNotify, relayURL := debouncer.ShallNotify(dummyURL)
	assert.True(shallNotify)
	assert.NotEmpty(relayURL)
	assert.Contains(relayURL, appBaseURL)

	shallNotify, relayURL2 := debouncer.ShallNotify(dummyURL)
	assert.False(shallNotify)
	assert.Empty(relayURL2)

	url, err := u.Parse(relayURL)
	assert.NoError(err)

	r := e.GET(url.Path).Expect()
	r.Status(http.StatusTemporaryRedirect)
	r.Header("Location").Equal(dummyURL.String())

	// Should forward to any page, also unknown ones. This is done in order to not invalidate
	// forwardings after being requested once resp. in order to not have to keep state
	r = e.GET("/" + encode(dummyURL2.String())).Expect()
	r.Status(http.StatusTemporaryRedirect)
	r.Header("Location").Equal(dummyURL2.String())

	assert.True(debouncer.ShallNotify(dummyURL))
}

func setup(t *testing.T) (*WebDebouncer, *httpexpect.Expect) {
	assert := assert.New(t)

	debouncer, err := NewWebDebouncer("/just-a-path-lacking-host-and-scheme")
	assert.Error(err, ErrInvalidURL)
	assert.Nil(debouncer)

	debouncer, err = NewWebDebouncer(appBaseURL)
	assert.NoError(err)

	handler := debouncer.handlerFunc

	e := httpexpect.WithConfig(httpexpect.Config{
		BaseURL: appBaseURL,
		Client: &http.Client{
			Transport: httpexpect.NewBinder(http.HandlerFunc(handler)),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})

	return debouncer, e
}
