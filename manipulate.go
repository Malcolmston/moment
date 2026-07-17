package moment

import (
	"time"
)

// Unit identifies a calendar or clock unit used by Add, Subtract, StartOf,
// EndOf, Set and Diff.
type Unit string

// Supported units. Manipulation and diff helpers also accept common moment.js
// aliases (plural and short forms) via normalizeUnit.
const (
	Year        Unit = "year"
	Month       Unit = "month"
	Week        Unit = "week"
	Day         Unit = "day"
	Date        Unit = "date"
	Hour        Unit = "hour"
	Minute      Unit = "minute"
	Second      Unit = "second"
	Millisecond Unit = "millisecond"
)

// normalizeUnit maps moment.js aliases to the canonical Unit values.
func normalizeUnit(u Unit) Unit {
	switch u {
	case "year", "years", "y":
		return Year
	case "month", "months", "M":
		return Month
	case "week", "weeks", "w":
		return Week
	case "day", "days", "d":
		return Day
	case "date", "dates", "D":
		return Date
	case "hour", "hours", "h":
		return Hour
	case "minute", "minutes", "m":
		return Minute
	case "second", "seconds", "s":
		return Second
	case "millisecond", "milliseconds", "ms":
		return Millisecond
	default:
		return u
	}
}

// Add returns a new Moment advanced by n of the given unit. Calendar units
// (years, months, weeks, days) use civil arithmetic via AddDate; clock units
// use fixed durations.
func (m Moment) Add(n int, unit Unit) Moment {
	switch normalizeUnit(unit) {
	case Year:
		m.t = m.t.AddDate(n, 0, 0)
	case Month:
		m.t = m.t.AddDate(0, n, 0)
	case Week:
		m.t = m.t.AddDate(0, 0, n*7)
	case Day, Date:
		m.t = m.t.AddDate(0, 0, n)
	case Hour:
		m.t = m.t.Add(time.Duration(n) * time.Hour)
	case Minute:
		m.t = m.t.Add(time.Duration(n) * time.Minute)
	case Second:
		m.t = m.t.Add(time.Duration(n) * time.Second)
	case Millisecond:
		m.t = m.t.Add(time.Duration(n) * time.Millisecond)
	}
	return m
}

// Subtract returns a new Moment moved back by n of the given unit.
func (m Moment) Subtract(n int, unit Unit) Moment {
	return m.Add(-n, unit)
}

// AddDuration returns a new Moment advanced by an arbitrary time.Duration.
func (m Moment) AddDuration(d time.Duration) Moment {
	m.t = m.t.Add(d)
	return m
}

// StartOf returns a new Moment truncated to the start of the given unit,
// preserving the time zone. Week starts on Sunday, matching moment.js's
// default locale.
func (m Moment) StartOf(unit Unit) Moment {
	y, mo, d := m.t.Date()
	h, mi, s := m.t.Clock()
	loc := m.t.Location()
	switch normalizeUnit(unit) {
	case Year:
		m.t = time.Date(y, time.January, 1, 0, 0, 0, 0, loc)
	case Month:
		m.t = time.Date(y, mo, 1, 0, 0, 0, 0, loc)
	case Week:
		start := time.Date(y, mo, d, 0, 0, 0, 0, loc)
		m.t = start.AddDate(0, 0, -int(start.Weekday()))
	case Day, Date:
		m.t = time.Date(y, mo, d, 0, 0, 0, 0, loc)
	case Hour:
		m.t = time.Date(y, mo, d, h, 0, 0, 0, loc)
	case Minute:
		m.t = time.Date(y, mo, d, h, mi, 0, 0, loc)
	case Second:
		m.t = time.Date(y, mo, d, h, mi, s, 0, loc)
	case Millisecond:
		m.t = m.t.Truncate(time.Millisecond)
	}
	return m
}

// EndOf returns a new Moment set to the last nanosecond of the given unit.
func (m Moment) EndOf(unit Unit) Moment {
	u := normalizeUnit(unit)
	next := m.StartOf(u)
	switch u {
	case Year:
		next.t = next.t.AddDate(1, 0, 0)
	case Month:
		next.t = next.t.AddDate(0, 1, 0)
	case Week:
		next.t = next.t.AddDate(0, 0, 7)
	case Day, Date:
		next.t = next.t.AddDate(0, 0, 1)
	case Hour:
		next.t = next.t.Add(time.Hour)
	case Minute:
		next.t = next.t.Add(time.Minute)
	case Second:
		next.t = next.t.Add(time.Second)
	case Millisecond:
		next.t = next.t.Add(time.Millisecond)
	default:
		return m
	}
	next.t = next.t.Add(-time.Nanosecond)
	return next
}

// Set returns a new Moment with a single component replaced by value.
func (m Moment) Set(unit Unit, value int) Moment {
	y, mo, d := m.t.Date()
	h, mi, s := m.t.Clock()
	ns := m.t.Nanosecond()
	loc := m.t.Location()
	switch normalizeUnit(unit) {
	case Year:
		y = value
	case Month:
		mo = time.Month(value)
	case Day, Date:
		d = value
	case Hour:
		h = value
	case Minute:
		mi = value
	case Second:
		s = value
	case Millisecond:
		ns = value * int(time.Millisecond)
	}
	m.t = time.Date(y, mo, d, h, mi, s, ns, loc)
	return m
}
