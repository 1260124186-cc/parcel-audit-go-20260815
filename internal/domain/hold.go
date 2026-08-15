package domain

import "time"

func EffectiveReadyAt(shipment Shipment, defaultHoldMinutes int, now time.Time) time.Time {
	if shipment.HoldUntil != nil {
		return shipment.HoldUntil.UTC()
	}
	return now.UTC().Add(time.Duration(defaultHoldMinutes) * time.Minute)
}
