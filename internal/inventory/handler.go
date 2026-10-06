package inventory

import (
	"net/http"
	"strconv"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/ctxlog"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
	"github.com/jackc/pgx/v5/pgtype"
)

type InventoryHandler struct {
	queries *db.Queries
}

func NewInventoryHandler(queries *db.Queries) *InventoryHandler {
	return &InventoryHandler{queries: queries}
}

func (h *InventoryHandler) CreateAllocation(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateAllocationParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create allocation", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	allocation, err := h.queries.CreateAllocation(r.Context(), params)
	if err != nil {
		logger.Error("failed to create allocation in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("allocation created successfully", "allocation_id", allocation.ID)
	shared.WriteJSON(w, http.StatusCreated, allocation)
}

func (h *InventoryHandler) DeleteAllocation(w http.ResponseWriter, r *http.Request) {
	// logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	_, err = h.queries.DeleteAllocation(r.Context(), id)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InventoryHandler) FindReservableBalances(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	arg := db.FindReservableBalancesParams{}
	balances, err := h.queries.FindReservableBalances(r.Context(), arg)
	if err != nil {
		logger.Error("failed to find reservable balances", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, balances)
}

func (h *InventoryHandler) GetAllocationsByOrderLine(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("order_line_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	allocations, err := h.queries.GetAllocationsByOrderLine(r.Context(), id)
	if err != nil {
		logger.Error("failed to get allocations by order line", "order_line_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, allocations)
}

func (h *InventoryHandler) IncrementAllocationPicked(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.IncrementAllocationPickedParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for increment allocation picked", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	allocation, err := h.queries.IncrementAllocationPicked(r.Context(), params)
	if err != nil {
		logger.Error("failed to increment allocation picked", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("allocation picked incremented successfully", "allocation_id", allocation.ID)
	shared.WriteJSON(w, http.StatusOK, allocation)
}

func (h *InventoryHandler) InsertMovement(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.InsertMovementParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for insert movement", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	movement, err := h.queries.InsertMovement(r.Context(), params)
	if err != nil {
		logger.Error("failed to insert movement in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("movement inserted successfully", "movement_id", movement.ID)
	shared.WriteJSON(w, http.StatusCreated, movement)
}

func (h *InventoryHandler) ListBalancesByNode(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("node_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	balances, err := h.queries.ListBalancesByNode(r.Context(), id)
	if err != nil {
		logger.Error("failed to list balances by node", "node_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, balances)
}

func (h *InventoryHandler) ListMovementsByNode(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("node_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	movements, err := h.queries.ListMovementsByNode(r.Context(), db.ListMovementsByNodeParams{
		SourceNodeID: pgtype.Int8{Int64: id, Valid: true},
	})
	if err != nil {
		logger.Error("failed to list movements by node", "node_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, movements)
}

func (h *InventoryHandler) ListMovementsBySKU(w http.ResponseWriter, r *http.Request) {
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

	skus, err := h.queries.ListMovementsBySKU(r.Context(), db.ListMovementsBySKUParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to list movements by sku", "limit", limit, "offset", offset, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, skus)
}

func (h *InventoryHandler) LockBalanceForUpdate(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	balance, err := h.queries.LockBalanceForUpdate(r.Context(), db.LockBalanceForUpdateParams{
		NodeID: id,
	})
	if err != nil {
		logger.Error("failed to lock balance for update", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, balance)
}

func (h *InventoryHandler) ProjectDestinationAddition(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.ProjectDestinationAdditionParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for project destination addition", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	err := h.queries.ProjectDestinationAddition(r.Context(), params)
	if err != nil {
		logger.Error("failed to project destination addition", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *InventoryHandler) ProjectSourceDeduction(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.ProjectSourceDeductionParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for project source deduction", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	result, err := h.queries.ProjectSourceDeduction(r.Context(), params)
	if err != nil {
		logger.Error("failed to project source deduction", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, result)
}

func (h *InventoryHandler) UpdateBalanceAllocation(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.UpdateBalanceAllocationParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for update balance allocation", "error", err)
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	result, err := h.queries.UpdateBalanceAllocation(r.Context(), params)
	if err != nil {
		logger.Error("failed to update balance allocation", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, result)
}

