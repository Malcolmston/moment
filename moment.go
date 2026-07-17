package moment

import (
	"time"
)

// Clock is the source of the current time. It is injected into Moment values so
// that relative-time helpers such as FromNow and Calendar are deterministic in
// tests. The standard implementation is backed by time.Now.
type Clock interface {
	// Now returns the current time.
	Now() time.Time
}

// ClockFunc adapts an ordinary function to the Clock interface.
type ClockFunc func() time.Time

// Now calls the underlying function.
func (f ClockFunc) Now() time.Time { return f() }

// systemClock is the default, real-world clock.
var systemClock Clock = ClockFunc(time.Now)

// FixedClock returns a Clock that always reports the same instant. It is handy
// for deterministic tests.
func FixedClock(t time.Time) Clock {
	return ClockFunc(func() time.Time { return t })
}

// Moment is an immutable wrapper around time.Time offering a moment.js-style
// API for parsing, formatting, manipulation and comparison. All manipulation
// methods return a new Moment and never mutate the receiver.
type Moment struct {
	t        time.Time
	clock    Clock
	loc      *Locale
	invalid  bool
	creation *CreationData
}

// CreationData records how a Moment was constructed by a parsing function,
// mirroring moment.js's creationData(). It is returned by Moment.CreationData.
type CreationData struct {
	// Input is the original string (or a description of the input) supplied to
	// the parser.
	Input string
	// Format is the format string that produced the Moment, if any.
	Format string
	// Locale is the locale name in effect during parsing.
	Locale string
	// IsUTC reports whether the Moment was created in UTC.
	IsUTC bool
	// Valid reports whether parsing succeeded.
	Valid bool
}

// clockOf returns the effective clock for a Moment, falling back to the system
// clock for zero-value receivers.
func (m Moment) clockOf() Clock {
	if m.clock == nil {
		return systemClock
	}
	return m.clock
}

// New wraps an existing time.Time in a Moment using the system clock.
func New(t time.Time) Moment {
	return Moment{t: t, clock: systemClock}
}

// FromTime is an alias for New, wrapping a time.Time value.
func FromTime(t time.Time) Moment {
	return New(t)
}

// Now returns a Moment for the current instant according to the system clock.
func Now() Moment {
	return Moment{t: systemClock.Now(), clock: systemClock}
}

// NowWith returns a Moment for the current instant according to the supplied
// clock. The returned Moment carries that clock so relative-time helpers stay
// deterministic.
func NowWith(clock Clock) Moment {
	if clock == nil {
		clock = systemClock
	}
	return Moment{t: clock.Now(), clock: clock}
}

// Unix returns a Moment for the given Unix timestamp in seconds, in UTC.
func Unix(sec int64) Moment {
	return New(time.Unix(sec, 0).UTC())
}

// UnixMilli returns a Moment for the given Unix timestamp in milliseconds, in
// UTC.
func UnixMilli(ms int64) Moment {
	return New(time.UnixMilli(ms).UTC())
}

// DateTime builds a Moment from calendar components in the given location. A
// nil location is treated as UTC.
func DateTime(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) Moment {
	if loc == nil {
		loc = time.UTC
	}
	return New(time.Date(year, month, day, hour, min, sec, nsec, loc))
}

// WithClock returns a copy of the Moment that uses the supplied clock for
// relative-time helpers. A nil clock resets to the system clock.
func (m Moment) WithClock(clock Clock) Moment {
	if clock == nil {
		clock = systemClock
	}
	m.clock = clock
	return m
}

// Clone returns an independent copy of the Moment.
func (m Moment) Clone() Moment {
	return m
}

// Time returns the underlying time.Time value.
func (m Moment) Time() time.Time {
	return m.t
}

// IsZero reports whether the underlying instant is the zero time.
func (m Moment) IsZero() bool {
	return m.t.IsZero()
}

// IsValid reports whether the Moment holds a usable instant. It is false for
// the zero value and for Moments returned by failed parses (see Invalid),
// mirroring moment.js's isValid.
func (m Moment) IsValid() bool {
	return !m.invalid && !m.t.IsZero()
}

// Invalid returns an explicitly invalid Moment, the moral equivalent of
// moment.invalid(). Its IsValid reports false and Format returns "Invalid
// date".
func Invalid() Moment {
	return Moment{invalid: true, clock: systemClock, creation: &CreationData{Valid: false}}
}

// CreationData returns how the Moment was produced by a parsing constructor, or
// nil when it was not created by one.
func (m Moment) CreationData() *CreationData {
	return m.creation
}

// Unix returns the Moment as a Unix timestamp in seconds.
func (m Moment) Unix() int64 { return m.t.Unix() }

// UnixMilli returns the Moment as a Unix timestamp in milliseconds.
func (m Moment) UnixMilli() int64 { return m.t.UnixMilli() }

// ValueOf returns the Moment as a Unix timestamp in milliseconds, matching
// moment.js's valueOf.
func (m Moment) ValueOf() int64 { return m.t.UnixMilli() }

// Location returns the time zone associated with the Moment.
func (m Moment) Location() *time.Location { return m.t.Location() }

// In returns a copy of the Moment representing the same instant in the given
// location. A nil location is treated as UTC.
func (m Moment) In(loc *time.Location) Moment {
	if loc == nil {
		loc = time.UTC
	}
	m.t = m.t.In(loc)
	return m
}

// UTC returns a copy of the Moment in the UTC location.
func (m Moment) UTC() Moment { return m.In(time.UTC) }

// Local returns a copy of the Moment in the local location.
func (m Moment) Local() Moment { return m.In(time.Local) }

// String implements fmt.Stringer, returning the ISO-8601 representation.
func (m Moment) String() string { return m.ISO() }
