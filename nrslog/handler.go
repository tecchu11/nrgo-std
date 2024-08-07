package nrslog

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// Attributer can convert itself into a newrelic log attributes.
type Attributer interface {
	// NRAttribute return newrelic log event attribute(key-value pair)
	NRAttribute() map[string]string
}

type handler struct {
	slog.Handler
	app         *newrelic.Application
	onlyForward bool
}

// NewHandler creates [slog.Handler] can send log event to newrelic.
// Default output format is json.
func NewHandler(app *newrelic.Application, opt ...HandlerOption) slog.Handler {
	h := handler{
		Handler:     slog.NewJSONHandler(os.Stdout, nil),
		app:         app,
		onlyForward: false,
	}
	return &h
}

// Handle send log event to newrelic and write record with parent handler.
// Only send log event to newrelic if onlyForward is true.
func (h *handler) Handle(ctx context.Context, record slog.Record) error {
	txn := newrelic.FromContext(ctx)
	attrs := make(map[string]any)
	record.Attrs(func(a slog.Attr) bool {
		if nra, ok := a.Value.Any().(Attributer); ok {
			for k, v := range nra.NRAttribute() {
				attrs[strings.Join([]string{a.Key, k}, ".")] = v
			}
		} else {
			attrs[a.Key] = a.Value.Any()
		}
		return true
	})
	log := newrelic.LogData{
		Timestamp:  record.Time.UnixMilli(),
		Severity:   record.Level.String(),
		Message:    record.Message,
		Attributes: attrs,
	}
	if txn != nil {
		txn.RecordLog(log)
	} else {
		h.app.RecordLog(log)
	}
	if h.onlyForward {
		return nil
	}
	return h.Handler.Handle(ctx, record)
}
