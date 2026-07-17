package moment

import "time"

// FromArray builds a Moment from a component slice in UTC, following moment.js's
// array form [year, month, day, hour, minute, second, millisecond]. As in
// moment, the month is 0-based (0 == January). Missing trailing components
// default to the start of their range (day and month default to 1). An empty
// slice yields an invalid Moment.
func FromArray(parts []int) Moment {
	return FromArrayInLocation(parts, time.UTC)
}

// FromArrayInLocation is FromArray in the given location. A nil location is
// treated as UTC.
func FromArrayInLocation(parts []int, loc *time.Location) Moment {
	if loc == nil {
		loc = time.UTC
	}
	if len(parts) == 0 {
		return Invalid()
	}
	get := func(i, def int) int {
		if i < len(parts) {
			return parts[i]
		}
		return def
	}
	year := get(0, 0)
	month := get(1, 0) // 0-based, moment style
	day := get(2, 1)
	hour := get(3, 0)
	min := get(4, 0)
	sec := get(5, 0)
	ms := get(6, 0)
	t := time.Date(year, time.Month(month+1), day, hour, min, sec, ms*int(time.Millisecond), loc)
	return New(t)
}

// FromObject builds a Moment from a map of component names to values in UTC,
// mirroring moment.js's object form. Recognized keys are "year"/"years"/"y",
// "month"/"months"/"M" (0-based), "day"/"days"/"date"/"D", "hour"/"hours"/"h",
// "minute"/"minutes"/"m", "second"/"seconds"/"s" and
// "millisecond"/"milliseconds"/"ms". Unset components default to the start of
// their range.
func FromObject(obj map[string]int) Moment {
	return FromObjectInLocation(obj, time.UTC)
}

// FromObjectInLocation is FromObject in the given location. A nil location is
// treated as UTC.
func FromObjectInLocation(obj map[string]int, loc *time.Location) Moment {
	if loc == nil {
		loc = time.UTC
	}
	pick := func(def int, keys ...string) int {
		for _, k := range keys {
			if v, ok := obj[k]; ok {
				return v
			}
		}
		return def
	}
	year := pick(0, "year", "years", "y")
	month := pick(0, "month", "months", "M") // 0-based, moment style
	day := pick(1, "day", "days", "date", "dates", "D", "d")
	hour := pick(0, "hour", "hours", "h")
	min := pick(0, "minute", "minutes", "m")
	sec := pick(0, "second", "seconds", "s")
	ms := pick(0, "millisecond", "milliseconds", "ms")
	t := time.Date(year, time.Month(month+1), day, hour, min, sec, ms*int(time.Millisecond), loc)
	return New(t)
}

// Max returns the latest of the given Moments. Invalid Moments are ignored;
// if no valid Moment is supplied the result is Invalid.
func Max(moments ...Moment) Moment {
	out := Invalid()
	for _, m := range moments {
		if !m.IsValid() {
			continue
		}
		if !out.IsValid() || m.IsAfter(out) {
			out = m
		}
	}
	return out
}

// Min returns the earliest of the given Moments. Invalid Moments are ignored;
// if no valid Moment is supplied the result is Invalid.
func Min(moments ...Moment) Moment {
	out := Invalid()
	for _, m := range moments {
		if !m.IsValid() {
			continue
		}
		if !out.IsValid() || m.IsBefore(out) {
			out = m
		}
	}
	return out
}
