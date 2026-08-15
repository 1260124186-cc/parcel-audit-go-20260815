package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestRunReportsValidationError(t *testing.T) {
	var output bytes.Buffer
	err := run(context.Background(), strings.NewReader(`{
		"routes": [{"id": "north", "capacity": 1}],
		"shipments": [{"id": "", "route_id": "north"}]
	}`), &output)

	if err == nil || !strings.Contains(err.Error(), "validation error: shipments: shipment id is required") {
		t.Fatalf("run() error = %v, want validation message", err)
	}
}
