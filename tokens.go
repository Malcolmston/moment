package moment

import (
	"strings"
	"time"
)

// tokenPair maps a moment.js-style token to its equivalent fragment in the Go
// reference layout (Mon Jan 2 15:04:05 MST 2006).
type tokenPair struct {
	moment string
	golang string
}

// tokenTable lists the supported moment tokens ordered longest-first so that
// the scanner performs greedy, unambiguous matching (e.g. YYYY before YY).
var tokenTable = []tokenPair{
	{"YYYY", "2006"},
	{"MMMM", "January"},
	{"dddd", "Monday"},
	{"MMM", "Jan"},
	{"ddd", "Mon"},
	{"SSS", "000"},
	{"YY", "06"},
	{"MM", "01"},
	{"DD", "02"},
	{"HH", "15"},
	{"hh", "03"},
	{"mm", "04"},
	{"ss", "05"},
	{"ZZ", "-0700"},
	{"M", "1"},
	{"D", "2"},
	{"H", "15"}, // Go has no unpadded 24-hour token; use padded as the closest match.
	{"h", "3"},
	{"m", "4"},
	{"s", "5"},
	{"A", "PM"},
	{"a", "pm"},
	{"Z", "-07:00"},
}

// SupportedTokens returns the moment-style tokens understood by the format
// translator, longest-first.
func SupportedTokens() []string {
	out := make([]string, 0, len(tokenTable))
	for _, p := range tokenTable {
		out = append(out, p.moment)
	}
	return out
}

// tokenToLayout converts a moment-style format string into a Go reference
// layout. Text wrapped in square brackets is treated as a literal and copied
// through verbatim, matching moment.js escaping.
func tokenToLayout(format string) string {
	var b strings.Builder
	for i := 0; i < len(format); {
		if format[i] == '[' {
			if end := strings.IndexByte(format[i+1:], ']'); end >= 0 {
				b.WriteString(format[i+1 : i+1+end])
				i += end + 2
				continue
			}
		}
		matched := false
		for _, p := range tokenTable {
			if strings.HasPrefix(format[i:], p.moment) {
				b.WriteString(p.golang)
				i += len(p.moment)
				matched = true
				break
			}
		}
		if !matched {
			b.WriteByte(format[i])
			i++
		}
	}
	return b.String()
}

// Format renders the Moment using a moment.js-style token string, for example
// "YYYY-MM-DD HH:mm:ss" or "dddd, MMMM D, YYYY".
func (m Moment) Format(format string) string {
	return m.t.Format(tokenToLayout(format))
}

// FormatLayout renders the Moment using a raw Go reference layout.
func (m Moment) FormatLayout(layout string) string {
	return m.t.Format(layout)
}

// ISO returns the RFC3339 (ISO-8601) representation of the Moment.
func (m Moment) ISO() string {
	return m.t.Format(time.RFC3339Nano)
}

// isoLayouts is the set of standard layouts attempted by Parse, in priority
// order.
var isoLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
	"2006-01-02",
	time.RFC1123Z,
	time.RFC1123,
	time.RFC822Z,
	time.RFC822,
	time.RFC850,
	time.ANSIC,
	time.UnixDate,
	"2006/01/02",
	"01/02/2006",
	"15:04:05",
	"15:04",
}

// Parse attempts to interpret value using a list of common ISO-8601 and RFC
// layouts. It returns an error if none of them match.
func Parse(value string) (Moment, error) {
	var firstErr error
	for _, layout := range isoLayouts {
		if t, err := time.Parse(layout, value); err == nil {
			return New(t), nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	return Moment{}, firstErr
}

// ParseFormat parses value using a moment.js-style token format string.
func ParseFormat(value, format string) (Moment, error) {
	t, err := time.Parse(tokenToLayout(format), value)
	if err != nil {
		return Moment{}, err
	}
	return New(t), nil
}

// ParseInLocation parses value using a moment-style token format, interpreting
// zone-less input as being in loc. A nil location is treated as UTC.
func ParseInLocation(value, format string, loc *time.Location) (Moment, error) {
	if loc == nil {
		loc = time.UTC
	}
	t, err := time.ParseInLocation(tokenToLayout(format), value, loc)
	if err != nil {
		return Moment{}, err
	}
	return New(t), nil
}

// ParseLayout parses value using a raw Go reference layout.
func ParseLayout(value, layout string) (Moment, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Moment{}, err
	}
	return New(t), nil
}
