package api

import (
	"io/fs"
	"net/http"

	"digital-eval-system/services/go-node/internal/chain"
	"digital-eval-system/services/go-node/internal/core"
	"digital-eval-system/services/go-node/internal/pybridge"
	"digital-eval-system/services/go-node/internal/storage"
)

type Handler struct {
	store      storage.Storage
	chain      *chain.Chain
	signerID   string
	registry   *core.ServiceRegistry
	embeddedUI fs.FS // embedded frontend static files (may be nil)
}

func NewHandlerWithRegistry(store storage.Storage, signerID string, registry *core.ServiceRegistry, embeddedUI fs.FS) *Handler {
	return &Handler{
		store:      store,
		chain:      chain.NewChain(store),
		signerID:   signerID,
		registry:   registry,
		embeddedUI: embeddedUI,
	}
}

func (h *Handler) pyClient() *pybridge.Client {
	return h.registry.MustGet("pybridge").(*pybridge.Client)
}

func (h *Handler) WithRouter() http.Handler {
	return NewRouter(h)
}
