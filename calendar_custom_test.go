package moment

import (
	"testing"
	"time"
)

func TestCalendarWith(t *testing.T) {
	ref := DateTime(2017, 7, 14, 12, 0, 0, 0, time.UTC)

	// Same day with a custom template.
	sameDay := ref.Add(3, Hour)
	if got := sameDay.CalendarWith(ref, CalendarFormats{SameDay: "[Right now]"}); got != "Right now" {
		t.Errorf("SameDay override = %q, want %q", got, "Right now")
	}

	// Empty fields fall back to the locale default.
	nextDay := ref.Add(1, Day)
	if got, want := nextDay.CalendarWith(ref, CalendarFormats{}), nextDay.Calendar(ref); got != want {
		t.Errorf("empty override = %q, want default %q", got, want)
	}

	// Custom SameElse for far dates.
	far := ref.Add(40, Day)
	if got := far.CalendarWith(ref, CalendarFormats{SameElse: "YYYY-MM-DD"}); got != "2017-08-23" {
		t.Errorf("SameElse override = %q, want 2017-08-23", got)
	}
}
