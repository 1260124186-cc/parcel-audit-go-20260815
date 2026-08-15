package domain

import (
	"testing"
	"time"
)

func TestEffectiveReadyAtUsesPlanDefault(t *testing.T) {
	now := time.Date(2026, 8, 15, 9, 0, 0, 0, time.UTC)

	got := EffectiveReadyAt(Shipment{}, 20, now)
	want := now.Add(20 * time.Minute)
	if !got.Equal(want) {
		t.Fatalf("EffectiveReadyAt() = %s, want %s", got, want)
	}
}

func TestEffectiveReadyAtPrefersShipmentHoldUntil(t *testing.T) {
	now := time.Date(2026, 8, 15, 9, 0, 0, 0, time.UTC)
	holdUntil := now.Add(5 * time.Minute)

	got := EffectiveReadyAt(Shipment{HoldUntil: &holdUntil}, 20, now)
	if !got.Equal(holdUntil) {
		t.Fatalf("EffectiveReadyAt() = %s, want shipment hold until %s", got, holdUntil)
	}
}
