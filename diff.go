package moment

import (
	"math"
	"time"
)

// monthDiff computes a - b measured in fractional months, following moment.js's
// algorithm so that partial months interpolate smoothly. The month anchors are
// advanced with addMonths, which clamps the day of the month, so end-of-month
// boundaries agree with moment (for example February 29 minus January 30 is
// exactly one month).
func monthDiff(a, b time.Time) float64 {
	// moment normalizes so the earlier-in-month operand is subtracted, then
	// negates, which keeps the fractional adjustment well behaved.
	if a.Day() < b.Day() {
		return -monthDiff(b, a)
	}
	wholeMonths := (b.Year()-a.Year())*12 + (int(b.Month()) - int(a.Month()))
	anchor := addMonths(a, wholeMonths)
	var adjust float64
	if b.Sub(anchor) < 0 {
		anchor2 := addMonths(a, wholeMonths-1)
		adjust = float64(b.Sub(anchor)) / float64(anchor.Sub(anchor2))
	} else {
		anchor2 := addMonths(a, wholeMonths+1)
		adjust = float64(b.Sub(anchor)) / float64(anchor2.Sub(anchor))
	}
	return -(float64(wholeMonths) + adjust)
}

// Diff returns the signed difference m - other expressed in the given unit as a
// floating-point value. Years and months use civil arithmetic; smaller units
// use exact durations.
func (m Moment) Diff(other Moment, unit Unit) float64 {
	switch normalizeUnit(unit) {
	case Year:
		return monthDiff(m.t, other.t) / 12
	case Quarter:
		return monthDiff(m.t, other.t) / 3
	case Month:
		return monthDiff(m.t, other.t)
	case Week, ISOWeek:
		return m.t.Sub(other.t).Hours() / (24 * 7)
	case Day, Date, DayOfYear:
		return m.t.Sub(other.t).Hours() / 24
	case Hour:
		return m.t.Sub(other.t).Hours()
	case Minute:
		return m.t.Sub(other.t).Minutes()
	case Second:
		return m.t.Sub(other.t).Seconds()
	case Millisecond:
		return float64(m.t.Sub(other.t).Milliseconds())
	default:
		return 0
	}
}

// DiffInt returns Diff truncated toward zero, matching moment.js's default
// integer diff behaviour.
func (m Moment) DiffInt(other Moment, unit Unit) int {
	return int(math.Trunc(m.Diff(other, unit)))
}

// DiffDuration returns the signed difference m - other as a time.Duration.
func (m Moment) DiffDuration(other Moment) time.Duration {
	return m.t.Sub(other.t)
}
