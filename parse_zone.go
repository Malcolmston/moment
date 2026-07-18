package moment

import "time"

// ParseZone parses an ISO-8601 or RFC date/time string and keeps the UTC offset
// carried by the input, mirroring moment.js's moment.parseZone. Whereas some
// constructors normalise to UTC or local time, the Moment returned here is
// pinned to a fixed zone equal to the parsed offset, so Format, ISO and
// UTCOffset report that offset. Inputs without an explicit offset are treated as
// UTC. It returns an invalid Moment and an error when no layout matches.
func ParseZone(value string) (Moment, error) {
	m, err := Parse(value)
	if err != nil {
		return m, err
	}
	name, off := m.t.Zone()
	m.t = m.t.In(time.FixedZone(name, off))
	if m.creation != nil {
		m.creation.Format = "parseZone"
	}
	return m, nil
}
