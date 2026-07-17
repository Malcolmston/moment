package moment

import "time"

// Year returns the year.
func (m Moment) Year() int { return m.t.Year() }

// Month returns the month.
func (m Moment) Month() time.Month { return m.t.Month() }

// Date returns the day of the month (1-31).
func (m Moment) Date() int { return m.t.Day() }

// Day is an alias for Date, returning the day of the month.
func (m Moment) Day() int { return m.t.Day() }

// Hour returns the hour within the day (0-23).
func (m Moment) Hour() int { return m.t.Hour() }

// Minute returns the minute offset within the hour (0-59).
func (m Moment) Minute() int { return m.t.Minute() }

// Second returns the second offset within the minute (0-59).
func (m Moment) Second() int { return m.t.Second() }

// Millisecond returns the millisecond offset within the second (0-999).
func (m Moment) Millisecond() int { return m.t.Nanosecond() / int(time.Millisecond) }

// Nanosecond returns the nanosecond offset within the second.
func (m Moment) Nanosecond() int { return m.t.Nanosecond() }

// Weekday returns the day of the week.
func (m Moment) Weekday() time.Weekday { return m.t.Weekday() }

// DayOfYear returns the 1-based day of the year (1-366).
func (m Moment) DayOfYear() int { return m.t.YearDay() }

// ISOWeek returns the ISO 8601 year and week number.
func (m Moment) ISOWeek() (year, week int) { return m.t.ISOWeek() }

// IsBefore reports whether m is strictly before other.
func (m Moment) IsBefore(other Moment) bool { return m.t.Before(other.t) }

// IsAfter reports whether m is strictly after other.
func (m Moment) IsAfter(other Moment) bool { return m.t.After(other.t) }

// IsSame reports whether m and other represent the same instant.
func (m Moment) IsSame(other Moment) bool { return m.t.Equal(other.t) }

// IsSameOrBefore reports whether m is before or equal to other.
func (m Moment) IsSameOrBefore(other Moment) bool { return m.IsBefore(other) || m.IsSame(other) }

// IsSameOrAfter reports whether m is after or equal to other.
func (m Moment) IsSameOrAfter(other Moment) bool { return m.IsAfter(other) || m.IsSame(other) }

// IsBetween reports whether m lies within (start, end). When inclusive is true
// the bounds themselves count as inside the range. The bounds may be supplied
// in either order.
func (m Moment) IsBetween(start, end Moment, inclusive bool) bool {
	lo, hi := start, end
	if hi.IsBefore(lo) {
		lo, hi = hi, lo
	}
	if inclusive {
		return m.IsSameOrAfter(lo) && m.IsSameOrBefore(hi)
	}
	return m.IsAfter(lo) && m.IsBefore(hi)
}

// IsSameUnit reports whether m and other fall in the same calendar unit (for
// example the same day or the same month), comparing after truncation.
func (m Moment) IsSameUnit(other Moment, unit Unit) bool {
	return m.StartOf(unit).IsSame(other.In(m.Location()).StartOf(unit))
}

// IsLeapYear reports whether the Moment falls in a leap year.
func (m Moment) IsLeapYear() bool {
	y := m.t.Year()
	return y%4 == 0 && (y%100 != 0 || y%400 == 0)
}
