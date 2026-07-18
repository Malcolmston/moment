package moment

// This file encodes known-answer test vectors taken directly from the upstream
// moment.js test suite (github.com/moment/moment, src/test/moment/*.js and
// src/test/duration/*.js). Each case pins this Go port to the exact behavior
// moment.js asserts. Vectors were transcribed from the live upstream sources.
//
// moment.js array construction moment([year, month, day, ...]) uses a 0-based
// month; the port's FromArray follows the same convention, so the array inputs
// below match the upstream tests verbatim. Month getters here are compared in
// the port's 1-based time.Month space by subtracting 1 to reach moment's
// 0-based value where the upstream assertion is 0-based.

import (
	"testing"
	"time"
)

// fromArr builds a UTC Moment from a moment.js-style [year, month0, day, ...]
// component slice.
func fromArr(parts ...int) Moment { return FromArray(parts) }

// month0 returns the port Moment's 0-based month, matching moment.js month().
func month0(m Moment) int { return int(m.Month()) - 1 }

// TestParityQuarterGetter mirrors quarter.js "library quarter getter".
func TestParityQuarterGetter(t *testing.T) {
	cases := []struct {
		y, mo, d int
		want     int
	}{
		{1985, 1, 4, 1},   // Feb  4 1985 is Q1
		{2029, 8, 18, 3},  // Sep 18 2029 is Q3
		{2013, 3, 24, 2},  // Apr 24 2013 is Q2
		{2015, 2, 5, 1},   // Mar  5 2015 is Q1
		{1970, 0, 2, 1},   // Jan  2 1970 is Q1
		{2001, 11, 12, 4}, // Dec 12 2001 is Q4
		{2000, 0, 2, 1},   // Jan  2 2000 is Q1
	}
	for _, c := range cases {
		if got := fromArr(c.y, c.mo, c.d).Quarter(); got != c.want {
			t.Errorf("quarter(%d-%d-%d) = %d, want %d", c.y, c.mo+1, c.d, got, c.want)
		}
	}
}

// TestParityQuarterSetter mirrors quarter.js "quarter setter singular",
// "only month changes" and the "bubble to next/previous year" cases.
func TestParityQuarterSetter(t *testing.T) {
	base := fromArr(2014, 4, 11) // May 11 2014
	for _, c := range []struct{ q, wantMonth0 int }{
		{2, 4}, {3, 7}, {1, 1}, {4, 10},
	} {
		if got := month0(base.SetQuarter(c.q)); got != c.wantMonth0 {
			t.Errorf("SetQuarter(%d).month() = %d, want %d", c.q, got, c.wantMonth0)
		}
	}

	// "quarter setter only month changes": q(4) keeps all other components.
	full := fromArr(2014, 4, 11, 1, 2, 3, 4).SetQuarter(4)
	if full.Year() != 2014 || month0(full) != 10 || full.Date() != 11 ||
		full.Hour() != 1 || full.Minute() != 2 || full.Second() != 3 || full.Millisecond() != 4 {
		t.Errorf("SetQuarter(4) altered non-month components: %s.%03d", full.Format("YYYY-MM-DD HH:mm:ss"), full.Millisecond())
	}

	// "bubble to next year": q(7) => 2015, month 7 (0-based).
	next := fromArr(2014, 4, 11, 1, 2, 3, 4).SetQuarter(7)
	if next.Year() != 2015 || month0(next) != 7 || next.Date() != 11 {
		t.Errorf("SetQuarter(7) = %s, want 2015 month0 7 date 11", next.Format("YYYY-MM-DD"))
	}

	// "bubble to previous year": q(-3) => 2013, month 1 (0-based).
	prev := fromArr(2014, 4, 11, 1, 2, 3, 4).SetQuarter(-3)
	if prev.Year() != 2013 || month0(prev) != 1 || prev.Date() != 11 {
		t.Errorf("SetQuarter(-3) = %s, want 2013 month0 1 date 11", prev.Format("YYYY-MM-DD"))
	}
}

// TestParityQuarterDiff mirrors quarter.js "quarter diff".
func TestParityQuarterDiff(t *testing.T) {
	p := func(s string) Moment { m, _ := Parse(s); return m }
	if got := p("2014-01-01").DiffInt(p("2014-04-01"), Quarter); got != -1 {
		t.Errorf("diff quarter = %d, want -1", got)
	}
	if got := p("2014-04-01").DiffInt(p("2014-01-01"), Quarter); got != 1 {
		t.Errorf("diff quarter = %d, want 1", got)
	}
	if got := p("2014-05-01").DiffInt(p("2014-01-01"), Quarter); got != 1 {
		t.Errorf("diff quarter = %d, want 1", got)
	}
	if got := p("2014-05-01").Diff(p("2014-01-01"), Quarter); got < 1.333 || got > 1.334 {
		t.Errorf("diff quarter float = %v, want ~1.3333", got)
	}
	if got := p("2015-01-01").DiffInt(p("2014-01-01"), Quarter); got != 4 {
		t.Errorf("diff quarter = %d, want 4", got)
	}
}

// TestParityDaysInMonth mirrors days_in_month.js "days in month" and the leap
// year cases.
func TestParityDaysInMonth(t *testing.T) {
	// Non-February months across a broad year range.
	days := [12]int{31, 0, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	for year := 1899; year < 2100; year++ {
		for mo := 0; mo < 12; mo++ {
			if mo == 1 {
				continue
			}
			if got := fromArr(year, mo).DaysInMonth(); got != days[mo] {
				t.Fatalf("daysInMonth(%d-%d) = %d, want %d", year, mo+1, got, days[mo])
			}
		}
	}
	// February across leap and common years.
	for _, c := range []struct {
		year, want int
	}{{2010, 28}, {2100, 28}, {2008, 29}, {2000, 29}} {
		if got := fromArr(c.year, 1).DaysInMonth(); got != c.want {
			t.Errorf("Feb %d daysInMonth = %d, want %d", c.year, got, c.want)
		}
	}
	// 2012 (leap) full month-length table, tested at first and last day.
	table := [12]int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	for i, want := range table {
		if got := fromArr(2012, i).DaysInMonth(); got != want {
			t.Errorf("2012-%02d first-day daysInMonth = %d, want %d", i+1, got, want)
		}
		if got := fromArr(2012, i, want).DaysInMonth(); got != want {
			t.Errorf("2012-%02d last-day daysInMonth = %d, want %d", i+1, got, want)
		}
	}
}

// TestParityStartEndOf mirrors start_end_of.js for year, quarter and month.
func TestParityStartEndOf(t *testing.T) {
	// startOf('year') from Feb 2 2011 03:04:05.006.
	base := New(time.Date(2011, time.February, 2, 3, 4, 5, 6*int(time.Millisecond), time.UTC))
	so := base.StartOf(Year)
	if so.Year() != 2011 || month0(so) != 0 || so.Date() != 1 ||
		so.Hour() != 0 || so.Minute() != 0 || so.Second() != 0 || so.Millisecond() != 0 {
		t.Errorf("startOf(year) = %s", so.Format("YYYY-MM-DD HH:mm:ss"))
	}
	eo := base.EndOf(Year)
	if eo.Year() != 2011 || month0(eo) != 11 || eo.Date() != 31 ||
		eo.Hour() != 23 || eo.Minute() != 59 || eo.Second() != 59 || eo.Millisecond() != 999 {
		t.Errorf("endOf(year) = %s.%03d", eo.Format("YYYY-MM-DD HH:mm:ss"), eo.Millisecond())
	}

	// startOf/endOf('quarter') from May 2 2011 (Q2).
	q := New(time.Date(2011, time.May, 2, 3, 4, 5, 6*int(time.Millisecond), time.UTC))
	qs := q.StartOf(Quarter)
	if qs.Year() != 2011 || qs.Quarter() != 2 || month0(qs) != 3 || qs.Date() != 1 || qs.Hour() != 0 {
		t.Errorf("startOf(quarter) = %s", qs.Format("YYYY-MM-DD HH:mm:ss"))
	}
	qe := q.EndOf(Quarter)
	if qe.Year() != 2011 || qe.Quarter() != 2 || month0(qe) != 5 || qe.Date() != 30 ||
		qe.Hour() != 23 || qe.Minute() != 59 || qe.Second() != 59 || qe.Millisecond() != 999 {
		t.Errorf("endOf(quarter) = %s.%03d", qe.Format("YYYY-MM-DD HH:mm:ss"), qe.Millisecond())
	}

	// startOf('month') from Feb 2 2011.
	ms := base.StartOf(Month)
	if ms.Date() != 1 || month0(ms) != 1 || ms.Hour() != 0 || ms.Millisecond() != 0 {
		t.Errorf("startOf(month) = %s", ms.Format("YYYY-MM-DD HH:mm:ss"))
	}
}

// TestParityAddSubtractClamp mirrors add_subtract.js: month and quarter
// arithmetic clamps the day of the month rather than overflowing.
func TestParityAddSubtractClamp(t *testing.T) {
	// Jan 31 2010 + 1 month => Feb 28 2010.
	b := fromArr(2010, 0, 31).Add(1, Month)
	if month0(b) != 1 || b.Date() != 28 {
		t.Errorf("Jan31 + 1M = %s, want Feb 28", b.Format("YYYY-MM-DD"))
	}
	// Feb 28 2010 - 1 month => Jan 28 2010.
	c := fromArr(2010, 1, 28).Subtract(1, Month)
	if month0(c) != 0 || c.Date() != 28 {
		t.Errorf("Feb28 - 1M = %s, want Jan 28", c.Format("YYYY-MM-DD"))
	}
	// Feb 28 2010 - 1 quarter => Nov 28 2009.
	d := fromArr(2010, 1, 28).Subtract(1, Quarter)
	if d.Year() != 2009 || month0(d) != 10 || d.Date() != 28 {
		t.Errorf("Feb28 2010 - 1Q = %s, want Nov 28 2009", d.Format("YYYY-MM-DD"))
	}
	// Feb 29 2012 + 1 year => Feb 28 2013 (leap-day clamp).
	y := fromArr(2012, 1, 29).Add(1, Year)
	if y.Format("YYYY-MM-DD") != "2013-02-28" {
		t.Errorf("Feb29 2012 + 1y = %s, want 2013-02-28", y.Format("YYYY-MM-DD"))
	}
}

// TestParityMonthSetterClamp mirrors getters_setters.js "month edge case":
// setting May 31 to April yields April 30, not May 1.
func TestParityMonthSetterClamp(t *testing.T) {
	a, err := ParseFormat("20130531", "YYYYMMDD")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	got := a.Set(Month, 4) // April, 1-based
	if month0(got) != 3 {
		t.Errorf("May31 set month April => month0 %d, want 3", month0(got))
	}
	if got.Date() != 30 {
		t.Errorf("May31 set month April => date %d, want 30", got.Date())
	}
}

// TestParityGetters mirrors getters_setters.js "getters".
func TestParityGetters(t *testing.T) {
	a := fromArr(2011, 9, 12, 6, 7, 8, 9)
	if a.Year() != 2011 {
		t.Errorf("year = %d", a.Year())
	}
	if month0(a) != 9 {
		t.Errorf("month = %d", month0(a))
	}
	if a.Date() != 12 {
		t.Errorf("date = %d", a.Date())
	}
	if int(a.Weekday()) != 3 {
		t.Errorf("day = %d, want 3", int(a.Weekday()))
	}
	if a.Hour() != 6 || a.Minute() != 7 || a.Second() != 8 || a.Millisecond() != 9 {
		t.Errorf("clock = %02d:%02d:%02d.%03d", a.Hour(), a.Minute(), a.Second(), a.Millisecond())
	}
	// get(unit) equivalence.
	if a.Get(Year) != 2011 || a.Get(Date) != 12 || a.Get(Hour) != 6 || a.Get(Millisecond) != 9 {
		t.Errorf("Get mismatch")
	}
}

// TestParityISOWeek mirrors weeks.js iso-week vectors for several year starts.
func TestParityISOWeek(t *testing.T) {
	cases := []struct {
		y, mo, d int
		want     int
	}{
		// year starting Sunday (2012)
		{2012, 0, 1, 52}, {2012, 0, 2, 1}, {2012, 0, 8, 1}, {2012, 0, 9, 2}, {2012, 0, 15, 2},
		// year starting Monday (2007)
		{2007, 0, 1, 1}, {2007, 0, 7, 1}, {2007, 0, 8, 2}, {2007, 0, 14, 2}, {2007, 0, 15, 3},
		// year starting Tuesday (2008)
		{2007, 11, 31, 1}, {2008, 0, 1, 1}, {2008, 0, 6, 1}, {2008, 0, 7, 2}, {2008, 0, 13, 2}, {2008, 0, 14, 3},
		// year ending/ starting around 2002-2003
		{2002, 11, 30, 1}, {2003, 0, 1, 1}, {2003, 0, 5, 1}, {2003, 0, 6, 2},
	}
	for _, c := range cases {
		if got := fromArr(c.y, c.mo, c.d).ISOWeekNumber(); got != c.want {
			t.Errorf("isoWeek(%d-%02d-%02d) = %d, want %d", c.y, c.mo+1, c.d, got, c.want)
		}
	}
}

// TestParityFormatBrackets mirrors format.js bracket-escaping vectors.
func TestParityFormatBrackets(t *testing.T) {
	b := New(time.Date(2009, time.February, 14, 15, 25, 50, 123*int(time.Millisecond), time.UTC))
	cases := []struct{ format, want string }{
		{"YY", "09"},
		{"[day]", "day"},
		{"[day] YY [YY]", "day 09 YY"},
		{"[YY", "[09"},
		{"[[YY]]", "[YY]"},
		{"[[]", "["},
		{"Q", "1"},
	}
	for _, c := range cases {
		if got := b.Format(c.format); got != c.want {
			t.Errorf("Format(%q) = %q, want %q", c.format, got, c.want)
		}
	}
}

// TestParityFormatFractionalSeconds mirrors format.js "milliseconds" runs of S.
// The instant has exactly 123 ms, so wider runs right-pad with zeros.
func TestParityFormatFractionalSeconds(t *testing.T) {
	m := New(time.Date(2009, time.February, 14, 0, 0, 0, 123*int(time.Millisecond), time.UTC))
	cases := []struct{ format, want string }{
		{"S", "1"},
		{"SS", "12"},
		{"SSS", "123"},
		{"SSSS", "1230"},
		{"SSSSS", "12300"},
		{"SSSSSS", "123000"},
		{"SSSSSSS", "1230000"},
		{"SSSSSSSS", "12300000"},
		{"SSSSSSSSS", "123000000"},
	}
	for _, c := range cases {
		if got := m.Format(c.format); got != c.want {
			t.Errorf("Format(%q) = %q, want %q", c.format, got, c.want)
		}
	}
}

// TestParityFormatHourMinuteTokens mirrors format.js hmm/hmmss/Hmm/Hmmss and
// k/kk vectors.
func TestParityFormatHourMinuteTokens(t *testing.T) {
	p := func(hms string) Moment { m, _ := ParseFormat(hms, "HH:mm:ss"); return m }
	cases := []struct {
		hms, format, want string
	}{
		{"12:34:56", "hmm", "1234"},
		{"01:34:56", "hmm", "134"},
		{"13:34:56", "hmm", "134"},
		{"12:34:56", "hmmss", "123456"},
		{"01:34:56", "hmmss", "13456"},
		{"13:34:56", "hmmss", "13456"},
		{"12:34:56", "Hmm", "1234"},
		{"01:34:56", "Hmm", "134"},
		{"13:34:56", "Hmm", "1334"},
		{"12:34:56", "Hmmss", "123456"},
		{"01:34:56", "Hmmss", "13456"},
		{"18:34:56", "Hmmss", "183456"},
		{"01:23:45", "k", "1"},
		{"12:34:56", "k", "12"},
		{"01:23:45", "kk", "01"},
		{"00:34:56", "kk", "24"},
		{"00:00:00", "kk", "24"},
	}
	for _, c := range cases {
		if got := p(c.hms).Format(c.format); got != c.want {
			t.Errorf("%s Format(%q) = %q, want %q", c.hms, c.format, got, c.want)
		}
	}
}

// TestParityFormatUnixTokens mirrors format.js X and x tokens.
func TestParityFormatUnixTokens(t *testing.T) {
	if got := Unix(1234567890).Format("X"); got != "1234567890" {
		t.Errorf("Format(X) = %q, want 1234567890", got)
	}
	if got := UnixMilli(1234567890123).Format("x"); got != "1234567890123" {
		t.Errorf("Format(x) = %q, want 1234567890123", got)
	}
}

// TestParityFormatInvalid mirrors format.js "invalid" formatting.
func TestParityFormatInvalid(t *testing.T) {
	if got := Invalid().Format("YYYY-MM-DD"); got != "Invalid date" {
		t.Errorf("Invalid().Format = %q, want %q", got, "Invalid date")
	}
}

// TestParityMinMax mirrors min_max.js ordering (min/max pick the earliest and
// latest of a set regardless of argument order).
func TestParityMinMax(t *testing.T) {
	now := New(time.Date(2017, 6, 15, 0, 0, 0, 0, time.UTC))
	future := now.Add(1, Month)
	past := now.Subtract(1, Month)
	if !Min(now, future, past).IsSame(past) {
		t.Errorf("min(now, future, past) != past")
	}
	if !Min(future, past, now).IsSame(past) {
		t.Errorf("min(future, past, now) != past")
	}
	if !Min(now).IsSame(now) {
		t.Errorf("min(now) != now")
	}
	if !Max(now, future, past).IsSame(future) {
		t.Errorf("max(now, future, past) != future")
	}
	if !Max(past, future, now).IsSame(future) {
		t.Errorf("max(past, future, now) != future")
	}
	if !Max(now).IsSame(now) {
		t.Errorf("max(now) != now")
	}
}

// TestParityFromTo mirrors from_to.js: from/to describe one Moment relative to
// another with the English default locale. The comparison is deterministic
// because it is purely between two fixed Moments.
func TestParityFromTo(t *testing.T) {
	start := New(time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC))
	cases := []struct {
		got, want string
	}{
		{start.From(start.Add(5, Second)), "a few seconds ago"},
		{start.From(start.Add(1, Minute)), "a minute ago"},
		{start.From(start.Add(5, Minute)), "5 minutes ago"},
		{start.From(start.Subtract(5, Second)), "in a few seconds"},
		{start.From(start.Subtract(1, Minute)), "in a minute"},
		{start.From(start.Subtract(5, Minute)), "in 5 minutes"},
		{start.To(start.Subtract(5, Second)), "a few seconds ago"},
		{start.To(start.Add(1, Minute)), "in a minute"},
	}
	for i, c := range cases {
		if c.got != c.want {
			t.Errorf("from/to case %d = %q, want %q", i, c.got, c.want)
		}
	}
}

// TestParityHumanizeThresholds mirrors relative_time.js default rounding
// thresholds (44 min => minutes, 45 min => an hour, and so on).
func TestParityHumanizeThresholds(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{44 * time.Minute, "44 minutes"},
		{45 * time.Minute, "an hour"},
		{21 * time.Hour, "21 hours"},
		{22 * time.Hour, "a day"},
		{25 * 24 * time.Hour, "25 days"},
		{26 * 24 * time.Hour, "a month"},
		{330 * 24 * time.Hour, "a year"},
	}
	for _, c := range cases {
		if got := Humanize(c.d); got != c.want {
			t.Errorf("Humanize(%v) = %q, want %q", c.d, got, c.want)
		}
	}
}

// TestParityDurationISOString mirrors duration.js "serialization to ISO 8601
// duration strings", including moment's mixed-sign per-group rendering.
func TestParityDurationISOString(t *testing.T) {
	cases := []struct {
		obj  map[string]int
		want string
	}{
		{map[string]int{"y": 1, "M": 2, "d": 3, "h": 4, "m": 5, "s": 6}, "P1Y2M3DT4H5M6S"},
		{map[string]int{"M": -1}, "-P1M"},
		{map[string]int{"m": -1}, "-PT1M"},
		{map[string]int{"y": -1, "M": 1}, "-P11M"},
		{map[string]int{"y": -1, "h": 1}, "-P1YT-1H"},
		{map[string]int{"y": -1, "h": 1, "m": -1}, "-P1YT-59M"},
		{map[string]int{"y": -1, "h": 1, "s": -1}, "-P1YT-59M-59S"},
		{map[string]int{"y": -1, "h": -1, "s": 1}, "-P1YT59M59S"},
		{map[string]int{"y": -1, "d": 2}, "-P1Y-2D"},
		{map[string]int{"M": 1}, "P1M"},
		{map[string]int{"y": 1, "M": 1}, "P1Y1M"},
		{map[string]int{}, "P0D"},
		{map[string]int{"M": 16, "d": 40, "s": 86465}, "P1Y4M40DT24H1M5S"},
		{map[string]int{"ms": 31952}, "PT31.952S"},
	}
	for _, c := range cases {
		if got := NewDurationFromObject(c.obj).ISOString(); got != c.want {
			t.Errorf("duration(%v).toISOString() = %q, want %q", c.obj, got, c.want)
		}
	}
}

// TestParityDurationAs mirrors duration.js asMonths/asDays/asHours vectors,
// including moment's rounding of the months-to-days conversion.
func TestParityDurationAs(t *testing.T) {
	if got := NewDuration(1, Year).AsMonths(); got != 12 {
		t.Errorf("1 year asMonths = %v, want 12", got)
	}
	if got := NewDuration(400, Year).AsMonths(); got != 4800 {
		t.Errorf("400 years asMonths = %v, want 4800", got)
	}
	for _, c := range []struct {
		years int
		days  float64
	}{{1, 365}, {2, 730}, {3, 1096}, {4, 1461}} {
		if got := NewDuration(c.years, Year).AsDays(); got != c.days {
			t.Errorf("%d years asDays = %v, want %v", c.years, got, c.days)
		}
	}
	if got := NewDuration(1, Year).AsHours(); got != 8760 {
		t.Errorf("1 year asHours = %v, want 8760", got)
	}
	if got := NewDuration(1, Quarter).AsMonths(); got != 3 {
		t.Errorf("1 quarter asMonths = %v, want 3", got)
	}
}
