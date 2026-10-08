package inbound_orders

import (
	"net/http"
	"strconv"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/ctxlog"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
)

type InboundOrderHandler struct {
	queries *db.Queries
}

func NewInboundOrderHandler(queries *db.Queries) *InboundOrderHandler {
	return &InboundOrderHandler{queries: queries}
}

func (h *InboundOrderHandler) CreateInboundOrder(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateInboundOrderParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create inbound order", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	order, err := h.queries.CreateInboundOrder(r.Context(), params)
	if err != nil {
		logger.Error("failed to create inbound order in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("inbound order created successfully", "order_id", order.ID)
	shared.WriteJSON(w, http.StatusCreated, order)
}

func (h *InboundOrderHandler) GetInboundOrderByID(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	order, err := h.queries.GetInboundOrderByID(r.Context(), id)
	if err != nil {
		logger.Warn("inbound order not found", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, order)
}

func (h *InboundOrderHandler) GetInboundOrderByPONumber(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	poNumber := r.PathValue("po_number")
	order, err := h.queries.GetInboundOrderByPONumber(r.Context(), poNumber)
	if err != nil {
		logger.Warn("inbound order not found", "po_number", poNumber, "error", err)
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, order)
}

func (h *InboundOrderHandler) GetInboundOrderLines(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	lines, err := h.queries.GetInboundOrderLines(r.Context(), id)
	if err != nil {
		logger.Error("failed to get inbound order lines", "inbound_order_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, lines)
}

func (h *InboundOrderHandler) CreateInboundOrderLine(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateInboundOrderLineParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create inbound order line", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	line, err := h.queries.CreateInboundOrderLine(r.Context(), params)
	if err != nil {
		logger.Error("failed to create inbound order line in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("inbound order line created successfully", "line_id", line.ID)
	shared.WriteJSON(w, http.StatusCreated, line)
}

func (h *InboundOrderHandler) UpdateInboundOrderStatus(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.UpdateInboundOrderStatusParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for update inbound order status", "error", err)
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	order, err := h.queries.UpdateInboundOrderStatus(r.Context(), params)
	if err != nil {
		logger.Error("failed to update inbound order status", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("inbound order status updated successfully", "order_id", order.ID)
	shared.WriteJSON(w, http.StatusOK, order)
}

func (h *InboundOrderHandler) CreateReceipt(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateReceiptParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create receipt", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	receipt, err := h.queries.CreateReceipt(r.Context(), params)
	if err != nil {
		logger.Error("failed to create receipt in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("receipt created successfully", "receipt_id", receipt.ID)
	shared.WriteJSON(w, http.StatusCreated, receipt)
}

func (h *InboundOrderHandler) CreateReceiptLine(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateReceiptLineParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create receipt line", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	line, err := h.queries.CreateReceiptLine(r.Context(), params)
	if err != nil {
		logger.Error("failed to create receipt line in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("receipt line created successfully", "line_id", line.ID)
	shared.WriteJSON(w, http.StatusCreated, line)
}

func (h *InboundOrderHandler) DeleteInboundOrder(w http.ResponseWriter, r *http.Request) {
	// logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	_, err = h.queries.DeleteInboundOrder(r.Context(), id)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InboundOrderHandler) IncrementInboundLineReceived(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("line_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	arg := db.IncrementInboundLineReceivedParams{
		ID: id,
	}

	line, err := h.queries.IncrementInboundLineReceived(r.Context(), arg)
	if err != nil {
		logger.Error("failed to increment inbound line received", "line_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("inbound line received incremented successfully", "line_id", line.ID)
	shared.WriteJSON(w, http.StatusOK, line)
}

func (h *InboundOrderHandler) UpdateReceiptStatus(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.UpdateReceiptStatusParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for update receipt status", "error", err)
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	receipt, err := h.queries.UpdateReceiptStatus(r.Context(), params)
	if err != nil {
		logger.Error("failed to update receipt status", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("receipt status updated successfully", "receipt_id", receipt.ID)
	shared.WriteJSON(w, http.StatusOK, receipt)
}

