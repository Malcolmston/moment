package moment

import "time"

// SetYear returns a new Moment with the year replaced by y, the setter
// companion to Year and the analogue of moment.js's year(value).
func (m Moment) SetYear(y int) Moment { return m.Set(Year, y) }

// SetMonth returns a new Moment with the month replaced by mo, the setter
// companion to Month and the analogue of moment.js's month(value). Overflow
// beyond a month's length rolls forward as in time.Date.
func (m Moment) SetMonth(mo time.Month) Moment { return m.Set(Month, int(mo)) }

// SetQuarter returns a new Moment moved to quarter q, matching moment.js's
// quarter(value). Following moment, the month keeps its position within the
// quarter (its month-mod-3 offset) rather than snapping to the first month, and
// values outside 1–4 bubble the year forward or backward.
func (m Moment) SetQuarter(q int) Moment { return m.Set(Quarter, q) }

// SetDate returns a new Moment with the day of the month replaced by d, the
// setter companion to Date and the analogue of moment.js's date(value).
func (m Moment) SetDate(d int) Moment { return m.Set(Date, d) }

// SetDayOfYear returns a new Moment set to the d-th day of its year (1-based),
// matching moment.js's dayOfYear(value).
func (m Moment) SetDayOfYear(d int) Moment { return m.Set(DayOfYear, d) }

// SetHour returns a new Moment with the hour replaced by h, matching
// moment.js's hour(value).
func (m Moment) SetHour(h int) Moment { return m.Set(Hour, h) }

// SetMinute returns a new Moment with the minute replaced by mi, matching
// moment.js's minute(value).
func (m Moment) SetMinute(mi int) Moment { return m.Set(Minute, mi) }

// SetSecond returns a new Moment with the second replaced by s, matching
// moment.js's second(value).
func (m Moment) SetSecond(s int) Moment { return m.Set(Second, s) }

// SetMillisecond returns a new Moment with the millisecond replaced by ms,
// matching moment.js's millisecond(value).
func (m Moment) SetMillisecond(ms int) Moment { return m.Set(Millisecond, ms) }

// SetWeekday returns a new Moment moved to weekday w within the same
// Sunday-based week, matching moment.js's day(value): the result stays in the
// current week and only the day of the week changes.
func (m Moment) SetWeekday(w time.Weekday) Moment {
	return m.Add(int(w)-int(m.t.Weekday()), Day)
}

// SetISOWeekday returns a new Moment moved to ISO weekday d (1 == Monday, 7 ==
// Sunday) within the same ISO week, matching moment.js's isoWeekday(value).
func (m Moment) SetISOWeekday(d int) Moment {
	return m.Add(d-m.ISOWeekday(), Day)
}
