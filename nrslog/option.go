package nrslog

import "log/slog"

// HandlerOption is option for handler.
type HandlerOption func(*handler)

// WithHandler configures specific handler.
func WithHandler(h slog.Handler) HandlerOption {
	return func(h *handler) {
		h.Handler = h
	}
}

// OnlyForward configures handler only send logs to newrelic.
func OnlyForward() HandlerOption {
	return func(h *handler) {
		h.onlyForward = true
	}
}
