package nrhttp

import (
	"net/http"

	"github.com/newrelic/go-agent/v3/newrelic"
)

const missingPattern = "ErrorHandler"

func Middleware(app *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			p := r.Pattern
			if p == "" {
				p = missingPattern
			}
			txn := app.StartTransaction(p)
			defer txn.End()
			txn.SetWebRequestHTTP(r)
			w = txn.SetWebResponse(w)
			ctx := newrelic.NewContext(r.Context(), txn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
