package service_test

import (
	"context"
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

	if _, err := auditor.Audit(ctx); err == nil {
		t.Fatal("Audit() error = nil, want cancelled context error")
	}
}

func assignmentIDs(report domain.Report) []string {
	ids := make([]string, 0, len(report.Assignments))
	for _, assignment := range report.Assignments {
		ids = append(ids, assignment.ShipmentID)
	}
	return ids
}
