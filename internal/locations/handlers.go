package locations

import (
	"net/http"
	"strconv"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
	"github.com/jackc/pgx/v5/pgtype"
)

type LocationHandler struct {
	queries *db.Queries
}

func NewLocationHandler(queries *db.Queries) *LocationHandler {
	return &LocationHandler{queries: queries}
}

func (h *LocationHandler) CreateStorageNode(w http.ResponseWriter, r *http.Request) {
	var params db.CreateStorageNodeParams
	if err := shared.ReadJSON(r, &params); err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	node, err := h.queries.CreateStorageNode(r.Context(), params)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusCreated, node)
}

func (h *LocationHandler) GetStorageNode(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	node, err := h.queries.GetStorageNodeByID(r.Context(), id)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, node)
}

func (h *LocationHandler) GetStorageNodeByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	node, err := h.queries.GetStorageNodeByCode(r.Context(), code)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	shared.WriteJSON(w, http.StatusOK, node)
}

func (h *LocationHandler) ListDirectChildNodes(w http.ResponseWriter, r *http.Request) {
	parentIdStr := r.PathValue("id")
	var parentId pgtype.Int8
	if parentIdStr != "" && parentIdStr != "root" {
		id, err := strconv.ParseInt(parentIdStr, 10, 64)
		if err != nil {
			shared.ErrorJSON(w, err)
			return
		}
		parentId = pgtype.Int8{Int64: id, Valid: true}
	} else {
		parentId = pgtype.Int8{Valid: false}
	}

	nodes, err := h.queries.ListDirectChildNodes(r.Context(), parentId)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, nodes)
}

func (h *LocationHandler) ListSubtreeNodes(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	nodes, err := h.queries.ListSubtreeNodes(r.Context(), id)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, nodes)
}

func (h *LocationHandler) UpdateStorageNodeDetails(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.UpdateStorageNodeDetailsParams
	if err := shared.ReadJSON(r, &params); err != nil {
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	node, err := h.queries.UpdateStorageNodeDetails(r.Context(), params)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, node)
}

func (h *LocationHandler) ReparentStorageNode(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	var params db.ReparentStorageNodeParams
	if err := shared.ReadJSON(r, &params); err != nil {
		shared.ErrorJSON(w, err)
		return
	}
	params.ID = id

	node, err := h.queries.ReparentStorageNode(r.Context(), params)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	shared.WriteJSON(w, http.StatusOK, node)
}

func (h *LocationHandler) DeleteStorageNode(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared.ErrorJSON(w, err)
		return
	}

	_, err = h.queries.DeleteStorageNode(r.Context(), id)
	if err != nil {
		shared.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
