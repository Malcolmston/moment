package moment

import "testing"

func TestDurationToJSONAndValueOf(t *testing.T) {
	d := NewDuration(90, Minute)
	if got, want := d.ToJSON(), "PT1H30M"; got != want {
		t.Errorf("ToJSON = %q, want %q", got, want)
	}
	if got, want := d.ValueOf(), int64(90*60*1000); got != want {
		t.Errorf("ValueOf = %d, want %d", got, want)
	}
	if got := NewDuration(-2, Second).ValueOf(); got != -2000 {
		t.Errorf("negative ValueOf = %d, want -2000", got)
	}
}

func TestNewDurationFromObject(t *testing.T) {
	d := NewDurationFromObject(map[string]int{"hours": 2, "minutes": 30})
	if got := d.AsMinutes(); got != 150 {
		t.Errorf("AsMinutes = %v, want 150", got)
	}
	// Aliases and multiple components.
	d2 := NewDurationFromObject(map[string]int{"y": 1, "M": 2, "d": 10})
	if got := d2.Years(); got != 1 {
		t.Errorf("Years = %d, want 1", got)
	}
	if got := d2.Months(); got != 2 {
		t.Errorf("Months = %d, want 2", got)
	}
	if got := d2.Days(); got != 10 {
		t.Errorf("Days = %d, want 10", got)
	}
	// Unknown keys ignored; empty map yields zero.
	if got := NewDurationFromObject(map[string]int{"parsec": 5}).ValueOf(); got != 0 {
		t.Errorf("unknown key should be ignored, got %d", got)
	}
	if got := NewDurationFromObject(nil).ValueOf(); got != 0 {
		t.Errorf("nil map should be zero, got %d", got)
	}
}

func BenchmarkNewDurationFromObject(b *testing.B) {
	obj := map[string]int{"years": 1, "months": 2, "days": 10, "hours": 5, "minutes": 30}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NewDurationFromObject(obj)
	}
}
