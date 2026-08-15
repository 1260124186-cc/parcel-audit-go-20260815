package store

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain"
)

func TestMemoryStoreOwnsPlanSnapshot(t *testing.T) {
	plan := domain.Plan{
		Routes:    []domain.Route{{ID: "north", Capacity: 1}},
		Shipments: []domain.Shipment{{ID: "p1", RouteID: "north", Labels: []string{"cold"}}},
	}
	memory := NewMemory(plan)
	plan.Shipments[0].Labels[0] = "changed"

	got, err := memory.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if want := []string{"cold"}; !reflect.DeepEqual(want, got.Shipments[0].Labels) {
		t.Fatalf("stored labels = %#v, want %#v", got.Shipments[0].Labels, want)
	}
}

func TestMemoryStoreLoadRespectsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	memory := NewMemory(domain.Plan{Routes: []domain.Route{{ID: "north", Capacity: 1}}})

	if _, err := memory.Load(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Load() error = %v, want context cancellation", err)
	}
}
