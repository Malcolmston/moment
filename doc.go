// Package moment provides a moment.js-style date and time API layered on top of
// the Go standard library's time package. It is implemented purely with the
// standard library: no cgo and no third-party dependencies.
//
// # The Moment type
//
// A Moment is an immutable wrapper around time.Time. Every manipulation method
// returns a new Moment and never mutates the receiver, so values are safe to
// share. Construct one with New, FromTime, Now, Unix, UnixMilli, DateTime, or
// by parsing a string:
//
//	m := moment.Now()
//	m := moment.Unix(1500000000)
//	m, err := moment.Parse("2017-07-14T02:40:00Z")
//	m, err := moment.ParseFormat("14/07/2017", "DD/MM/YYYY")
//
// # Deterministic clock
//
// Relative-time helpers such as FromNow and Calendar consult a Clock. The
// default is backed by time.Now, but tests can inject a fixed clock so the
// output is deterministic:
//
//	clock := moment.FixedClock(reference)
//	m := moment.New(t).WithClock(clock)
//	fmt.Println(m.FromNow())
//
// # Format tokens
//
// Format and ParseFormat use moment.js-style tokens which are translated to Go
// reference layouts internally. The supported tokens are:
//
//	YYYY  four-digit year        2006
//	YY    two-digit year         06
//	MMMM  full month name        January
//	MMM   short month name       Jan
//	MM    two-digit month        01
//	M     month                  1
//	DD    two-digit day          02
//	D     day of month           2
//	dddd  full weekday name      Monday
//	ddd   short weekday name     Mon
//	HH    two-digit 24h hour     15
//	H     24h hour               (padded; Go has no unpadded form)
//	hh    two-digit 12h hour     03
//	h     12h hour               3
//	mm    two-digit minute       04
//	m     minute                 4
//	ss    two-digit second       05
//	s     second                 5
//	SSS   milliseconds           000
//	A     upper meridiem         PM
//	a     lower meridiem         pm
//	Z     numeric zone offset    -07:00
//	ZZ    numeric zone offset    -0700
//
// Text wrapped in square brackets is emitted literally, for example
// "YYYY [year]". FormatLayout and ParseLayout accept a raw Go layout instead.
//
// # Manipulation
//
// Add and Subtract move a Moment by a number of units; StartOf and EndOf snap
// to unit boundaries; Set replaces a single component; Clone copies the value.
// Units are the Unit constants (Year, Month, Week, Day, Hour, Minute, Second,
// Millisecond) and common moment.js aliases such as "days" or "h" are accepted.
//
//	m.Add(3, moment.Day)
//	m.Subtract(1, moment.Month)
//	m.StartOf(moment.Month)
//	m.EndOf(moment.Year)
//	m.Set(moment.Hour, 9)
//
// # Comparison and query
//
// IsBefore, IsAfter, IsSame, IsSameOrBefore, IsSameOrAfter and IsBetween
// compare Moments. Getters include Year, Month, Date, Hour, Minute, Second,
// Millisecond, Weekday, DayOfYear and ISOWeek.
//
// # Difference and relative time
//
// Diff reports the signed difference between two Moments in any unit as a
// float64; DiffInt truncates toward zero and DiffDuration returns a
// time.Duration. FromNow, From, To and ToNow produce humanized phrases such as
// "in 3 days" or "2 hours ago", Humanize renders a bare duration, and Calendar
// yields strings like "Today at 2:30 PM".
//
// # Time zones
//
// In, UTC and Local reinterpret a Moment in another time zone using
// time.Location values obtained from time.UTC, time.Local or
// time.LoadLocation.
package moment
