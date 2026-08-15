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
			{ID: "accepted-after", RouteID: "north", Labels: []string{"bulk"}, HoldUntil: &now},
		},
	}
	auditor := service.NewAuditor(store.NewMemory(plan), func() time.Time { return now })

	report, err := auditor.Audit(context.Background())
	if err != nil {
		t.Fatalf("Audit() error = %v", err)
	}
	if got := assignmentIDs(report); !reflect.DeepEqual([]string{"held", "accepted", "accepted-after"}, got) {
		t.Fatalf("assignment ids = %#v, want [held accepted accepted-after]", got)
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
	if want := []domain.RouteLoad{{RouteID: "north", Used: 2, Capacity: 2}}; !reflect.DeepEqual(want, report.RouteLoads) {
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

// 混合审计场景：同一路线依次处理可发运、被阻止、后续可发运货件。
// 被阻止的货件不得占用容量，后续可发运货件才能在容量耗尽时被正确拒绝，
// 最终路线负载须与真正被接收的货件数量一致。
func TestAuditBlockedShipmentDoesNotConsumeCapacity(t *testing.T) {
	now := time.Date(2026, 8, 15, 9, 0, 0, 0, time.UTC)
	plan := domain.Plan{
		Routes: []domain.Route{{ID: "north", Capacity: 2}},
		Shipments: []domain.Shipment{
			{ID: "s1", RouteID: "north", Labels: []string{"bulk"}, HoldUntil: &now},
			{ID: "s2", RouteID: "north", Labels: []string{"quarantine"}, HoldUntil: &now},
			{ID: "s3", RouteID: "north", Labels: []string{"bulk"}, HoldUntil: &now},
			{ID: "s4", RouteID: "north", Labels: []string{"bulk"}, HoldUntil: &now},
		},
	}
	auditor := service.NewAuditor(store.NewMemory(plan), func() time.Time { return now })

	report, err := auditor.Audit(context.Background())
	if err != nil {
		t.Fatalf("Audit() error = %v", err)
	}
	if got := assignmentIDs(report); !reflect.DeepEqual([]string{"s1", "s3"}, got) {
		t.Fatalf("assignment ids = %#v, want [s1 s3]", got)
	}
	if got := rejectionIDs(report); !reflect.DeepEqual([]string{"s2", "s4"}, got) {
		t.Fatalf("rejection ids = %#v, want [s2 s4]", got)
	}
	if report.Rejections[0].Reason != "shipment is blocked" {
		t.Fatalf("s2 reason = %q, want shipment is blocked", report.Rejections[0].Reason)
	}
	if report.Rejections[1].Reason != "route capacity is unavailable" {
		t.Fatalf("s4 reason = %q, want route capacity is unavailable", report.Rejections[1].Reason)
	}
	if want := []domain.RouteLoad{{RouteID: "north", Used: 2, Capacity: 2}}; !reflect.DeepEqual(want, report.RouteLoads) {
		t.Fatalf("route loads = %#v, want %#v", report.RouteLoads, want)
	}
}

func assignmentIDs(report domain.Report) []string {
	ids := make([]string, 0, len(report.Assignments))
	for _, assignment := range report.Assignments {
		ids = append(ids, assignment.ShipmentID)
	}
	return ids
}

func rejectionIDs(report domain.Report) []string {
	ids := make([]string, 0, len(report.Rejections))
	for _, rejection := range report.Rejections {
		ids = append(ids, rejection.ShipmentID)
	}
	return ids
}
