package moment

import "time"

// safeIndex returns names[i] when i is in range, and the empty string
// otherwise. It guards the locale-listing helpers against short slices.
func lqSafeIndex(names []string, i int) string {
	if i < 0 || i >= len(names) {
		return ""
	}
	return names[i]
}

// Months returns a copy of the twelve month names (January first) for the named
// locale, mirroring moment.js's moment.months(). An unknown locale falls back
// to the global default.
func Months(locale string) []string {
	return append([]string(nil), mustLocale(locale).Months...)
}

// MonthsShort returns a copy of the twelve abbreviated month names (January
// first) for the named locale, mirroring moment.js's moment.monthsShort().
func MonthsShort(locale string) []string {
	return append([]string(nil), mustLocale(locale).MonthsShort...)
}

// Weekdays returns a copy of the seven weekday names (Sunday first, index 0 ==
// Sunday) for the named locale, mirroring moment.js's moment.weekdays().
func Weekdays(locale string) []string {
	return append([]string(nil), mustLocale(locale).Weekdays...)
}

// WeekdaysShort returns a copy of the seven abbreviated weekday names (Sunday
// first) for the named locale, mirroring moment.js's moment.weekdaysShort().
func WeekdaysShort(locale string) []string {
	return append([]string(nil), mustLocale(locale).WeekdaysShort...)
}

// WeekdaysMin returns a copy of the seven minimal weekday names (Sunday first)
// for the named locale, mirroring moment.js's moment.weekdaysMin().
func WeekdaysMin(locale string) []string {
	return append([]string(nil), mustLocale(locale).WeekdaysMin...)
}

// MonthName returns the full name of month m in the named locale.
func MonthName(locale string, m time.Month) string {
	return lqSafeIndex(mustLocale(locale).Months, int(m)-1)
}

// MonthShortName returns the abbreviated name of month m in the named locale.
func MonthShortName(locale string, m time.Month) string {
	return lqSafeIndex(mustLocale(locale).MonthsShort, int(m)-1)
}

// WeekdayName returns the full name of weekday w in the named locale.
func WeekdayName(locale string, w time.Weekday) string {
	return lqSafeIndex(mustLocale(locale).Weekdays, int(w))
}

// WeekdayShortName returns the abbreviated name of weekday w in the named
// locale.
func WeekdayShortName(locale string, w time.Weekday) string {
	return lqSafeIndex(mustLocale(locale).WeekdaysShort, int(w))
}

// WeekdayMinName returns the minimal name of weekday w in the named locale.
func WeekdayMinName(locale string, w time.Weekday) string {
	return lqSafeIndex(mustLocale(locale).WeekdaysMin, int(w))
}

// Ordinal renders n as an ordinal in the named locale, for example "1st" in
// English or "1." in German, mirroring moment.js's localeData().ordinal.
func Ordinal(locale string, n int) string {
	return mustLocale(locale).ordinal(n, "D")
}

// FirstDayOfWeek returns the named locale's first day of the week (0 == Sunday,
// 1 == Monday), mirroring moment.js's localeData().firstDayOfWeek.
func FirstDayOfWeek(locale string) int {
	dow, _ := mustLocale(locale).weekRules()
	return dow
}

// FirstWeekContainsDate returns the January date that always falls in week 1
// under the named locale's week-numbering rules (moment's "doy"), mirroring
// localeData().firstDayOfYear.
func FirstWeekContainsDate(locale string) int {
	_, doy := mustLocale(locale).weekRules()
	return doy
}

// LongDateFormat returns the named locale's pattern for one of the long-date
// tokens LT, LTS, L, LL, LLL or LLLL, mirroring moment.js's
// localeData().longDateFormat. An unknown token yields the empty string.
func LongDateFormat(locale, token string) string {
	f := mustLocale(locale).LongDateFormats
	switch token {
	case "LT":
		return f.LT
	case "LTS":
		return f.LTS
	case "L":
		return f.L
	case "LL":
		return f.LL
	case "LLL":
		return f.LLL
	case "LLLL":
		return f.LLLL
	default:
		return ""
	}
}

// Meridiem returns the named locale's meridiem string for the given clock time;
// isLower requests the lower-case form. It mirrors moment.js's
// localeData().meridiem and honours locales (such as Chinese and Japanese) that
// use time-of-day words instead of AM/PM.
func Meridiem(locale string, hour, minute int, isLower bool) string {
	return mustLocale(locale).meridiem(hour, minute, isLower)
}

// NormalizeUnit resolves a moment.js unit alias (plural or short form such as
// "days", "d" or "M") to its canonical Unit, reporting whether it was
// recognized. It is the analogue of moment.js's moment.normalizeUnits.
func NormalizeUnit(unit string) (Unit, bool) {
	u := normalizeUnit(Unit(unit))
	switch u {
	case Year, Quarter, Month, Week, ISOWeek, Day, Date, DayOfYear, Hour, Minute, Second, Millisecond:
		return u, true
	default:
		return u, false
	}
}
