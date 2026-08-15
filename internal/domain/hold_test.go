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
