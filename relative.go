package moment

import (
	"fmt"
	"math"
	"time"
)

// round rounds a float to the nearest integer, halves away from zero.
func round(f float64) int64 {
	return int64(math.Round(f))
}

// relativeUnit describes the magnitude phrase for a duration, without any
// past/future suffix.
func relativeMagnitude(d time.Duration) string {
	abs := d
	if abs < 0 {
		abs = -abs
	}
	seconds := round(abs.Seconds())
	minutes := round(abs.Minutes())
	hours := round(abs.Hours())
	days := round(abs.Hours() / 24)
	months := round(abs.Hours() / 24 / 30)
	years := round(abs.Hours() / 24 / 365)

	switch {
	case seconds < 45:
		return "a few seconds"
	case seconds < 90:
		return "a minute"
	case minutes < 45:
		return fmt.Sprintf("%d minutes", minutes)
	case minutes < 90:
		return "an hour"
	case hours < 22:
		return fmt.Sprintf("%d hours", hours)
	case hours < 36:
		return "a day"
	case days < 26:
		return fmt.Sprintf("%d days", days)
	case days < 46:
		return "a month"
	case days < 320:
		return fmt.Sprintf("%d months", months)
	case days < 548:
		return "a year"
	default:
		return fmt.Sprintf("%d years", years)
	}
}

// Humanize renders a duration as a human-friendly phrase such as
// "2 hours" or "a few seconds", without any past/future suffix.
func Humanize(d time.Duration) string {
	return relativeMagnitude(d)
}

// From returns a humanized description of m relative to other, for example
// "in 3 days" or "2 hours ago".
func (m Moment) From(other Moment) string {
	return humanizeRelative(m.t.Sub(other.t))
}

// FromNow returns a humanized description of m relative to the current time as
// reported by the Moment's clock.
func (m Moment) FromNow() string {
	return humanizeRelative(m.t.Sub(m.clockOf().Now()))
}

// To returns a humanized description of other relative to m. It is the inverse
// of From.
func (m Moment) To(other Moment) string {
	return humanizeRelative(other.t.Sub(m.t))
}

// ToNow returns a humanized description of the current time relative to m.
func (m Moment) ToNow() string {
	return humanizeRelative(m.clockOf().Now().Sub(m.t))
}

// humanizeRelative adds the appropriate "in"/"ago" affix. A positive delta is
// in the future.
func humanizeRelative(delta time.Duration) string {
	phrase := relativeMagnitude(delta)
	if delta >= 0 {
		return "in " + phrase
	}
	return phrase + " ago"
}

// Calendar returns a moment.js-style calendar string describing m relative to
// reference, such as "Today at 2:30 PM" or "Last Friday at 9:00 AM". Days
// outside a one-week window fall back to a numeric date.
func (m Moment) Calendar(reference Moment) string {
	ref := reference.In(m.Location())
	startM := m.StartOf(Day)
	startRef := ref.StartOf(Day)
	dayDelta := int(math.Round(startM.t.Sub(startRef.t).Hours() / 24))
	clock := m.Format("h:mm A")
	switch {
	case dayDelta == 0:
		return "Today at " + clock
	case dayDelta == 1:
		return "Tomorrow at " + clock
	case dayDelta == -1:
		return "Yesterday at " + clock
	case dayDelta > 1 && dayDelta < 7:
		return m.Format("dddd") + " at " + clock
	case dayDelta < -1 && dayDelta > -7:
		return "Last " + m.Format("dddd") + " at " + clock
	default:
		return m.Format("MM/DD/YYYY")
	}
}

// CalendarNow is Calendar relative to the Moment's clock.
func (m Moment) CalendarNow() string {
	return m.Calendar(NowWith(m.clockOf()))
}
