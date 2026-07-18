package moment

import "math"

// CalendarWith is Calendar with caller-supplied templates, mirroring moment.js's
// calendar(referenceDay, formats). Each non-empty field of formats overrides the
// corresponding locale template for that bucket (same day, next day, and so on);
// empty fields fall back to the Moment locale's default calendar templates. This
// lets callers customise the phrasing without registering a whole locale.
func (m Moment) CalendarWith(reference Moment, formats CalendarFormats) string {
	loc := m.localeOf()
	def := loc.Calendar
	pick := func(override, fallback string) string {
		if override != "" {
			return override
		}
		return fallback
	}
	ref := reference.In(m.Location())
	dayDelta := int(math.Round(m.StartOf(Day).t.Sub(ref.StartOf(Day).t).Hours() / 24))
	var pattern string
	switch {
	case dayDelta == 0:
		pattern = pick(formats.SameDay, def.SameDay)
	case dayDelta == 1:
		pattern = pick(formats.NextDay, def.NextDay)
	case dayDelta == -1:
		pattern = pick(formats.LastDay, def.LastDay)
	case dayDelta > 1 && dayDelta < 7:
		pattern = pick(formats.NextWeek, def.NextWeek)
	case dayDelta < -1 && dayDelta > -7:
		pattern = pick(formats.LastWeek, def.LastWeek)
	default:
		pattern = pick(formats.SameElse, def.SameElse)
	}
	return m.formatWith(pattern, loc)
}
