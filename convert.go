package moment

import "time"

// isoZoneLayout renders an instant with millisecond precision while keeping the
// zone offset, matching moment.js toISOString(true).
const isoZoneLayout = "2006-01-02T15:04:05.000Z07:00"

// ToArray returns the Moment's calendar and clock components as a slice in the
// moment.js toArray order: [year, month, day, hour, minute, second,
// millisecond]. As in moment.js the month is 0-based (0 == January), matching
// FromArray. An invalid Moment yields nil.
func (m Moment) ToArray() []int {
	if !m.IsValid() {
		return nil
	}
	return []int{
		m.t.Year(),
		int(m.t.Month()) - 1,
		m.t.Day(),
		m.t.Hour(),
		m.t.Minute(),
		m.t.Second(),
		m.t.Nanosecond() / int(time.Millisecond),
	}
}

// ToObject returns the Moment's components keyed by moment.js's toObject names:
// "years", "months" (0-based), "date", "hours", "minutes", "seconds" and
// "milliseconds". The keys mirror FromObject. An invalid Moment yields nil.
func (m Moment) ToObject() map[string]int {
	if !m.IsValid() {
		return nil
	}
	return map[string]int{
		"years":        m.t.Year(),
		"months":       int(m.t.Month()) - 1,
		"date":         m.t.Day(),
		"hours":        m.t.Hour(),
		"minutes":      m.t.Minute(),
		"seconds":      m.t.Second(),
		"milliseconds": m.t.Nanosecond() / int(time.Millisecond),
	}
}

// ToDate returns the underlying time.Time value, mirroring moment.js's
// toDate(). It is an alias for Time.
func (m Moment) ToDate() time.Time { return m.t }

// ToJSON returns the Moment serialized for JSON as a UTC ISO-8601 string with
// millisecond precision, matching moment.js's toJSON (which aliases
// toISOString). An invalid Moment yields the empty string, mirroring moment's
// null.
func (m Moment) ToJSON() string {
	if !m.IsValid() {
		return ""
	}
	return m.ToISOString()
}

// ToISOStringZone returns the moment.js toISOString(true) representation: the
// instant with millisecond precision but keeping the Moment's own zone offset
// rather than converting to UTC, e.g. "2017-07-14T02:40:00.000-05:00". An
// invalid Moment yields the empty string.
func (m Moment) ToISOStringZone() string {
	if !m.IsValid() {
		return ""
	}
	return m.t.Format(isoZoneLayout)
}

// Get returns a single component of the Moment identified by unit, the getter
// companion to Set and the analogue of moment.js's get. Recognized units are
// Year, Quarter, Month, Week, ISOWeek, Day/Date, DayOfYear, Hour, Minute,
// Second and Millisecond (plus their aliases). Unlike moment's 0-based month,
// Get(Month) returns the 1-based month to stay consistent with Month and
// Set(Month, …). Unknown units yield 0.
func (m Moment) Get(unit Unit) int {
	switch normalizeUnit(unit) {
	case Year:
		return m.t.Year()
	case Quarter:
		return m.Quarter()
	case Month:
		return int(m.t.Month())
	case Week:
		return m.Week()
	case ISOWeek:
		return m.ISOWeekNumber()
	case Day, Date:
		return m.t.Day()
	case DayOfYear:
		return m.t.YearDay()
	case Hour:
		return m.t.Hour()
	case Minute:
		return m.t.Minute()
	case Second:
		return m.t.Second()
	case Millisecond:
		return m.t.Nanosecond() / int(time.Millisecond)
	default:
		return 0
	}
}

// DaysInYear returns the number of days in the Moment's year: 366 in a leap
// year and 365 otherwise.
func (m Moment) DaysInYear() int {
	return daysInYearOf(m.t.Year())
}

// IsUTC reports whether the Moment's location is UTC, mirroring moment.js's
// isUTC.
func (m Moment) IsUTC() bool {
	return m.t.Location() == time.UTC
}

// IsLocal reports whether the Moment's location is the system local zone,
// mirroring moment.js's isLocal.
func (m Moment) IsLocal() bool {
	return m.t.Location() == time.Local
}

// ZoneAbbr returns the abbreviated zone name in effect at the Moment, such as
// "UTC" or "EST", matching moment.js's zoneAbbr. The abbreviation comes from
// the Moment's time.Location and may be a numeric offset for fixed zones.
func (m Moment) ZoneAbbr() string {
	name, _ := m.t.Zone()
	return name
}

// ZoneName returns the Moment's full time-zone identifier, for example "UTC" or
// "America/New_York". It is the loaded location name; use ZoneAbbr for the
// short form that moment.js prints.
func (m Moment) ZoneName() string {
	return m.t.Location().String()
}

// LocaleWeekday returns the day of the week numbered from the Moment locale's
// first day of the week (0 == the locale's first weekday), matching moment.js's
// locale-aware weekday getter. Compare ISOWeekday, which is always Monday-based.
func (m Moment) LocaleWeekday() int {
	return m.localeWeekday(m.localeOf())
}

// IsMoment reports whether v is a Moment value, mirroring moment.js's
// moment.isMoment.
func IsMoment(v any) bool {
	_, ok := v.(Moment)
	return ok
}

// IsDuration reports whether v is a Duration value, mirroring moment.js's
// moment.isDuration.
func IsDuration(v any) bool {
	_, ok := v.(Duration)
	return ok
}

// IsDate reports whether v is a time.Time value, the Go analogue of moment.js's
// moment.isDate (which tests for a JavaScript Date).
func IsDate(v any) bool {
	_, ok := v.(time.Time)
	return ok
}
