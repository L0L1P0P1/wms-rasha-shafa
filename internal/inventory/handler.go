package inventory

import (
	"net/http"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/ctxlog"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
)

type InventoryHandler struct {
	queries *db.Queries
}

func NewInventoryHandler(queries *db.Queries) *InventoryHandler {
	return &InventoryHandler{queries: queries}
}

func (h *InventoryHandler) something(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

}
