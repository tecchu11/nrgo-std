package nrhttp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/tecchu11/nrgo-std/nrhttp"
)

func TestMiddleware_ApplicationIsNotNil(t *testing.T) {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("test"),
		newrelic.ConfigLicense("0000000000000000000000000000000000000000"),
	)
	if err != nil {
		t.Fatal(err)
	}
	middlewared := nrhttp.Middleware(app)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(nil)
		if err != nil {
			t.Fatal(err)
		}
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest("", "/", nil)

	middlewared.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Errorf("response code must be 200 but got %d", w.Code)
	}
	if strings.TrimSpace(w.Body.String()) != "null" {
		t.Errorf("response body must be null but got %v", w.Body.String())
	}
}
