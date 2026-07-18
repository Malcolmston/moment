package moment

// ToJSON returns the Duration serialized for JSON as an ISO-8601 duration
// string, mirroring moment.js's Duration.toJSON (which aliases toISOString).
func (d Duration) ToJSON() string { return d.ISOString() }

// ValueOf returns the Duration's total length in whole milliseconds, mirroring
// moment.js's Duration.valueOf. Month and day components are converted using
// the same mean-month length As uses, then truncated toward zero.
func (d Duration) ValueOf() int64 {
	return int64(d.As(Millisecond))
}

// durObjectAliases maps the accepted moment.js object keys to canonical units,
// following moment.duration({...}).
var durObjectAliases = map[string]Unit{
	"year": Year, "years": Year, "y": Year,
	"quarter": Quarter, "quarters": Quarter, "Q": Quarter,
	"month": Month, "months": Month, "M": Month,
	"week": Week, "weeks": Week, "w": Week,
	"day": Day, "days": Day, "d": Day,
	"hour": Hour, "hours": Hour, "h": Hour,
	"minute": Minute, "minutes": Minute, "m": Minute,
	"second": Second, "seconds": Second, "s": Second,
	"millisecond": Millisecond, "milliseconds": Millisecond, "ms": Millisecond,
}

// NewDurationFromObject builds a Duration from a map of unit names to amounts,
// mirroring moment.js's moment.duration({...}) object form. Recognized keys are
// the moment.js unit names and their aliases, for example "years"/"y",
// "months"/"M", "weeks"/"w", "days"/"d", "hours"/"h", "minutes"/"m",
// "seconds"/"s" and "milliseconds"/"ms". Unknown keys are ignored; an empty map
// yields a zero Duration.
func NewDurationFromObject(obj map[string]int) Duration {
	var d Duration
	for key, n := range obj {
		if unit, ok := durObjectAliases[key]; ok {
			d = d.Add(NewDuration(n, unit))
		}
	}
	return d
}
