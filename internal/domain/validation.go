package domain

import "fmt"

type ValidationError struct {
	Field  string
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

func ValidatePlan(plan Plan) error {
	if plan.DefaultHoldMinutes < 0 {
		return &ValidationError{Field: "default_hold_minutes", Reason: "must not be negative"}
	}

	routes := make(map[string]struct{}, len(plan.Routes))
	for _, route := range plan.Routes {
		if route.ID == "" {
			return &ValidationError{Field: "routes", Reason: "route id is required"}
		}
		if route.Capacity < 1 {
			return &ValidationError{Field: "routes." + route.ID, Reason: "capacity must be positive"}
		}
		if _, exists := routes[route.ID]; exists {
			return &ValidationError{Field: "routes." + route.ID, Reason: "route id must be unique"}
		}
		routes[route.ID] = struct{}{}
	}

	for _, shipment := range plan.Shipments {
		if shipment.ID == "" {
			return &ValidationError{Field: "shipments", Reason: "shipment id is required"}
		}
		if shipment.RouteID == "" {
			return &ValidationError{Field: "shipments." + shipment.ID, Reason: "route id is required"}
		}
	}

	return nil
}
