package domain

func ClonePlan(plan Plan) Plan {
	copyPlan := Plan{
		DefaultHoldMinutes: plan.DefaultHoldMinutes,
		Routes:             append([]Route(nil), plan.Routes...),
		Shipments:          make([]Shipment, 0, len(plan.Shipments)),
	}

	for _, shipment := range plan.Shipments {
		copyShipment := shipment
		copyShipment.Labels = append([]string(nil), shipment.Labels...)
		if shipment.HoldUntil != nil {
			holdUntil := *shipment.HoldUntil
			copyShipment.HoldUntil = &holdUntil
		}
		copyPlan.Shipments = append(copyPlan.Shipments, copyShipment)
	}

	return copyPlan
}
