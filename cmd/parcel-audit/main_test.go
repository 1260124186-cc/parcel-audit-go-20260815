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

func TestRunKeepsRuntimeFailureContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	var output bytes.Buffer
	err := run(ctx, strings.NewReader(`{
		"routes": [{"id": "north", "capacity": 1}]
	}`), &output)

	if err == nil {
		t.Fatal("run() error = nil, want runtime failure")
	}
	if strings.Contains(err.Error(), "validation error:") {
		t.Fatalf("run() error = %v, must not be reported as validation error", err)
	}
	if !strings.Contains(err.Error(), "audit failed: load plan: context canceled") {
		t.Fatalf("run() error = %v, want runtime failure context", err)
	}
}
