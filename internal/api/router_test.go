package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/L0L1P0P1/wms-rasha-shafa/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var errNoDatabase = errors.New("no database in router test")

type stubDB struct{}

func (stubDB) Exec(context.Context, string, ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, errNoDatabase
}

func (stubDB) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errNoDatabase
}

func (stubDB) QueryRow(context.Context, string, ...any) pgx.Row {
	return errRow{errNoDatabase}
}

type errRow struct{ err error }

func (r errRow) Scan(...any) error { return r.err }

// TestRouterCoversOpenAPIPaths asserts that every path/method pair documented in
// openapi.yaml reaches a handler. ServeMux answers unmatched requests with a
// plain-text 404, while every handler — including the 501 stubs — answers with
// JSON, so a JSON content type proves the route is wired.
func TestRouterCoversOpenAPIPaths(t *testing.T) {
	router := NewRouter(db.New(stubDB{}))

	cases := []struct {
		method string
		path   string
	}{
		{"GET", "/api/ping"},

		{"POST", "/api/v1/skus"},
		{"GET", "/api/v1/skus"},
		{"PATCH", "/api/v1/skus"},
		{"GET", "/api/v1/skus/100001"},
		{"PATCH", "/api/v1/skus/100001"},
		{"DELETE", "/api/v1/skus/100001"},
		{"POST", "/api/v1/skus/100001/packaging-units"},
		{"GET", "/api/v1/skus/100001/packaging-units"},
		{"GET", "/api/v1/packaging-units/100001"},
		{"PATCH", "/api/v1/packaging-units/100001"},
		{"DELETE", "/api/v1/packaging-units/100001"},
		{"GET", "/api/v1/packaging-units/by-barcode/0123456789"},
		{"POST", "/api/v1/skus/100001/lots"},
		{"GET", "/api/v1/skus/100001/lots"},
		{"GET", "/api/v1/skus/100001/lots/LOT-1"},
		{"GET", "/api/v1/lots/1"},
		{"PATCH", "/api/v1/lots/1/status"},

		{"POST", "/api/v1/storage-nodes"},
		{"GET", "/api/v1/storage-nodes"},
		{"GET", "/api/v1/storage-nodes/1"},
		{"PATCH", "/api/v1/storage-nodes/1"},
		{"DELETE", "/api/v1/storage-nodes/1"},
		{"GET", "/api/v1/storage-nodes/by-code/WH-1"},
		{"GET", "/api/v1/storage-nodes/root/children"},
		{"GET", "/api/v1/storage-nodes/1/children"},
		{"GET", "/api/v1/storage-nodes/1/subtree"},
		{"PUT", "/api/v1/storage-nodes/1/parent"},
		{"POST", "/api/v1/storage-nodes/1/relocations"},

		{"GET", "/api/v1/inventory/balances/node/1"},
		{"GET", "/api/v1/inventory/balances/reservable"},
		{"POST", "/api/v1/inventory/balances/lock"},
		{"PATCH", "/api/v1/inventory/balances/1/allocation"},
		{"POST", "/api/v1/inventory/balances/project/destination-addition"},
		{"POST", "/api/v1/inventory/balances/project/source-deduction"},
		{"POST", "/api/v1/inventory/allocations"},
		{"GET", "/api/v1/inventory/allocations/order-line/1"},
		{"PATCH", "/api/v1/inventory/allocations/1/picked"},
		{"DELETE", "/api/v1/inventory/allocations/1"},
		{"POST", "/api/v1/inventory/movements"},
		{"GET", "/api/v1/inventory/movements/node/1"},
		{"GET", "/api/v1/inventory/movements/sku/100001"},

		{"POST", "/api/v1/inbound-orders"},
		{"GET", "/api/v1/inbound-orders"},
		{"GET", "/api/v1/inbound-orders/1"},
		{"DELETE", "/api/v1/inbound-orders/1"},
		{"GET", "/api/v1/inbound-orders/by-po-number/PO-1"},
		{"POST", "/api/v1/inbound-orders/1/lines"},
		{"GET", "/api/v1/inbound-orders/1/lines"},
		{"POST", "/api/v1/inbound-orders/1/lines/1/received"},
		{"PATCH", "/api/v1/inbound-orders/1/status"},
		{"POST", "/api/v1/inbound-orders/1/receipts"},
		{"POST", "/api/v1/receipts/1/lines"},
		{"PATCH", "/api/v1/receipts/1/status"},

		{"POST", "/api/v1/outbound-orders"},
		{"GET", "/api/v1/outbound-orders"},
		{"GET", "/api/v1/outbound-orders/1"},
		{"DELETE", "/api/v1/outbound-orders/1"},
		{"GET", "/api/v1/outbound-orders/by-order-number/SO-1"},
		{"POST", "/api/v1/outbound-orders/1/lines"},
		{"GET", "/api/v1/outbound-orders/1/lines"},
		{"POST", "/api/v1/outbound-orders/1/lines/1/fulfillment"},
		{"PATCH", "/api/v1/outbound-orders/1/status"},

		{"POST", "/api/v1/handling-units"},
		{"POST", "/api/v1/handling-units/1/contents"},
		{"PUT", "/api/v1/handling-units/1/shipment"},
		{"POST", "/api/v1/shipments"},

		{"POST", "/api/v1/tasks"},
		{"GET", "/api/v1/tasks"},
		{"GET", "/api/v1/tasks/1"},
		{"POST", "/api/v1/tasks/1/assign"},
		{"POST", "/api/v1/tasks/1/complete"},
		{"POST", "/api/v1/tasks/1/cancel"},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader("{}"))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if got := rec.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Errorf("%s %s: not routed to a handler (status %d, content-type %q)",
				tc.method, tc.path, rec.Code, got)
		}
	}
}

// TestRouterRejectsUnknownActions guards the catch-all `.../{id}/{action}`
// patterns used to work around ServeMux pattern conflicts.
func TestRouterRejectsUnknownActions(t *testing.T) {
	router := NewRouter(db.New(stubDB{}))

	for _, path := range []string{
		"/api/v1/storage-nodes/1/bogus",
		"/api/v1/inbound-orders/1/bogus",
		"/api/v1/outbound-orders/1/bogus",
	} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))

		if rec.Code != http.StatusNotFound {
			t.Errorf("GET %s: got %d, want 404", path, rec.Code)
		}
	}
}
