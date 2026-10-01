package skus

import (
	"net/http"
	"strconv"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/ctxlog"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
)

// helper function to extract sku id from url param and luhn check
func GetSkuIDParam(r *http.Request, paramName string) (int64, error) {
	val, err := strconv.ParseInt(r.PathValue(paramName), 10, 64)
	if err != nil {
		return -1, err
	}

	if err := shared.VerifyLuhn(val/10, val%10); err != nil {
		return val, err
	}

	return val, nil
}

type SkuHandler struct {
	queries *db.Queries
}

func NewSkuHandler(queries *db.Queries) *SkuHandler {
	return &SkuHandler{queries: queries}
}

// Creates SKU, returns the SKU inserted into the database
// with a luhn checked id
func (h *SkuHandler) CreateSKU(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateSKUParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create sku", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	sku, err := h.queries.CreateSKU(r.Context(), params)
	if err != nil {
		logger.Error("failed to create sku in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// add luhn checksum to the displayed id
	sku.ID = sku.ID*10 + shared.CalculateLuhn(sku.ID)

	logger.Info("sku created successfully", "sku_id", sku.ID)
	shared.WriteJSON(w, http.StatusCreated, sku)
}

func (h *SkuHandler) GetSKU(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	id, err := GetSkuIDParam(r, "id")

	if err != nil {
		logger.Warn("invalid sku id", "id", id, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	// query on the id without the checksum (aka dbid)
	sku, err := h.queries.GetSKUByID(r.Context(), id/10)

	if err != nil {
		logger.Warn("sku not found", "dbid", id/10, "error", err)
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, sku)
}

func (h *SkuHandler) ListSKUs(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	q := r.URL.Query()
	limitStr := q.Get("limit")
	offsetStr := q.Get("offset")

	limit := int32(100)
	offset := int32(0)

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = int32(l)
		} else {
			logger.Warn("invalid limit query parameter, using default", "raw_limit", limitStr)
		}
	}
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = int32(o)
		} else {
			logger.Warn("invalid offset query parameter, using default", "raw_offset", offsetStr)
		}
	}

	skus, err := h.queries.ListSKUs(r.Context(), db.ListSKUsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to list skus", "limit", limit, "offset", offset, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for _, sku := range skus {
		sku.ID = sku.ID*10 + shared.CalculateLuhn(sku.ID)
	}

	shared.WriteJSON(w, http.StatusOK, skus)
}

func (h *SkuHandler) UpdateSKU(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	id, err := GetSkuIDParam(r, "id")
	if err != nil {
		logger.Warn("invalid sku id in path for update", "id", id, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	var params db.UpdateSKUParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for update sku", "sku_id", id, "error", err)
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id / 10

	sku, err := h.queries.UpdateSKU(r.Context(), params)
	if err != nil {
		logger.Error("failed to update sku in database", "dbid", id/10, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// add luhn checksum to the displayed id
	sku.ID = sku.ID*10 + shared.CalculateLuhn(sku.ID)

	logger.Info("sku updated successfully", "sku_id", sku.ID)
	shared.WriteJSON(w, http.StatusOK, sku)
}

func (h *SkuHandler) DeleteSKU(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	id, err := GetSkuIDParam(r, "id")

	if err != nil {
		logger.Warn("invalid sku id for delete", "id", id, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	_, err = h.queries.DeleteSKU(r.Context(), id/10)
	if err != nil {
		logger.Error("failed to delete sku from database", "dbid", id/10, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("sku deleted successfully", "dbid", id/10)
	w.WriteHeader(http.StatusNoContent)
}

// Packaging Units

func (h *SkuHandler) CreatePackagingUnit(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreatePackagingUnitParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for packaging unit", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	// verify luhn
	if err := shared.VerifyLuhn(params.SkuID/10, params.SkuID%10); err != nil {
		logger.Warn("Invalid sku id in params", "sku_id", params.SkuID, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	// remove luhn check
	params.SkuID /= 10

	unit, err := h.queries.CreatePackagingUnit(r.Context(), params)
	if err != nil {
		logger.Error("failed to create packaging unit in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// add luhn check
	unit.SkuID = 10*unit.SkuID + shared.CalculateLuhn(unit.SkuID)

	logger.Info("packaging unit created successfully", "packaging_unit_id", unit.ID)
	shared.WriteJSON(w, http.StatusCreated, unit)
}

func (h *SkuHandler) ListPackagingUnitsBySKU(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	skuId, err := GetSkuIDParam(r, "sku_id")

	if err != nil {
		logger.Warn("invalid sku_id in path for listing packaging units", "sku_id", skuId, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	units, err := h.queries.ListPackagingUnitsBySKU(r.Context(), skuId/10)
	if err != nil {
		logger.Error("failed to list packaging units for sku", "sku_dbid", skuId/10, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	for _, unit := range units {
		unit.SkuID = 10*unit.SkuID + shared.CalculateLuhn(unit.SkuID)
	}

	shared.WriteJSON(w, http.StatusOK, units)
}

// Create a Lot for an sku
func (h *SkuHandler) CreateLot(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateLotParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create lot", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	if err := shared.VerifyLuhn(params.SkuID/10, params.SkuID%10); err != nil {
		logger.Warn("Invalid sku_id for creating a lot", "sku_id", params.SkuID, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	// remove luhn check
	params.SkuID /= 10

	lot, err := h.queries.CreateLot(r.Context(), params)
	if err != nil {
		logger.Error("failed to create lot in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// add luhn checksum
	lot.SkuID = 10*lot.SkuID + shared.CalculateLuhn(lot.SkuID)

	logger.Info("lot created successfully", "lot_id", lot.ID)
	shared.WriteJSON(w, http.StatusCreated, lot)
}

// List all of the lots of a sku
func (h *SkuHandler) ListLotsBySKU(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	skuId, err := GetSkuIDParam(r, "sku_id")
	if err != nil {
		logger.Warn("invalid sku_id in path for listing lots", "sku_id", skuId, "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	lots, err := h.queries.ListLotsBySKU(r.Context(), skuId/10)

	if err != nil {
		logger.Error("failed to list lots for sku", "sku_dbid", skuId/10, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// add luhn checks back
	for _, lot := range lots {
		lot.SkuID = 10*lot.SkuID + shared.CalculateLuhn(lot.SkuID)
	}

	shared.WriteJSON(w, http.StatusOK, lots)
}
