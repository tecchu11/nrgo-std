package nrslog

import (
	"log/slog"

	"github.com/newrelic/go-agent/v3/newrelic"
)

type Handler = handler

func (h *Handler) GetHandler() slog.Handler {
	return h.hn
}

func (h *Handler) GetAPP() *newrelic.Application {
	return h.app
}

func (h *Handler) GetOnlyForward() bool {
	return h.onlyForward
}
