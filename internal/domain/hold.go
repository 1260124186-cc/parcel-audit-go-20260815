package domain

import "time"

func EffectiveReadyAt(shipment Shipment, defaultHoldMinutes int, now time.Time) time.Time {
	// 显式保留截止时间优先；否则遵循计划的默认保留时长
	if shipment.HoldUntil != nil {
		return shipment.HoldUntil.UTC()
	}
	return now.Add(time.Duration(defaultHoldMinutes) * time.Minute).UTC()
}
