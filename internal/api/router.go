package api

import (
	"net/http"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/inbound_orders"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/inventory"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/locations"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/outbound_orders"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/shared"
	"github.com/L0L1P0P1/wms-rasha-shafa/internal/skus"
)

func NewRouter(queries *db.Queries) http.Handler {
	skuHandler := skus.NewSkuHandler(queries)
	locationHandler := locations.NewLocationHandler(queries)
	inventoryHandler := inventory.NewInventoryHandler(queries)
	inboundOrderHandler := inbound_orders.NewInboundOrderHandler(queries)
	outboundOrderHandler := outbound_orders.NewOutboundOrderHandler(queries)

	mux := http.NewServeMux()

	registerHealthRoutes(mux)
	registerSkuRoutes(mux, skuHandler)
	registerStorageNodeRoutes(mux, locationHandler)
	registerInventoryRoutes(mux, inventoryHandler)
	registerInboundOrderRoutes(mux, inboundOrderHandler)
	registerOutboundOrderRoutes(mux, outboundOrderHandler)
	registerHandlingUnitRoutes(mux, outboundOrderHandler)
	registerTaskRoutes(mux)

	return mux
}

func registerHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/ping", ping)
}

func registerSkuRoutes(mux *http.ServeMux, h *skus.SkuHandler) {
	// SKUs
	mux.HandleFunc("POST /api/v1/skus", h.CreateSKU)
	mux.HandleFunc("GET /api/v1/skus", h.ListSKUs)
	mux.HandleFunc("GET /api/v1/skus/{id}", h.GetSKU)
	mux.HandleFunc("PATCH /api/v1/skus/{id}", h.UpdateSKU)
	mux.HandleFunc("DELETE /api/v1/skus/{id}", h.DeleteSKU)
	mux.HandleFunc("PATCH /api/v1/skus", notImplemented("searchSkusByName"))

	// Packaging units nested under a SKU
	mux.HandleFunc("POST /api/v1/skus/{sku_id}/packaging-units", h.CreatePackagingUnit)
	mux.HandleFunc("GET /api/v1/skus/{sku_id}/packaging-units", h.ListPackagingUnitsBySKU)

	// Packaging units addressed directly
	mux.HandleFunc("GET /api/v1/packaging-units/{id}", notImplemented("getPackagingUnitById"))
	mux.HandleFunc("PATCH /api/v1/packaging-units/{id}", notImplemented("updatePackagingUnit"))
	mux.HandleFunc("DELETE /api/v1/packaging-units/{id}", notImplemented("deletePackagingUnit"))
	mux.HandleFunc("GET /api/v1/packaging-units/by-barcode/{barcode}", notImplemented("getPackagingUnitByBarcode"))

	// Lots nested under a SKU
	mux.HandleFunc("POST /api/v1/skus/{sku_id}/lots", h.CreateLot)
	mux.HandleFunc("GET /api/v1/skus/{sku_id}/lots", h.ListLotsBySKU)
	mux.HandleFunc("GET /api/v1/skus/{sku_id}/lots/{lot_number}", notImplemented("getLotByNumber"))

	// Lots addressed directly
	mux.HandleFunc("GET /api/v1/lots/{id}", notImplemented("getLotById"))
	mux.HandleFunc("PATCH /api/v1/lots/{id}/status", notImplemented("updateLotStatus"))
}

func registerStorageNodeRoutes(mux *http.ServeMux, h *locations.LocationHandler) {
	mux.HandleFunc("POST /api/v1/storage-nodes", h.CreateStorageNode)
	mux.HandleFunc("GET /api/v1/storage-nodes", h.ListDirectChildNodes)

	// Literal segments outrank wildcards in ServeMux, so `root` and `by-code`
	// are never captured as {id}.
	mux.HandleFunc("GET /api/v1/storage-nodes/root/children", h.ListDirectChildNodes)
	mux.HandleFunc("GET /api/v1/storage-nodes/by-code/{code}", h.GetStorageNodeByCode)
	mux.HandleFunc("GET /api/v1/storage-nodes/{id}", h.GetStorageNode)
	mux.HandleFunc("PATCH /api/v1/storage-nodes/{id}", h.UpdateStorageNodeDetails)
	mux.HandleFunc("DELETE /api/v1/storage-nodes/{id}", h.DeleteStorageNode)
	mux.HandleFunc("PUT /api/v1/storage-nodes/{id}/parent", h.ReparentStorageNode)
	mux.HandleFunc("POST /api/v1/storage-nodes/{id}/relocations", notImplemented("recordNodeRelocation"))

	// Fallback for `GET /api/v1/storage-nodes/{id}/children` and
	// `GET /api/v1/storage-nodes/{id}/subtree`, which cannot be registered
	// directly next to `by-code/{code}`: `/api/v1/storage-nodes/by-code/children`
	// would match both patterns and neither is more specific, which makes
	// ServeMux panic at registration time.
	mux.Handle("GET /api/v1/storage-nodes/{id}/{action}", byTrailingAction(map[string]http.HandlerFunc{
		"children": h.ListDirectChildNodes,
		"subtree":  h.ListSubtreeNodes,
	}))
}

func registerInventoryRoutes(mux *http.ServeMux, h *inventory.InventoryHandler) {
	// Balances
	mux.HandleFunc("GET /api/v1/inventory/balances/node/{node_id}", h.ListBalancesByNode)
	mux.HandleFunc("GET /api/v1/inventory/balances/reservable", h.FindReservableBalances)
	mux.HandleFunc("POST /api/v1/inventory/balances/lock", h.LockBalanceForUpdate)
	mux.HandleFunc("PATCH /api/v1/inventory/balances/{id}/allocation", h.UpdateBalanceAllocation)
	mux.HandleFunc("POST /api/v1/inventory/balances/project/destination-addition", h.ProjectDestinationAddition)
	mux.HandleFunc("POST /api/v1/inventory/balances/project/source-deduction", h.ProjectSourceDeduction)

	// Allocations
	mux.HandleFunc("POST /api/v1/inventory/allocations", h.CreateAllocation)
	mux.HandleFunc("GET /api/v1/inventory/allocations/order-line/{order_line_id}", h.GetAllocationsByOrderLine)
	mux.HandleFunc("PATCH /api/v1/inventory/allocations/{id}/picked", h.IncrementAllocationPicked)
	mux.HandleFunc("DELETE /api/v1/inventory/allocations/{id}", h.DeleteAllocation)

	// Movements
	mux.HandleFunc("POST /api/v1/inventory/movements", h.InsertMovement)
	mux.HandleFunc("GET /api/v1/inventory/movements/node/{node_id}", h.ListMovementsByNode)
	mux.HandleFunc("GET /api/v1/inventory/movements/sku/{sku_id}", h.ListMovementsBySKU)
}

func registerInboundOrderRoutes(mux *http.ServeMux, h *inbound_orders.InboundOrderHandler) {
	mux.HandleFunc("POST /api/v1/inbound-orders", h.CreateInboundOrder)
	mux.HandleFunc("GET /api/v1/inbound-orders", notImplemented("listInboundOrders"))
	mux.HandleFunc("GET /api/v1/inbound-orders/by-po-number/{po_number}", h.GetInboundOrderByPONumber)
	mux.HandleFunc("GET /api/v1/inbound-orders/{id}", h.GetInboundOrderByID)
	mux.HandleFunc("DELETE /api/v1/inbound-orders/{id}", h.DeleteInboundOrder)
	mux.HandleFunc("POST /api/v1/inbound-orders/{id}/lines", h.CreateInboundOrderLine)
	mux.HandleFunc("POST /api/v1/inbound-orders/{id}/lines/{line_id}/received", h.IncrementInboundLineReceived)
	mux.HandleFunc("PATCH /api/v1/inbound-orders/{id}/status", h.UpdateInboundOrderStatus)
	mux.HandleFunc("POST /api/v1/inbound-orders/{id}/receipts", h.CreateReceipt)

	// Fallback for `GET /api/v1/inbound-orders/{id}/lines`, which collides with
	// `by-po-number/{po_number}` (see registerStorageNodeRoutes).
	mux.Handle("GET /api/v1/inbound-orders/{id}/{action}", byTrailingAction(map[string]http.HandlerFunc{
		"lines": h.GetInboundOrderLines,
	}))

	// Receipts addressed directly
	mux.HandleFunc("POST /api/v1/receipts/{id}/lines", h.CreateReceiptLine)
	mux.HandleFunc("PATCH /api/v1/receipts/{id}/status", h.UpdateReceiptStatus)
}

func registerOutboundOrderRoutes(mux *http.ServeMux, h *outbound_orders.OutboundOrderHandler) {
	mux.HandleFunc("POST /api/v1/outbound-orders", h.CreateOutboundOrder)
	mux.HandleFunc("GET /api/v1/outbound-orders", h.ListOrdersByStatus)
	mux.HandleFunc("GET /api/v1/outbound-orders/by-order-number/{order_number}", h.GetOutboundOrderByNumber)
	mux.HandleFunc("GET /api/v1/outbound-orders/{id}", h.GetOutboundOrderByID)
	mux.HandleFunc("DELETE /api/v1/outbound-orders/{id}", h.DeleteOutboundOrder)
	mux.HandleFunc("POST /api/v1/outbound-orders/{id}/lines", h.CreateOutboundOrderLine)
	mux.HandleFunc("POST /api/v1/outbound-orders/{id}/lines/{line_id}/fulfillment", h.IncrementOrderLineFulfillment)
	mux.HandleFunc("PATCH /api/v1/outbound-orders/{id}/status", h.UpdateOutboundOrderStatus)

	// Fallback for `GET /api/v1/outbound-orders/{id}/lines`, which collides with
	// `by-order-number/{order_number}` (see registerStorageNodeRoutes).
	mux.Handle("GET /api/v1/outbound-orders/{id}/{action}", byTrailingAction(map[string]http.HandlerFunc{
		"lines": h.GetOutboundOrderLines,
	}))
}

func registerHandlingUnitRoutes(mux *http.ServeMux, h *outbound_orders.OutboundOrderHandler) {
	mux.HandleFunc("POST /api/v1/handling-units", h.CreateHandlingUnit)
	mux.HandleFunc("POST /api/v1/handling-units/{handling_unit_id}/contents", h.PackHandlingUnitContent)
	mux.HandleFunc("POST /api/v1/shipments", h.CreateShipment)
	mux.HandleFunc("PUT /api/v1/handling-units/{handling_unit_id}/shipment", notImplemented("assignHandlingUnitToShipment"))
}

func registerTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/tasks", notImplemented("createWarehouseTask"))
	mux.HandleFunc("GET /api/v1/tasks", notImplemented("listPendingTasks"))
	mux.HandleFunc("GET /api/v1/tasks/{id}", notImplemented("getWarehouseTaskById"))
	mux.HandleFunc("POST /api/v1/tasks/{id}/assign", notImplemented("assignTask"))
	mux.HandleFunc("POST /api/v1/tasks/{id}/complete", notImplemented("completeTask"))
	mux.HandleFunc("POST /api/v1/tasks/{id}/cancel", notImplemented("cancelTask"))
}

func notImplemented(operationID string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shared.ErrorJSON(w, errNotImplemented{OperationID: operationID}, http.StatusNotImplemented)
	}
}

// byTrailingAction serves routes registered as `.../{id}/{action}` by dispatching
// on the trailing literal segment. It exists only for the three collections where
// a literal-prefix route (`by-code`, `by-po-number`, `by-order-number`) sits next
// to a sibling-with-action route (`{id}/children`, `{id}/lines`); those two shapes
// overlap and cannot both be registered, because neither pattern is more specific
// than the other and ServeMux panics. An unknown action is a 404.
func byTrailingAction(actions map[string]http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler, ok := actions[r.PathValue("action")]
		if !ok {
			shared.ErrorJSON(w, errNoRoute{Path: r.URL.Path}, http.StatusNotFound)
			return
		}

		handler(w, r)
	}
}

type errNotImplemented struct {
	OperationID string
}

func (e errNotImplemented) Error() string {
	return "not implemented: " + e.OperationID
}

type errNoRoute struct {
	Path string
}

func (e errNoRoute) Error() string {
	return "no route for " + e.Path
}
