package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/domain"
	"github.com/1260124186-cc/parcel-audit-go-20260815/internal/store"
)

type Auditor struct {
	routes store.RouteStore
	now    func() time.Time
}

func NewAuditor(routes store.RouteStore, now func() time.Time) *Auditor {
	return &Auditor{routes: routes, now: now}
}

func (a *Auditor) Audit(ctx context.Context) (domain.Report, error) {
	if err := ctx.Err(); err != nil {
		return domain.Report{}, err
	}

	plan, err := a.routes.Load(ctx)
	if err != nil {
		return domain.Report{}, fmt.Errorf("load plan: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.Report{}, err
	}
	if err := domain.ValidatePlan(plan); err != nil {
		return domain.Report{}, fmt.Errorf("validate plan: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.Report{}, err
	}

	report := domain.Report{}
	now := a.now().UTC()
	for _, shipment := range plan.Shipments {
		if err := ctx.Err(); err != nil {
			return domain.Report{}, err
		}
		assignment, rejection, err := a.auditShipment(ctx, plan.DefaultHoldMinutes, shipment, now)
		if err != nil {
			return domain.Report{}, err
		}
		if err := ctx.Err(); err != nil {
			return domain.Report{}, err
		}
		if rejection != nil {
			report.Rejections = append(report.Rejections, *rejection)
			continue
		}
		report.Assignments = append(report.Assignments, *assignment)
	}
	report.RouteLoads, err = a.routes.RouteLoads(ctx)
	if err != nil {
		return domain.Report{}, fmt.Errorf("read route loads: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return domain.Report{}, err
	}
	return report, nil
}

func (a *Auditor) auditShipment(ctx context.Context, defaultHoldMinutes int, shipment domain.Shipment, now time.Time) (*domain.Assignment, *domain.Rejection, error) {
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	labels := domain.NormalizeLabels(shipment.Labels)
	readyAt := domain.EffectiveReadyAt(shipment, defaultHoldMinutes, now)
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if readyAt.After(now) {
		return &domain.Assignment{
			ShipmentID: shipment.ID,
			RouteID:    shipment.RouteID,
			State:      "held",
			ReadyAt:    readyAt,
			Labels:     labels,
		}, nil, nil
	}

	finalize, err := a.routes.Reserve(ctx, shipment.RouteID, shipment.ID)
	if err != nil {
		if errors.Is(err, store.ErrRouteUnavailable) {
			return nil, &domain.Rejection{ShipmentID: shipment.ID, Reason: "route capacity is unavailable"}, nil
		}
		return nil, nil, fmt.Errorf("reserve %s: %w", shipment.ID, err)
	}

	committed := false
	defer func() { finalize(committed) }()
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	if domain.HasBlockingLabel(labels) {
		return nil, &domain.Rejection{ShipmentID: shipment.ID, Reason: "shipment is blocked"}, nil
	}
	if err := ctx.Err(); err != nil {
		return nil, nil, err
	}
	committed = true
	return &domain.Assignment{
		ShipmentID: shipment.ID,
		RouteID:    shipment.RouteID,
		State:      "accepted",
		ReadyAt:    readyAt,
		Labels:     labels,
	}, nil, nil
}
