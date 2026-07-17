package moment

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// formatTokens lists every supported moment.js format token, ordered
// longest-first so the scanner matches greedily and unambiguously. Runs of the
// fractional-second token "S" are handled separately by the scanner.
var formatTokens = []string{
	"LLLL", "LTS", "LLL", "LL", "LT", "L",
	"MMMM", "MMM", "MM", "Mo", "M",
	"Qo", "Q",
	"DDDo", "DDDD", "DDD", "Do", "DD", "D",
	"dddd", "ddd", "dd", "do", "d",
	"e", "E",
	"wo", "ww", "w",
	"Wo", "WW", "W",
	"YYYY", "YY", "Y",
	"gggg", "gg",
	"GGGG", "GG",
	"A", "a",
	"HH", "H", "hh", "h", "kk", "k",
	"mm", "m",
	"ss", "s",
	"zz", "z", "ZZ", "Z",
	"X", "x",
}

// SupportedTokens returns the moment-style format tokens understood by Format
// and ParseFormat, ordered longest-first. In addition, a run of "S" characters
// renders that many digits of the fractional second.
func SupportedTokens() []string {
	out := make([]string, len(formatTokens))
	copy(out, formatTokens)
	return out
}

// Format renders the Moment using a moment.js-style token string, resolving
// tokens against the Moment's effective locale. Text wrapped in square brackets
// is emitted literally, e.g. "YYYY [year]".
func (m Moment) Format(format string) string {
	if m.invalid {
		return "Invalid date"
	}
	return m.formatWith(format, m.localeOf())
}

// formatWith renders format using the given locale, expanding long-date tokens
// recursively.
func (m Moment) formatWith(format string, loc *Locale) string {
	var b strings.Builder
	for i := 0; i < len(format); {
		if format[i] == '[' {
			if end := strings.IndexByte(format[i+1:], ']'); end >= 0 {
				b.WriteString(format[i+1 : i+1+end])
				i += end + 2
				continue
			}
		}
		if format[i] == 'S' {
			j := i
			for j < len(format) && format[j] == 'S' {
				j++
			}
			b.WriteString(fractionalDigits(m.t.Nanosecond(), j-i))
			i = j
			continue
		}
		matched := false
		for _, tok := range formatTokens {
			if strings.HasPrefix(format[i:], tok) {
				b.WriteString(m.renderToken(tok, loc))
				i += len(tok)
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

// fractionalDigits returns the first n digits of the fractional-second value
// represented by nsec (0–999999999), matching moment's S…SSSSSSSSS tokens.
func fractionalDigits(nsec, n int) string {
	full := fmt.Sprintf("%09d", nsec)
	if n <= 9 {
		return full[:n]
	}
	return full + strings.Repeat("0", n-9)
}

// renderToken renders a single recognized token against loc.
func (m Moment) renderToken(tok string, loc *Locale) string {
	t := m.t
	switch tok {
	case "LTS", "LT", "L", "LL", "LLL", "LLLL":
		return m.formatWith(loc.longDateFormat(tok), loc)
	case "M":
		return strconv.Itoa(int(t.Month()))
	case "Mo":
		return loc.ordinal(int(t.Month()), "M")
	case "MM":
		return fmt.Sprintf("%02d", int(t.Month()))
	case "MMM":
		return loc.MonthsShort[int(t.Month())-1]
	case "MMMM":
		return loc.Months[int(t.Month())-1]
	case "Q":
		return strconv.Itoa(m.Quarter())
	case "Qo":
		return loc.ordinal(m.Quarter(), "Q")
	case "D":
		return strconv.Itoa(t.Day())
	case "Do":
		return loc.ordinal(t.Day(), "D")
	case "DD":
		return fmt.Sprintf("%02d", t.Day())
	case "DDD":
		return strconv.Itoa(t.YearDay())
	case "DDDo":
		return loc.ordinal(t.YearDay(), "DDD")
	case "DDDD":
		return fmt.Sprintf("%03d", t.YearDay())
	case "d":
		return strconv.Itoa(int(t.Weekday()))
	case "do":
		return loc.ordinal(int(t.Weekday()), "d")
	case "dd":
		return loc.WeekdaysMin[int(t.Weekday())]
	case "ddd":
		return loc.WeekdaysShort[int(t.Weekday())]
	case "dddd":
		return loc.Weekdays[int(t.Weekday())]
	case "e":
		return strconv.Itoa(m.localeWeekday(loc))
	case "E":
		return strconv.Itoa(m.ISOWeekday())
	case "w":
		w, _ := m.localeWeek(loc)
		return strconv.Itoa(w)
	case "wo":
		w, _ := m.localeWeek(loc)
		return loc.ordinal(w, "w")
	case "ww":
		w, _ := m.localeWeek(loc)
		return fmt.Sprintf("%02d", w)
	case "W":
		_, w := t.ISOWeek()
		return strconv.Itoa(w)
	case "Wo":
		_, w := t.ISOWeek()
		return loc.ordinal(w, "W")
	case "WW":
		_, w := t.ISOWeek()
		return fmt.Sprintf("%02d", w)
	case "YYYY":
		return fmt.Sprintf("%04d", t.Year())
	case "YY":
		return fmt.Sprintf("%02d", t.Year()%100)
	case "Y":
		return strconv.Itoa(t.Year())
	case "gg":
		_, gy := m.localeWeek(loc)
		return fmt.Sprintf("%02d", gy%100)
	case "gggg":
		_, gy := m.localeWeek(loc)
		return fmt.Sprintf("%04d", gy)
	case "GG":
		gy, _ := t.ISOWeek()
		return fmt.Sprintf("%02d", gy%100)
	case "GGGG":
		gy, _ := t.ISOWeek()
		return fmt.Sprintf("%04d", gy)
	case "A":
		return loc.meridiem(t.Hour(), t.Minute(), false)
	case "a":
		return loc.meridiem(t.Hour(), t.Minute(), true)
	case "H":
		return strconv.Itoa(t.Hour())
	case "HH":
		return fmt.Sprintf("%02d", t.Hour())
	case "h":
		return strconv.Itoa(hour12(t.Hour()))
	case "hh":
		return fmt.Sprintf("%02d", hour12(t.Hour()))
	case "k":
		return strconv.Itoa(hour24k(t.Hour()))
	case "kk":
		return fmt.Sprintf("%02d", hour24k(t.Hour()))
	case "m":
		return strconv.Itoa(t.Minute())
	case "mm":
		return fmt.Sprintf("%02d", t.Minute())
	case "s":
		return strconv.Itoa(t.Second())
	case "ss":
		return fmt.Sprintf("%02d", t.Second())
	case "z", "zz":
		name, _ := t.Zone()
		return name
	case "Z":
		return offsetString(t, true)
	case "ZZ":
		return offsetString(t, false)
	case "X":
		return strconv.FormatInt(t.Unix(), 10)
	case "x":
		return strconv.FormatInt(t.UnixMilli(), 10)
	}
	return tok
}

// longDateFormat returns the pattern for a long-date token (LT, LTS, L, LL,
// LLL, LLLL), or the token itself if unknown.
func (l *Locale) longDateFormat(tok string) string {
	switch tok {
	case "LT":
		return l.LongDateFormats.LT
	case "LTS":
		return l.LongDateFormats.LTS
	case "L":
		return l.LongDateFormats.L
	case "LL":
		return l.LongDateFormats.LL
	case "LLL":
		return l.LongDateFormats.LLL
	case "LLLL":
		return l.LongDateFormats.LLLL
	}
	return tok
}

// hour12 converts a 24-hour value to the 1–12 range used by the h/hh tokens.
func hour12(h int) int {
	h %= 12
	if h == 0 {
		return 12
	}
	return h
}

// hour24k converts a 24-hour value to the 1–24 range used by the k/kk tokens.
func hour24k(h int) int {
	if h == 0 {
		return 24
	}
	return h
}

// offsetString renders the Moment's zone offset. When colon is true it uses the
// "+07:00" form (token Z); otherwise the "+0700" form (token ZZ).
func offsetString(t time.Time, colon bool) string {
	_, off := t.Zone()
	sign := "+"
	if off < 0 {
		sign = "-"
		off = -off
	}
	h := off / 3600
	mnt := (off % 3600) / 60
	if colon {
		return fmt.Sprintf("%s%02d:%02d", sign, h, mnt)
	}
	return fmt.Sprintf("%s%02d%02d", sign, h, mnt)
}

// FormatLayout renders the Moment using a raw Go reference layout
// (for example "2006-01-02 15:04:05").
func (m Moment) FormatLayout(layout string) string {
	return m.t.Format(layout)
}

// ISO returns the RFC3339 (ISO-8601) representation of the Moment, preserving
// its location and any fractional seconds.
func (m Moment) ISO() string {
	return m.t.Format(time.RFC3339Nano)
}

// ToISOString returns the moment.js toISOString representation: the instant in
// UTC with millisecond precision and a trailing "Z", e.g.
// "2017-07-14T02:40:00.000Z".
func (m Moment) ToISOString() string {
	return m.t.UTC().Format("2006-01-02T15:04:05.000Z07:00")
}
