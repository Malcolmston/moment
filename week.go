package moment

import "time"

// floorDiv returns the floor of a/b for a positive divisor b, unlike Go's "/"
// operator which truncates toward zero.
func floorDiv(a, b int) int {
	q := a / b
	if (a%b != 0) && ((a < 0) != (b < 0)) {
		q--
	}
	return q
}

// daysInYearOf returns 366 for a leap year and 365 otherwise.
func daysInYearOf(year int) int {
	if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
		return 366
	}
	return 365
}

// firstWeekOffset implements moment's firstWeekOffset for the given week rules:
// dow is the first day of the week (0 == Sunday) and doy selects which January
// date is always in week 1.
func firstWeekOffset(year, dow, doy int) int {
	fwd := 7 + dow - doy
	janFwd := time.Date(year, time.January, fwd, 0, 0, 0, 0, time.UTC)
	fwdlw := (7 + int(janFwd.Weekday()) - dow) % 7
	return -fwdlw + fwd - 1
}

// weeksInYearOf returns the number of weeks in year under the given week rules.
func weeksInYearOf(year, dow, doy int) int {
	wo := firstWeekOffset(year, dow, doy)
	won := firstWeekOffset(year+1, dow, doy)
	return (daysInYearOf(year) - wo + won) / 7
}

// weekOfYear computes the locale week number and week-year for t under the given
// rules, following moment's algorithm.
func weekOfYear(t time.Time, dow, doy int) (week, weekYear int) {
	year := t.Year()
	wo := firstWeekOffset(year, dow, doy)
	week = floorDiv(t.YearDay()-wo-1, 7) + 1
	switch {
	case week < 1:
		weekYear = year - 1
		week += weeksInYearOf(weekYear, dow, doy)
	case week > weeksInYearOf(year, dow, doy):
		week -= weeksInYearOf(year, dow, doy)
		weekYear = year + 1
	default:
		weekYear = year
	}
	return week, weekYear
}

// localeWeek returns the Moment's week number and week-year under loc's rules.
func (m Moment) localeWeek(loc *Locale) (week, weekYear int) {
	dow, doy := loc.weekRules()
	return weekOfYear(m.t, dow, doy)
}

// localeWeekday returns the Moment's weekday index relative to loc's first day
// of the week (moment's "e" token: 0 == the locale's first weekday).
func (m Moment) localeWeekday(loc *Locale) int {
	dow, _ := loc.weekRules()
	return (int(m.t.Weekday()) - dow + 7) % 7
}
