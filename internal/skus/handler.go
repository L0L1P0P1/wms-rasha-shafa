package skus

import (
	"encoding/json"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"net/http"
)

type SkuHandler struct{}

func (*SkuHandler) GetSKUByID(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SKU_id string `json:"sku_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
}
