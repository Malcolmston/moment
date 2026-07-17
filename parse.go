package moment

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// isoLayouts is the set of standard layouts attempted by Parse and ParseISO, in
// priority order.
var isoLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04",
	"2006-01-02 15:04",
	"2006-01-02",
	"2006-W02-1",
	"20060102T150405Z0700",
	"20060102",
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
			m := New(t)
			m.creation = &CreationData{Input: value, Locale: GlobalLocale(), IsUTC: t.Location() == time.UTC, Valid: true}
			return m, nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	return Invalid(), firstErr
}

// ParseISO parses an ISO-8601 date/time string. It is a synonym for Parse that
// documents intent when the input is known to be ISO-8601.
func ParseISO(value string) (Moment, error) {
	return Parse(value)
}

// rfc2822Layouts lists the layouts ParseRFC2822 attempts.
var rfc2822Layouts = []string{
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Mon, 2 Jan 2006 15:04:05 -0700",
	"02 Jan 2006 15:04:05 -0700",
	"2 Jan 2006 15:04:05 -0700",
	"Mon, 02 Jan 2006 15:04:05 MST",
	"Mon, 2 Jan 2006 15:04:05 MST",
	"Mon, 02 Jan 2006 15:04 -0700",
	"Mon, 2 Jan 2006 15:04 -0700",
}

// ParseRFC2822 parses an RFC 2822 formatted date, e.g.
// "Tue, 01 Nov 2016 13:23:12 +0630".
func ParseRFC2822(value string) (Moment, error) {
	value = strings.TrimSpace(value)
	var firstErr error
	for _, layout := range rfc2822Layouts {
		if t, err := time.Parse(layout, value); err == nil {
			m := New(t)
			m.creation = &CreationData{Input: value, Format: "RFC2822", Locale: GlobalLocale(), Valid: true}
			return m, nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	return Invalid(), firstErr
}

// ParseFormat parses value using a moment.js-style token format string in the
// Moment's global locale, interpreting zone-less input as UTC.
func ParseFormat(value, format string) (Moment, error) {
	return parseWithFormat(value, format, mustLocale(""), time.UTC, false)
}

// ParseFormatStrict parses value against format in strict mode: the value must
// match the format exactly, with no extra or missing characters.
func ParseFormatStrict(value, format string) (Moment, error) {
	return parseWithFormat(value, format, mustLocale(""), time.UTC, true)
}

// ParseFormatLocale parses value using format and the named locale (so that
// localized month and weekday names are recognized).
func ParseFormatLocale(value, format, locale string) (Moment, error) {
	return parseWithFormat(value, format, mustLocale(locale), time.UTC, false)
}

// ParseInLocation parses value using a moment-style token format, interpreting
// zone-less input as being in loc. A nil location is treated as UTC.
func ParseInLocation(value, format string, loc *time.Location) (Moment, error) {
	if loc == nil {
		loc = time.UTC
	}
	return parseWithFormat(value, format, mustLocale(""), loc, false)
}

// ParseFormats tries each format in turn (moment.js's multiple-format parsing)
// and returns the first that succeeds. It returns an error if none match.
func ParseFormats(value string, formats []string) (Moment, error) {
	var firstErr error
	for _, f := range formats {
		if m, err := parseWithFormat(value, f, mustLocale(""), time.UTC, false); err == nil {
			return m, nil
		} else if firstErr == nil {
			firstErr = err
		}
	}
	if firstErr == nil {
		firstErr = errors.New("moment: no formats provided")
	}
	return Invalid(), firstErr
}

// ParseLayout parses value using a raw Go reference layout.
func ParseLayout(value, layout string) (Moment, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return Invalid(), err
	}
	return New(t), nil
}

// parseFields accumulates components as a format string is consumed.
type parseFields struct {
	year, month, day           int
	hour, minute, second, nsec int
	dayOfYear                  int
	yearSet, monthSet, daySet  bool
	hour12                     bool
	pmSet, pm                  bool
	offset                     int
	offsetSet                  bool
	unixMillis                 int64
	unixSet                    bool
}

// parseWithFormat parses value against a moment token format. base supplies the
// location for zone-less input; strict requires an exact match.
func parseWithFormat(value, format string, loc *Locale, base *time.Location, strict bool) (Moment, error) {
	var f parseFields
	f.month = 1
	f.day = 1
	vi := 0

	fail := func() (Moment, error) {
		return Invalid(), errors.New("moment: value " + strconv.Quote(value) + " does not match format " + strconv.Quote(format))
	}

	readInt := func(min, max int) (int, bool) {
		start := vi
		for vi < len(value) && vi-start < max && value[vi] >= '0' && value[vi] <= '9' {
			vi++
		}
		if vi-start < min {
			return 0, false
		}
		n, _ := strconv.Atoi(value[start:vi])
		return n, true
	}

	// widths select strict vs lenient digit counts.
	intTok := func(width int) (int, bool) {
		if strict {
			return readInt(width, width)
		}
		return readInt(1, width)
	}

	matchName := func(names []string) (int, bool) {
		bestIdx, bestLen := -1, 0
		for i, name := range names {
			if len(name) > bestLen && len(value)-vi >= len(name) &&
				strings.EqualFold(value[vi:vi+len(name)], name) {
				bestIdx, bestLen = i, len(name)
			}
		}
		if bestIdx < 0 {
			return 0, false
		}
		vi += bestLen
		return bestIdx, true
	}

	skipOrdinal := func() {
		for vi < len(value) {
			c := value[vi]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				vi++
			} else {
				break
			}
		}
	}

	handleToken := func(tok string) bool {
		switch tok {
		case "YYYY":
			n, ok := readInt(cond(strict, 4, 1), 4)
			if !ok {
				return false
			}
			f.year, f.yearSet = n, true
		case "YY":
			n, ok := readInt(cond(strict, 2, 1), 2)
			if !ok {
				return false
			}
			if n > 68 {
				n += 1900
			} else {
				n += 2000
			}
			f.year, f.yearSet = n, true
		case "Y":
			neg := false
			if vi < len(value) && (value[vi] == '-' || value[vi] == '+') {
				neg = value[vi] == '-'
				vi++
			}
			n, ok := readInt(1, 6)
			if !ok {
				return false
			}
			if neg {
				n = -n
			}
			f.year, f.yearSet = n, true
		case "gggg", "GGGG":
			if _, ok := readInt(cond(strict, 4, 1), 4); !ok {
				return false
			}
		case "gg", "GG":
			if _, ok := intTok(2); !ok {
				return false
			}
		case "M", "MM":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.month, f.monthSet = n, true
		case "Mo":
			n, ok := readInt(1, 2)
			if !ok {
				return false
			}
			skipOrdinal()
			f.month, f.monthSet = n, true
		case "MMM":
			i, ok := matchName(loc.MonthsShort)
			if !ok {
				return false
			}
			f.month, f.monthSet = i+1, true
		case "MMMM":
			i, ok := matchName(loc.Months)
			if !ok {
				return false
			}
			f.month, f.monthSet = i+1, true
		case "Q", "Qo":
			n, ok := readInt(1, 1)
			if !ok {
				return false
			}
			if tok == "Qo" {
				skipOrdinal()
			}
			f.month, f.monthSet = (n-1)*3+1, true
		case "D", "DD":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.day, f.daySet = n, true
		case "Do":
			n, ok := readInt(1, 2)
			if !ok {
				return false
			}
			skipOrdinal()
			f.day, f.daySet = n, true
		case "DDD", "DDDD":
			n, ok := readInt(cond(strict && tok == "DDDD", 3, 1), 3)
			if !ok {
				return false
			}
			f.dayOfYear = n
		case "DDDo":
			n, ok := readInt(1, 3)
			if !ok {
				return false
			}
			skipOrdinal()
			f.dayOfYear = n
		case "H", "HH":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.hour = n
		case "h", "hh":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.hour, f.hour12 = n, true
		case "k", "kk":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.hour = n % 24
		case "m", "mm":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.minute = n
		case "s", "ss":
			n, ok := intTok(2)
			if !ok {
				return false
			}
			f.second = n
		case "A", "a":
			if lower, ok := matchMeridiem(value, &vi, loc); ok {
				f.pmSet, f.pm = true, lower
			} else {
				return false
			}
		case "d", "e", "E":
			if _, ok := readInt(1, 1); !ok {
				return false
			}
		case "do":
			if _, ok := readInt(1, 1); !ok {
				return false
			}
			skipOrdinal()
		case "dd", "ddd", "dddd":
			names := loc.WeekdaysMin
			switch tok {
			case "ddd":
				names = loc.WeekdaysShort
			case "dddd":
				names = loc.Weekdays
			}
			if _, ok := matchName(names); !ok {
				return false
			}
		case "w", "ww", "W", "WW":
			if _, ok := intTok(2); !ok {
				return false
			}
		case "wo", "Wo":
			if _, ok := readInt(1, 2); !ok {
				return false
			}
			skipOrdinal()
		case "z", "zz":
			for vi < len(value) {
				c := value[vi]
				if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
					vi++
				} else {
					break
				}
			}
		case "Z", "ZZ":
			off, ok := parseOffset(value, &vi)
			if !ok {
				return false
			}
			f.offset, f.offsetSet = off, true
		case "X":
			start := vi
			for vi < len(value) && (value[vi] == '-' || (value[vi] >= '0' && value[vi] <= '9') || value[vi] == '.') {
				vi++
			}
			sec, err := strconv.ParseFloat(value[start:vi], 64)
			if err != nil {
				return false
			}
			f.unixMillis, f.unixSet = int64(sec*1000), true
		case "x":
			neg := false
			if vi < len(value) && value[vi] == '-' {
				neg = true
				vi++
			}
			n, ok := readInt(1, 18)
			if !ok {
				return false
			}
			if neg {
				n = -n
			}
			f.unixMillis, f.unixSet = int64(n), true
		default:
			return false
		}
		return true
	}

	for i := 0; i < len(format); {
		if format[i] == '[' {
			if end := strings.IndexByte(format[i+1:], ']'); end >= 0 {
				lit := format[i+1 : i+1+end]
				if !strings.HasPrefix(value[vi:], lit) {
					return fail()
				}
				vi += len(lit)
				i += end + 2
				continue
			}
		}
		if format[i] == 'S' {
			j := i
			for j < len(format) && format[j] == 'S' {
				j++
			}
			start := vi
			for vi < len(value) && vi-start < 9 && value[vi] >= '0' && value[vi] <= '9' {
				vi++
			}
			if digits := value[start:vi]; digits != "" {
				for len(digits) < 9 {
					digits += "0"
				}
				n, _ := strconv.Atoi(digits[:9])
				f.nsec = n
			}
			i = j
			continue
		}
		matched := false
		for _, tok := range formatTokens {
			if strings.HasPrefix(format[i:], tok) {
				if !handleToken(tok) {
					return fail()
				}
				i += len(tok)
				matched = true
				break
			}
		}
		if matched {
			continue
		}
		// Literal format character.
		c := format[i]
		if c == ' ' {
			for vi < len(value) && value[vi] == ' ' {
				vi++
			}
			i++
			continue
		}
		if vi >= len(value) || value[vi] != c {
			return fail()
		}
		vi++
		i++
	}

	if strict && vi != len(value) {
		return fail()
	}

	m, err := f.build(base)
	if err != nil {
		return Invalid(), err
	}
	m.loc = loc
	m.creation = &CreationData{Input: value, Format: format, Locale: loc.Name, IsUTC: base == time.UTC, Valid: true}
	return m, nil
}

// build assembles a Moment from the accumulated fields.
func (f parseFields) build(base *time.Location) (Moment, error) {
	if f.unixSet {
		return New(time.UnixMilli(f.unixMillis).In(base)), nil
	}
	loc := base
	if f.offsetSet {
		loc = time.FixedZone("", f.offset)
	}
	hour := f.hour
	if f.hour12 && f.pmSet {
		hour %= 12
		if f.pm {
			hour += 12
		}
	}
	year := f.year
	if !f.yearSet {
		year = time.Now().Year()
	}
	if f.dayOfYear > 0 && !f.monthSet && !f.daySet {
		t := time.Date(year, time.January, 1, hour, f.minute, f.second, f.nsec, loc).AddDate(0, 0, f.dayOfYear-1)
		return New(t), nil
	}
	t := time.Date(year, time.Month(f.month), f.day, hour, f.minute, f.second, f.nsec, loc)
	return New(t), nil
}

// cond returns a when c is true, otherwise b (a small ternary helper).
func cond(c bool, a, b int) int {
	if c {
		return a
	}
	return b
}

// matchMeridiem consumes a meridiem marker at value[*vi], returning whether it
// denotes PM. It recognizes the locale's own AM/PM strings as well as the ASCII
// forms.
func matchMeridiem(value string, vi *int, loc *Locale) (isPM, ok bool) {
	candidates := []struct {
		text string
		pm   bool
	}{
		{loc.PM, true}, {loc.AM, false},
		{"PM", true}, {"AM", false}, {"P", true}, {"A", false},
	}
	best := -1
	bestLen := 0
	for i, c := range candidates {
		if c.text == "" {
			continue
		}
		if len(c.text) > bestLen && len(value)-*vi >= len(c.text) &&
			strings.EqualFold(value[*vi:*vi+len(c.text)], c.text) {
			best, bestLen = i, len(c.text)
		}
	}
	if best < 0 {
		return false, false
	}
	*vi += bestLen
	return candidates[best].pm, true
}

// parseOffset consumes a zone offset (Z, ±HH:MM or ±HHMM) at value[*vi] and
// returns it in seconds east of UTC.
func parseOffset(value string, vi *int) (int, bool) {
	if *vi < len(value) && (value[*vi] == 'Z' || value[*vi] == 'z') {
		*vi++
		return 0, true
	}
	if *vi >= len(value) || (value[*vi] != '+' && value[*vi] != '-') {
		return 0, false
	}
	sign := 1
	if value[*vi] == '-' {
		sign = -1
	}
	*vi++
	if len(value)-*vi < 2 {
		return 0, false
	}
	hh, err := strconv.Atoi(value[*vi : *vi+2])
	if err != nil {
		return 0, false
	}
	*vi += 2
	if *vi < len(value) && value[*vi] == ':' {
		*vi++
	}
	mm := 0
	if len(value)-*vi >= 2 && value[*vi] >= '0' && value[*vi] <= '9' {
		mm, _ = strconv.Atoi(value[*vi : *vi+2])
		*vi += 2
	}
	return sign * (hh*3600 + mm*60), true
}
