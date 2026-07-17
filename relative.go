package moment

import (
	"math"
	"time"
)

// Humanize renders a bare time.Duration as a human-friendly phrase such as
// "2 hours" or "a few seconds", without any past/future suffix, using the
// global default locale.
func Humanize(d time.Duration) string {
	return DurationFromTime(d).Humanize(false)
}

// From returns a humanized, locale-aware description of m relative to other,
// for example "in 3 days" or "2 hours ago".
func (m Moment) From(other Moment) string {
	return DurationBetween(other, m).localeWith(m.localeOf()).Humanize(true)
}

// FromNow returns a humanized description of m relative to the current time as
// reported by the Moment's clock.
func (m Moment) FromNow() string {
	now := Moment{t: m.clockOf().Now()}
	return DurationBetween(now, m).localeWith(m.localeOf()).Humanize(true)
}

// To returns a humanized description of other relative to m. It is the inverse
// of From.
func (m Moment) To(other Moment) string {
	return DurationBetween(m, other).localeWith(m.localeOf()).Humanize(true)
}

// ToNow returns a humanized description of the current time relative to m.
func (m Moment) ToNow() string {
	now := Moment{t: m.clockOf().Now()}
	return DurationBetween(m, now).localeWith(m.localeOf()).Humanize(true)
}

// localeWith binds a locale to a Duration for the relative-time helpers.
func (d Duration) localeWith(l *Locale) Duration {
	d.loc = l
	return d
}

// Calendar returns a locale-aware calendar string describing m relative to
// reference, such as "Today at 2:30 PM" or "Last Friday at 9:00 AM". Days
// outside a one-week window fall back to the locale's numeric date (sameElse).
func (m Moment) Calendar(reference Moment) string {
	loc := m.localeOf()
	ref := reference.In(m.Location())
	startM := m.StartOf(Day)
	startRef := ref.StartOf(Day)
	dayDelta := int(math.Round(startM.t.Sub(startRef.t).Hours() / 24))
	var pattern string
	switch {
	case dayDelta == 0:
		pattern = loc.Calendar.SameDay
	case dayDelta == 1:
		pattern = loc.Calendar.NextDay
	case dayDelta == -1:
		pattern = loc.Calendar.LastDay
	case dayDelta > 1 && dayDelta < 7:
		pattern = loc.Calendar.NextWeek
	case dayDelta < -1 && dayDelta > -7:
		pattern = loc.Calendar.LastWeek
	default:
		pattern = loc.Calendar.SameElse
	}
	return m.formatWith(pattern, loc)
}

// CalendarNow is Calendar relative to the Moment's clock.
func (m Moment) CalendarNow() string {
	return m.Calendar(Moment{t: m.clockOf().Now(), clock: m.clockOf()})
}
