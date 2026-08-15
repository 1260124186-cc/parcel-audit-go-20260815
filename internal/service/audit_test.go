package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain"
	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/service"
	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/store"
)

func TestAuditAppliesHoldAndTracksCommittedReservations(t *testing.T) {
	now := time.Date(2026, 8, 15, 9, 0, 0, 0, time.UTC)
	plan := domain.Plan{
		DefaultHoldMinutes: 15,
		Routes:             []domain.Route{{ID: "north", Capacity: 2}},
		Shipments: []domain.Shipment{
			{ID: "held", RouteID: "north", Labels: []string{"Cold"}},
			{ID: "accepted", RouteID: "north", Labels: []string{"fragile", "fragile"}, HoldUntil: &now},
			{ID: "blocked", RouteID: "north", Labels: []string{"quarantine"}, HoldUntil: &now},
		},
	}
	auditor := service.NewAuditor(store.NewMemory(plan), func() time.Time { return now })

	report, err := auditor.Audit(context.Background())
	if err != nil {
		t.Fatalf("Audit() error = %v", err)
	}
	if got := assignmentIDs(report); !reflect.DeepEqual([]string{"held", "accepted"}, got) {
		t.Fatalf("assignment ids = %#v, want [held accepted]", got)
	}
	if report.Assignments[0].State != "held" || !report.Assignments[0].ReadyAt.Equal(now.Add(15*time.Minute)) {
		t.Fatalf("missing hold was not applied: %#v", report.Assignments[0])
	}
	if got := report.Assignments[1].Labels; !reflect.DeepEqual([]string{"fragile"}, got) {
		t.Fatalf("normalized labels = %#v, want [fragile]", got)
	}
	if want := []domain.Rejection{{ShipmentID: "blocked", Reason: "shipment is blocked"}}; !reflect.DeepEqual(want, report.Rejections) {
		t.Fatalf("rejections = %#v, want %#v", report.Rejections, want)
	}
	if want := []domain.RouteLoad{{RouteID: "north", Used: 1, Capacity: 2}}; !reflect.DeepEqual(want, report.RouteLoads) {
		t.Fatalf("route loads = %#v, want %#v", report.RouteLoads, want)
	}
}

func TestAuditRespectsCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	plan := domain.Plan{Routes: []domain.Route{{ID: "north", Capacity: 1}}}
	auditor := service.NewAuditor(store.NewMemory(plan), time.Now)

	if _, err := auditor.Audit(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Audit() error = %v, want cancelled context error", err)
	}
}

func TestAuditPassesCallerContextToStorage(t *testing.T) {
	plan := domain.Plan{Routes: []domain.Route{{ID: "north", Capacity: 1}}}
	routes := &contextRecordingStore{RouteStore: store.NewMemory(plan)}
	ctx := context.WithValue(context.Background(), contextKey{}, "audit")
	auditor := service.NewAuditor(routes, time.Now)

	if _, err := auditor.Audit(ctx); err != nil {
		t.Fatalf("Audit() error = %v", err)
	}
	if routes.loadContext != ctx {
		t.Fatal("Load() did not receive the caller context")
	}
	if routes.routeLoadsContext != ctx {
		t.Fatal("RouteLoads() did not receive the caller context")
	}
}

func TestAuditStopsAndReleasesReservationWhenCancelledDuringReserve(t *testing.T) {
	plan := domain.Plan{
		Routes: []domain.Route{{ID: "north", Capacity: 2}},
		Shipments: []domain.Shipment{
			{ID: "p1", RouteID: "north"},
			{ID: "p2", RouteID: "north"},
		},
	}
	ctx, cancel := context.WithCancel(context.Background())
	backingStore := store.NewMemory(plan)
	routes := &cancellingReserveStore{
		RouteStore: backingStore,
		cancel:     cancel,
	}
	auditor := service.NewAuditor(routes, time.Now)

	if _, err := auditor.Audit(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("Audit() error = %v, want cancelled context error", err)
	}
	if routes.reserveCalls != 1 {
		t.Fatalf("Reserve() calls = %d, want 1", routes.reserveCalls)
	}
	if routes.routeLoadsCalls != 0 {
		t.Fatalf("RouteLoads() calls = %d, want 0", routes.routeLoadsCalls)
	}

	loads, err := backingStore.RouteLoads(context.Background())
	if err != nil {
		t.Fatalf("RouteLoads() error = %v", err)
	}
	if want := []domain.RouteLoad{{RouteID: "north", Used: 0, Capacity: 2}}; !reflect.DeepEqual(loads, want) {
		t.Fatalf("route loads = %#v, want %#v", loads, want)
	}
}

type cancellingReserveStore struct {
	store.RouteStore
	cancel          context.CancelFunc
	reserveCalls    int
	routeLoadsCalls int
}

func (s *cancellingReserveStore) Reserve(ctx context.Context, routeID, shipmentID string) (func(bool), error) {
	s.reserveCalls++
	finalize, err := s.RouteStore.Reserve(ctx, routeID, shipmentID)
	s.cancel()
	return finalize, err
}

func (s *cancellingReserveStore) RouteLoads(ctx context.Context) ([]domain.RouteLoad, error) {
	s.routeLoadsCalls++
	return s.RouteStore.RouteLoads(ctx)
}

type contextKey struct{}

type contextRecordingStore struct {
	store.RouteStore
	loadContext       context.Context
	routeLoadsContext context.Context
}

func (s *contextRecordingStore) Load(ctx context.Context) (domain.Plan, error) {
	s.loadContext = ctx
	return s.RouteStore.Load(ctx)
}

func (s *contextRecordingStore) RouteLoads(ctx context.Context) ([]domain.RouteLoad, error) {
	s.routeLoadsContext = ctx
	return s.RouteStore.RouteLoads(ctx)
}

func assignmentIDs(report domain.Report) []string {
	ids := make([]string, 0, len(report.Assignments))
	for _, assignment := range report.Assignments {
		ids = append(ids, assignment.ShipmentID)
	}
	return ids
}
