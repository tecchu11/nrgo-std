package nrslog

import "log/slog"

// HandlerOption is option for handler.
type HandlerOption func(*handler)

// WithHandler configures specific handler.
func WithHandler(parent slog.Handler) HandlerOption {
	return func(h *handler) {
		h.hn = parent
	}
}

// OnlyForward configures handler only send logs to newrelic.
func OnlyForward() HandlerOption {
	return func(h *handler) {
		h.onlyForward = true
	}
}
