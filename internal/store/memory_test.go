package store

import (
	"context"
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

	got.Shipments[0].Labels[0] = "returned-plan-change"
	reloaded, err := memory.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() after changing returned plan error = %v", err)
	}
	if want := []string{"cold"}; !reflect.DeepEqual(want, reloaded.Shipments[0].Labels) {
		t.Fatalf("stored labels after returned plan change = %#v, want %#v", reloaded.Shipments[0].Labels, want)
	}
}
