package outbound_orders

import (
	"net/http"
	"strconv"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/ctxlog"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
)

type OutboundOrderHandler struct {
	queries *db.Queries
}

func NewOutboundOrderHandler(queries *db.Queries) *OutboundOrderHandler {
	return &OutboundOrderHandler{queries: queries}
}

func (h *OutboundOrderHandler) CreateHandlingUnit(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateHandlingUnitParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create handling unit", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	unit, err := h.queries.CreateHandlingUnit(r.Context(), params)
	if err != nil {
		logger.Error("failed to create handling unit in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, unit)
}

func (h *OutboundOrderHandler) CreateOutboundOrder(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateOutboundOrderParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create outbound order", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	order, err := h.queries.CreateOutboundOrder(r.Context(), params)
	if err != nil {
		logger.Error("failed to create outbound order in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("outbound order created successfully", "order_id", order.ID)
	shared.WriteJSON(w, http.StatusCreated, order)
}

func (h *OutboundOrderHandler) CreateOutboundOrderLine(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateOutboundOrderLineParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create outbound order line", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	line, err := h.queries.CreateOutboundOrderLine(r.Context(), params)
	if err != nil {
		logger.Error("failed to create outbound order line in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("outbound order line created successfully", "line_id", line.ID)
	shared.WriteJSON(w, http.StatusCreated, line)
}

func (h *OutboundOrderHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.CreateShipmentParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for create shipment", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	shipment, err := h.queries.CreateShipment(r.Context(), params)
	if err != nil {
		logger.Error("failed to create shipment in database", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("shipment created successfully", "shipment_id", shipment.ID)
	shared.WriteJSON(w, http.StatusCreated, shipment)
}

func (h *OutboundOrderHandler) DeleteOutboundOrder(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	_, err = h.queries.DeleteOutboundOrder(r.Context(), id)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *OutboundOrderHandler) GetOutboundOrderByID(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	order, err := h.queries.GetOutboundOrderByID(r.Context(), id)
	if err != nil {
		logger.Warn("outbound order not found", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, order)
}

func (h *OutboundOrderHandler) GetOutboundOrderByNumber(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	orderNumber := r.PathValue("order_number")
	order, err := h.queries.GetOutboundOrderByNumber(r.Context(), orderNumber)
	if err != nil {
		logger.Warn("outbound order not found", "order_number", orderNumber, "error", err)
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, order)
}

func (h *OutboundOrderHandler) GetOutboundOrderLines(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	lines, err := h.queries.GetOutboundOrderLines(r.Context(), id)
	if err != nil {
		logger.Error("failed to get outbound order lines", "order_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, lines)
}

func (h *OutboundOrderHandler) IncrementOrderLineFulfillment(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("line_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.IncrementOrderLineFulfillmentParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for increment order line fulfillment", "error", err)
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	line, err := h.queries.IncrementOrderLineFulfillment(r.Context(), params)
	if err != nil {
		logger.Error("failed to increment order line fulfillment", "line_id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("order line fulfillment incremented successfully", "line_id", line.ID)
	shared.WriteJSON(w, http.StatusOK, line)
}

func (h *OutboundOrderHandler) ListOrdersByStatus(w http.ResponseWriter, r *http.Request) {
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

	orders, err := h.queries.ListOrdersByStatus(r.Context(), db.ListOrdersByStatusParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to list orders by status", "limit", limit, "offset", offset, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, orders)
}

func (h *OutboundOrderHandler) PackHandlingUnitContent(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	var params db.PackHandlingUnitContentParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for pack handling unit content", "error", err)
		shared.ErrorJSON(w, err)
		return
	}

	_, err := h.queries.PackHandlingUnitContent(r.Context(), params)
	if err != nil {
		logger.Error("failed to pack handling unit content", "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OutboundOrderHandler) UpdateOutboundOrderStatus(w http.ResponseWriter, r *http.Request) {
	logger := ctxlog.FromContext(r.Context())

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.UpdateOutboundOrderStatusParams
	if err := shared.ReadJSON(r, &params); err != nil {
		logger.Warn("failed to decode request body for update outbound order status", "error", err)
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	order, err := h.queries.UpdateOutboundOrderStatus(r.Context(), params)
	if err != nil {
		logger.Error("failed to update outbound order status", "id", id, "error", err)
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	logger.Info("outbound order status updated successfully", "order_id", order.ID)
	shared.WriteJSON(w, http.StatusOK, order)
}

