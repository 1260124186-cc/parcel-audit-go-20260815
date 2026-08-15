package domain

import (
	"reflect"
	"testing"
)

func TestNormalizeLabelsDoesNotMutateInput(t *testing.T) {
	input := []string{" Cold ", "fragile", "cold"}
	got := NormalizeLabels(input)

	if want := []string{"cold", "fragile"}; !reflect.DeepEqual(want, got) {
		t.Fatalf("normalized labels = %#v, want %#v", got, want)
	}
	if want := []string{" Cold ", "fragile", "cold"}; !reflect.DeepEqual(want, input) {
		t.Fatalf("input labels changed = %#v, want %#v", input, want)
	}
}
