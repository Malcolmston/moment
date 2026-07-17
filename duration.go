package moment

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Duration is a length of time, the moment.js Duration analogue. It stores
// months, days and milliseconds separately (because months and days have
// variable lengths) and offers conversion (As), component access (Get),
// arithmetic, ISO-8601 serialization and locale-aware humanization.
type Duration struct {
	months int
	days   int
	ms     int64
	loc    *Locale
}

// NewDuration returns a Duration of n of the given unit. Recognized units are
// Year, Quarter, Month, Week, Day, Hour, Minute, Second and Millisecond (plus
// their moment.js aliases).
func NewDuration(n int, unit Unit) Duration {
	var d Duration
	switch normalizeUnit(unit) {
	case Year:
		d.months = n * 12
	case Quarter:
		d.months = n * 3
	case Month:
		d.months = n
	case Week, ISOWeek:
		d.days = n * 7
	case Day, Date, DayOfYear:
		d.days = n
	case Hour:
		d.ms = int64(n) * int64(time.Hour/time.Millisecond)
	case Minute:
		d.ms = int64(n) * int64(time.Minute/time.Millisecond)
	case Second:
		d.ms = int64(n) * 1000
	case Millisecond:
		d.ms = int64(n)
	}
	return d
}

// DurationFromTime returns a Duration equal to the given time.Duration, stored
// entirely as milliseconds.
func DurationFromTime(td time.Duration) Duration {
	return Duration{ms: td.Milliseconds()}
}

// DurationBetween returns the Duration from from to to (that is, to − from),
// measured in milliseconds.
func DurationBetween(from, to Moment) Duration {
	return Duration{ms: to.t.UnixMilli() - from.t.UnixMilli()}
}

// daysToMonths and monthsToDays convert between days and months using the mean
// Gregorian month length (146097/4800 days), matching moment.js.
func daysToMonths(days float64) float64   { return days * 4800 / 146097 }
func monthsToDays(months float64) float64 { return months * 146097 / 4800 }

// absFloor rounds toward zero's floor: floor for positive, ceil for negative.
func absFloor(x float64) int {
	if x < 0 {
		return int(math.Ceil(x))
	}
	return int(math.Floor(x))
}

// absCeil rounds away from zero: ceil for positive, floor for negative.
func absCeil(x float64) int {
	if x < 0 {
		return int(math.Floor(x))
	}
	return int(math.Ceil(x))
}

// As returns the whole Duration expressed in the given unit as a float64.
func (d Duration) As(unit Unit) float64 {
	u := normalizeUnit(unit)
	switch u {
	case Year, Quarter, Month:
		days := float64(d.days) + float64(d.ms)/864e5
		months := float64(d.months) + daysToMonths(days)
		switch u {
		case Year:
			return months / 12
		case Quarter:
			return months / 3
		default:
			return months
		}
	}
	days := float64(d.days) + math.Round(monthsToDays(float64(d.months)))
	switch u {
	case Week, ISOWeek:
		return days/7 + float64(d.ms)/6048e5
	case Day, Date, DayOfYear:
		return days + float64(d.ms)/864e5
	case Hour:
		return days*24 + float64(d.ms)/36e5
	case Minute:
		return days*1440 + float64(d.ms)/6e4
	case Second:
		return days*86400 + float64(d.ms)/1000
	case Millisecond:
		return math.Floor(days*864e5) + float64(d.ms)
	}
	return 0
}

// AsYears returns the Duration in fractional years.
func (d Duration) AsYears() float64 { return d.As(Year) }

// AsMonths returns the Duration in fractional months.
func (d Duration) AsMonths() float64 { return d.As(Month) }

// AsWeeks returns the Duration in fractional weeks.
func (d Duration) AsWeeks() float64 { return d.As(Week) }

// AsDays returns the Duration in fractional days.
func (d Duration) AsDays() float64 { return d.As(Day) }

// AsHours returns the Duration in fractional hours.
func (d Duration) AsHours() float64 { return d.As(Hour) }

// AsMinutes returns the Duration in fractional minutes.
func (d Duration) AsMinutes() float64 { return d.As(Minute) }

// AsSeconds returns the Duration in fractional seconds.
func (d Duration) AsSeconds() float64 { return d.As(Second) }

// AsMilliseconds returns the Duration in whole milliseconds.
func (d Duration) AsMilliseconds() float64 { return d.As(Millisecond) }

// ToDuration returns the Duration as a time.Duration (nanosecond precision from
// its total-milliseconds length).
func (d Duration) ToDuration() time.Duration {
	return time.Duration(int64(d.As(Millisecond))) * time.Millisecond
}

// bubble normalizes the Duration into calendar components, following moment's
// bubbling algorithm, and returns them.
func (d Duration) bubble() (years, months, days, hours, minutes, seconds, milliseconds int) {
	ms := float64(d.ms)
	dys := d.days
	mos := d.months
	if !((ms >= 0 && dys >= 0 && mos >= 0) || (ms <= 0 && dys <= 0 && mos <= 0)) {
		ms += float64(absCeil(monthsToDays(float64(mos))+float64(dys))) * 864e5
		dys = 0
		mos = 0
	}
	milliseconds = int(math.Mod(ms, 1000))
	secTotal := absFloor(ms / 1000)
	seconds = secTotal % 60
	minTotal := absFloor(float64(secTotal) / 60)
	minutes = minTotal % 60
	hourTotal := absFloor(float64(minTotal) / 60)
	hours = hourTotal % 24
	dys += absFloor(float64(hourTotal) / 24)
	monthsFromDays := absFloor(daysToMonths(float64(dys)))
	mos += monthsFromDays
	dys -= absCeil(monthsToDays(float64(monthsFromDays)))
	years = absFloor(float64(mos) / 12)
	months = mos % 12
	days = dys
	return
}

// Get returns the Duration's value for the given unit as a normalized component
// (for example, Get(Minute) is 0–59, and Get(Month) is 0–11).
func (d Duration) Get(unit Unit) int {
	y, mo, dd, h, mi, s, mil := d.bubble()
	switch normalizeUnit(unit) {
	case Year:
		return y
	case Month:
		return mo
	case Week, ISOWeek:
		return absFloor(float64(dd) / 7)
	case Day, Date, DayOfYear:
		return dd
	case Hour:
		return h
	case Minute:
		return mi
	case Second:
		return s
	case Millisecond:
		return mil
	}
	return 0
}

// Years returns the years component (Get(Year)).
func (d Duration) Years() int { return d.Get(Year) }

// Months returns the months component (0–11).
func (d Duration) Months() int { return d.Get(Month) }

// Weeks returns the whole-weeks component derived from the days component.
func (d Duration) Weeks() int { return d.Get(Week) }

// Days returns the days component.
func (d Duration) Days() int { return d.Get(Day) }

// Hours returns the hours component (0–23).
func (d Duration) Hours() int { return d.Get(Hour) }

// Minutes returns the minutes component (0–59).
func (d Duration) Minutes() int { return d.Get(Minute) }

// Seconds returns the seconds component (0–59).
func (d Duration) Seconds() int { return d.Get(Second) }

// Milliseconds returns the milliseconds component (0–999).
func (d Duration) Milliseconds() int { return d.Get(Millisecond) }

// Add returns the sum of the two Durations, component by component.
func (d Duration) Add(other Duration) Duration {
	d.months += other.months
	d.days += other.days
	d.ms += other.ms
	return d
}

// Subtract returns d minus other, component by component.
func (d Duration) Subtract(other Duration) Duration {
	d.months -= other.months
	d.days -= other.days
	d.ms -= other.ms
	return d
}

// Abs returns the Duration with all components made non-negative.
func (d Duration) Abs() Duration {
	if d.months < 0 {
		d.months = -d.months
	}
	if d.days < 0 {
		d.days = -d.days
	}
	if d.ms < 0 {
		d.ms = -d.ms
	}
	return d
}

// Clone returns an independent copy of the Duration.
func (d Duration) Clone() Duration { return d }

// Locale returns a copy of the Duration bound to the named locale, affecting
// Humanize. An unknown name falls back to the global default.
func (d Duration) Locale(name string) Duration {
	if l, ok := LookupLocale(name); ok {
		d.loc = l
	} else {
		d.loc = nil
	}
	return d
}

// localeOf returns the Duration's effective locale.
func (d Duration) localeOf() *Locale {
	if d.loc != nil {
		return d.loc
	}
	return mustLocale("")
}

// relThresholds holds the relative-time cutoffs used by Humanize. A zero w
// means the week unit is disabled (moment's null).
type relThresholds struct {
	ss, s, m, h, d, w, M int
}

var (
	thresholdMu       sync.RWMutex
	currentThresholds = relThresholds{ss: 44, s: 45, m: 45, h: 22, d: 26, w: 0, M: 11}
)

// SetRelativeTimeThreshold overrides the cutoff for one relative-time unit, as
// moment.relativeTimeThreshold does. Recognized units are "ss", "s", "m", "h",
// "d", "w" and "M". It reports whether the unit name was recognized.
func SetRelativeTimeThreshold(unit string, limit int) bool {
	thresholdMu.Lock()
	defer thresholdMu.Unlock()
	switch unit {
	case "ss":
		currentThresholds.ss = limit
	case "s":
		currentThresholds.s = limit
	case "m":
		currentThresholds.m = limit
	case "h":
		currentThresholds.h = limit
	case "d":
		currentThresholds.d = limit
	case "w":
		currentThresholds.w = limit
	case "M":
		currentThresholds.M = limit
	default:
		return false
	}
	return true
}

// RelativeTimeThreshold returns the current cutoff for a relative-time unit, or
// -1 for an unknown unit.
func RelativeTimeThreshold(unit string) int {
	thresholdMu.RLock()
	defer thresholdMu.RUnlock()
	switch unit {
	case "ss":
		return currentThresholds.ss
	case "s":
		return currentThresholds.s
	case "m":
		return currentThresholds.m
	case "h":
		return currentThresholds.h
	case "d":
		return currentThresholds.d
	case "w":
		return currentThresholds.w
	case "M":
		return currentThresholds.M
	}
	return -1
}

func getThresholds() relThresholds {
	thresholdMu.RLock()
	defer thresholdMu.RUnlock()
	return currentThresholds
}

// roundInt rounds a float64 to the nearest int, halves away from zero.
func roundInt(f float64) int { return int(math.Round(f)) }

// Humanize renders the Duration as a locale-aware phrase such as "2 days". When
// withSuffix is true the locale's future/past wrapper is applied, yielding
// phrases like "in 2 days" or "2 days ago".
func (d Duration) Humanize(withSuffix bool) string {
	loc := d.localeOf()
	th := getThresholds()
	ad := d.Abs()
	seconds := roundInt(ad.As(Second))
	minutes := roundInt(ad.As(Minute))
	hours := roundInt(ad.As(Hour))
	days := roundInt(ad.As(Day))
	weeks := roundInt(ad.As(Week))
	months := roundInt(ad.As(Month))
	years := roundInt(ad.As(Year))

	var key string
	var number int
	switch {
	case seconds <= th.ss:
		key, number = "s", seconds
	case seconds < th.s:
		key, number = "ss", seconds
	case minutes <= 1:
		key, number = "m", minutes
	case minutes < th.m:
		key, number = "mm", minutes
	case hours <= 1:
		key, number = "h", hours
	case hours < th.h:
		key, number = "hh", hours
	case days <= 1:
		key, number = "d", days
	case days < th.d:
		key, number = "dd", days
	case th.w > 0 && weeks <= 1:
		key, number = "w", weeks
	case th.w > 0 && weeks < th.w:
		key, number = "ww", weeks
	case months <= 1:
		key, number = "M", months
	case months < th.M:
		key, number = "MM", months
	case years <= 1:
		key, number = "y", years
	default:
		key, number = "yy", years
	}

	future := d.signMillis() >= 0
	output := loc.relativeTimeString(number, !withSuffix, key, future)
	if withSuffix {
		output = loc.pastFuture(future, output)
	}
	return output
}

// signMillis returns a signed magnitude used to decide future vs past. Zero is
// treated as future.
func (d Duration) signMillis() int64 {
	total := d.ms + int64(d.days)*864e5 + int64(math.Round(monthsToDays(float64(d.months))))*864e5
	return total
}

// relativeTimeString renders a single relative-time key for the locale.
func (l *Locale) relativeTimeString(number int, withoutSuffix bool, key string, isFuture bool) string {
	if l.RelativeTime.Pluralize != nil {
		return l.RelativeTime.Pluralize(number, withoutSuffix, key, isFuture)
	}
	tmpl := l.relativeField(key)
	return strings.Replace(tmpl, "%d", strconv.Itoa(number), 1)
}

// relativeField returns the template string for a relative-time key.
func (l *Locale) relativeField(key string) string {
	switch key {
	case "s":
		return l.RelativeTime.Second
	case "ss":
		return l.RelativeTime.Seconds
	case "m":
		return l.RelativeTime.Minute
	case "mm":
		return l.RelativeTime.Minutes
	case "h":
		return l.RelativeTime.Hour
	case "hh":
		return l.RelativeTime.Hours
	case "d":
		return l.RelativeTime.Day
	case "dd":
		return l.RelativeTime.Days
	case "w":
		return l.RelativeTime.Week
	case "ww":
		return l.RelativeTime.Weeks
	case "M":
		return l.RelativeTime.Month
	case "MM":
		return l.RelativeTime.Months
	case "y":
		return l.RelativeTime.Year
	case "yy":
		return l.RelativeTime.Years
	}
	return ""
}

// pastFuture wraps output with the locale's future or past template.
func (l *Locale) pastFuture(future bool, output string) string {
	tmpl := l.RelativeTime.Past
	if future {
		tmpl = l.RelativeTime.Future
	}
	return strings.Replace(tmpl, "%s", output, 1)
}

// ISOString serializes the Duration as an ISO-8601 duration such as
// "P1Y2M10DT2H30M". A zero Duration is "P0D". The output uses the stored
// months, days and milliseconds without converting between them.
func (d Duration) ISOString() string {
	seconds := math.Abs(float64(d.ms)) / 1000
	days := d.days
	if days < 0 {
		days = -days
	}
	months := d.months
	if months < 0 {
		months = -months
	}
	minutes := absFloor(seconds / 60)
	hours := absFloor(float64(minutes) / 60)
	seconds = math.Mod(seconds, 60)
	minutes %= 60
	years := months / 12
	months %= 12

	total := d.As(Second)
	if total == 0 {
		return "P0D"
	}
	sign := ""
	if total < 0 {
		sign = "-"
	}

	var b strings.Builder
	b.WriteString(sign)
	b.WriteByte('P')
	if years != 0 {
		fmt.Fprintf(&b, "%dY", years)
	}
	if months != 0 {
		fmt.Fprintf(&b, "%dM", months)
	}
	if days != 0 {
		fmt.Fprintf(&b, "%dD", days)
	}
	if hours != 0 || minutes != 0 || seconds != 0 {
		b.WriteByte('T')
		if hours != 0 {
			fmt.Fprintf(&b, "%dH", hours)
		}
		if minutes != 0 {
			fmt.Fprintf(&b, "%dM", minutes)
		}
		if seconds != 0 {
			b.WriteString(strconv.FormatFloat(seconds, 'f', -1, 64))
			b.WriteByte('S')
		}
	}
	return b.String()
}

// String returns the Duration's ISO-8601 representation.
func (d Duration) String() string { return d.ISOString() }

var isoDurationRe = regexp.MustCompile(`^(-)?P(?:(\d+(?:\.\d+)?)Y)?(?:(\d+(?:\.\d+)?)M)?(?:(\d+(?:\.\d+)?)W)?(?:(\d+(?:\.\d+)?)D)?(?:T(?:(\d+(?:\.\d+)?)H)?(?:(\d+(?:\.\d+)?)M)?(?:(\d+(?:\.\d+)?)S)?)?$`)

// ParseDuration parses an ISO-8601 duration string (for example
// "P1Y2M10DT2H30M" or "PT1.5S") into a Duration. Weeks (W) are folded into
// days. It returns an error for input that is not a valid ISO-8601 duration.
func ParseDuration(s string) (Duration, error) {
	m := isoDurationRe.FindStringSubmatch(s)
	if m == nil || s == "P" || s == "-P" || s == "PT" {
		return Duration{}, errors.New("moment: invalid ISO-8601 duration: " + strconv.Quote(s))
	}
	// Require at least one component so bare "P" is rejected.
	hasAny := false
	for _, g := range m[2:] {
		if g != "" {
			hasAny = true
			break
		}
	}
	if !hasAny {
		return Duration{}, errors.New("moment: invalid ISO-8601 duration: " + strconv.Quote(s))
	}
	num := func(g string) float64 {
		if g == "" {
			return 0
		}
		v, _ := strconv.ParseFloat(g, 64)
		return v
	}
	years := num(m[2])
	months := num(m[3])
	weeks := num(m[4])
	days := num(m[5])
	hours := num(m[6])
	minutes := num(m[7])
	seconds := num(m[8])

	var d Duration
	d.months = int(years*12 + months)
	d.days = int(weeks*7 + days)
	d.ms = int64(math.Round(hours*3600000 + minutes*60000 + seconds*1000))
	if m[1] == "-" {
		d.months = -d.months
		d.days = -d.days
		d.ms = -d.ms
	}
	return d, nil
}
