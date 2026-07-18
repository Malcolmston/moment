package moment

import (
	"testing"
	"time"
)

func TestComponentSetters(t *testing.T) {
	base := DateTime(2017, 7, 14, 2, 40, 9, 0, time.UTC)
	if got := base.SetYear(2020).Year(); got != 2020 {
		t.Errorf("SetYear = %d", got)
	}
	if got := base.SetMonth(time.December).Month(); got != time.December {
		t.Errorf("SetMonth = %v", got)
	}
	if got := base.SetQuarter(4).Month(); got != time.October {
		t.Errorf("SetQuarter month = %v, want October", got)
	}
	if got := base.SetDate(1).Date(); got != 1 {
		t.Errorf("SetDate = %d", got)
	}
	if got := base.SetDayOfYear(1); got.Month() != time.January || got.Date() != 1 {
		t.Errorf("SetDayOfYear(1) = %v", got.ISO())
	}
	if got := base.SetHour(23).Hour(); got != 23 {
		t.Errorf("SetHour = %d", got)
	}
	if got := base.SetMinute(5).Minute(); got != 5 {
		t.Errorf("SetMinute = %d", got)
	}
	if got := base.SetSecond(0).Second(); got != 0 {
		t.Errorf("SetSecond = %d", got)
	}
	if got := base.SetMillisecond(500).Millisecond(); got != 500 {
		t.Errorf("SetMillisecond = %d", got)
	}
	// Immutability: base must be unchanged.
	if base.Year() != 2017 || base.Hour() != 2 {
		t.Errorf("setters mutated the receiver")
	}
}

func TestWeekdaySetters(t *testing.T) {
	// 2017-07-14 is a Friday (weekday 5, ISO weekday 5).
	base := DateTime(2017, 7, 14, 12, 0, 0, 0, time.UTC)

	// SetWeekday stays within the Sunday-based week: Monday of that week is
	// 2017-07-10.
	if got := base.SetWeekday(time.Monday); got.Date() != 10 || got.Weekday() != time.Monday {
		t.Errorf("SetWeekday(Mon) = %v", got.ISO())
	}
	if got := base.SetWeekday(time.Sunday); got.Date() != 9 {
		t.Errorf("SetWeekday(Sun) = %v, want 2017-07-09", got.ISO())
	}
	// Time of day is preserved.
	if got := base.SetWeekday(time.Monday).Hour(); got != 12 {
		t.Errorf("SetWeekday hour = %d, want 12", got)
	}

	// SetISOWeekday: Monday (1) of the ISO week is 2017-07-10, Sunday (7) is
	// 2017-07-16.
	if got := base.SetISOWeekday(1); got.Date() != 10 {
		t.Errorf("SetISOWeekday(1) = %v, want 2017-07-10", got.ISO())
	}
	if got := base.SetISOWeekday(7); got.Date() != 16 {
		t.Errorf("SetISOWeekday(7) = %v, want 2017-07-16", got.ISO())
	}
}
