package moment

// This file encodes additional known-answer vectors transcribed from the
// upstream moment.js test suite (github.com/moment/moment, src/test/moment/*.js)
// that complement the cases in parity_test.go. It reuses the fromArr and month0
// helpers defined there. Every assertion pins this Go port to a value moment.js
// itself asserts.

import (
	"testing"
	"time"
)

// TestParityLeapYear mirrors is-leap-year.js "leap years".
func TestParityLeapYear(t *testing.T) {
	cases := []struct {
		year int
		want bool
	}{
		{1, false}, {4, true}, {100, false}, {400, true}, {1700, false},
		{1900, false}, {1904, true}, {2000, true}, {2008, true}, {2010, false},
		{2024, true}, {3448, true}, {4000, true}, {0, true}, {-400, true},
	}
	for _, c := range cases {
		got := New(time.Date(c.year, time.January, 1, 0, 0, 0, 0, time.UTC)).IsLeapYear()
		if got != c.want {
			t.Errorf("IsLeapYear(%d) = %v, want %v", c.year, got, c.want)
		}
	}
}

// TestParityFormatNegativeYears mirrors format.js "handle negative years":
// YYYY/YY place the sign outside the zero padding.
func TestParityFormatNegativeYears(t *testing.T) {
	cases := []struct {
		year             int
		wantYY, wantYYYY string
	}{
		{-1, "-01", "-0001"},
		{-12, "-12", "-0012"},
		{-123, "-23", "-0123"},
		{-1234, "-34", "-1234"},
		{-12345, "-45", "-12345"},
	}
	for _, c := range cases {
		m := New(time.Date(c.year, time.January, 1, 0, 0, 0, 0, time.UTC))
		if got := m.Format("YY"); got != c.wantYY {
			t.Errorf("year %d Format(YY) = %q, want %q", c.year, got, c.wantYY)
		}
		if got := m.Format("YYYY"); got != c.wantYYYY {
			t.Errorf("year %d Format(YYYY) = %q, want %q", c.year, got, c.wantYYYY)
		}
	}
}

// TestParityDiffKeys mirrors diff.js "diff key after" and "diff key before":
// integer diffs in each unit, both signs.
func TestParityDiffKeys(t *testing.T) {
	// "diff key after" — receiver earlier than argument, negative results.
	after := []struct {
		a, b []int
		unit Unit
		want int
	}{
		{[]int{2010}, []int{2011}, Year, -1},
		{[]int{2010}, []int{2010, 2}, Month, -2},
		{[]int{2010}, []int{2010, 0, 7}, Week, 0},
		{[]int{2010}, []int{2010, 0, 8}, Week, -1},
		{[]int{2010}, []int{2010, 0, 21}, Week, -2},
		{[]int{2010}, []int{2010, 0, 22}, Week, -3},
		{[]int{2010}, []int{2010, 0, 4}, Day, -3},
		{[]int{2010}, []int{2010, 0, 1, 4}, Hour, -4},
		{[]int{2010}, []int{2010, 0, 1, 0, 5}, Minute, -5},
		{[]int{2010}, []int{2010, 0, 1, 0, 0, 6}, Second, -6},
	}
	for _, c := range after {
		if got := fromArr(c.a...).DiffInt(fromArr(c.b...), c.unit); got != c.want {
			t.Errorf("diff after %v %v %s = %d, want %d", c.a, c.b, c.unit, got, c.want)
		}
	}

	// "diff key before" — receiver later than argument, positive results.
	before := []struct {
		a, b []int
		unit Unit
		want int
	}{
		{[]int{2011}, []int{2010}, Year, 1},
		{[]int{2010, 2}, []int{2010}, Month, 2},
		{[]int{2010, 0, 4}, []int{2010}, Day, 3},
		{[]int{2010, 0, 7}, []int{2010}, Week, 0},
		{[]int{2010, 0, 8}, []int{2010}, Week, 1},
		{[]int{2010, 0, 22}, []int{2010}, Week, 3},
		{[]int{2010, 0, 1, 4}, []int{2010}, Hour, 4},
		{[]int{2010, 0, 1, 0, 5}, []int{2010}, Minute, 5},
		{[]int{2010, 0, 1, 0, 0, 6}, []int{2010}, Second, 6},
	}
	for _, c := range before {
		if got := fromArr(c.a...).DiffInt(fromArr(c.b...), c.unit); got != c.want {
			t.Errorf("diff before %v %v %s = %d, want %d", c.a, c.b, c.unit, got, c.want)
		}
	}

	// "diff month": Jan 31 2011 to Mar 1 2011 is -1 month.
	if got := fromArr(2011, 0, 31).DiffInt(fromArr(2011, 2, 1), Month); got != -1 {
		t.Errorf("Jan31 diff Mar1 months = %d, want -1", got)
	}

	// Raw millisecond diffs (diff.js "diff").
	if got := UnixMilli(1000).DiffInt(UnixMilli(0), Millisecond); got != 1000 {
		t.Errorf("1000ms diff 0 = %d, want 1000", got)
	}
	if got := UnixMilli(0).DiffInt(UnixMilli(1000), Millisecond); got != -1000 {
		t.Errorf("0 diff 1000ms = %d, want -1000", got)
	}
}

// TestParityDiffEndOfMonth mirrors diff.js "end of month diff" and "end of
// month diff with time behind", which exercise moment's clamped month anchors.
func TestParityDiffEndOfMonth(t *testing.T) {
	d := func(s string) Moment {
		tt, _ := time.Parse("2006-01-02", s)
		return New(tt)
	}
	cases := []struct {
		a, b string
		want int
	}{
		{"2016-02-29", "2016-01-30", 1},
		{"2016-02-29", "2016-01-31", 1},
		{"2017-03-31", "2017-02-28", 1},
		{"2017-02-28", "2017-03-31", -1},
	}
	for _, c := range cases {
		if got := d(c.a).DiffInt(d(c.b), Month); got != c.want {
			t.Errorf("%s diff %s months = %d, want %d", c.a, c.b, got, c.want)
		}
	}
	// (May 31 + 1 month) to May 31 is 1 month.
	if got := d("2016-05-31").Add(1, Month).DiffInt(d("2016-05-31"), Month); got != 1 {
		t.Errorf("(May31+1M) diff May31 = %d, want 1", got)
	}
}

// TestParityAddSingleUnits mirrors add_subtract.js single-unit adds: each unit
// advances independently without carrying into neighbouring fields.
func TestParityAddSingleUnits(t *testing.T) {
	start := fromArr(2011, 9, 12, 6, 7, 8, 500) // Oct 12 2011 06:07:08.500
	if got := start.Add(50, Millisecond).Millisecond(); got != 550 {
		t.Errorf("add 50ms -> %d, want 550", got)
	}
	if got := start.Add(1, Second).Second(); got != 9 {
		t.Errorf("add 1s -> %d, want 9", got)
	}
	if got := start.Add(1, Minute).Minute(); got != 8 {
		t.Errorf("add 1m -> %d, want 8", got)
	}
	if got := start.Add(1, Hour).Hour(); got != 7 {
		t.Errorf("add 1h -> %d, want 7", got)
	}
	if got := start.Add(1, Day).Date(); got != 13 {
		t.Errorf("add 1d -> %d, want 13", got)
	}
	if got := month0(start.Add(1, Month)); got != 10 {
		t.Errorf("add 1M -> month0 %d, want 10", got)
	}
	if got := start.Add(1, Year).Year(); got != 2012 {
		t.Errorf("add 1y -> %d, want 2012", got)
	}
}

// TestParityISOWeekdayGetter mirrors weekday.js "iso weekday" getter cases.
func TestParityISOWeekdayGetter(t *testing.T) {
	cases := []struct {
		y, mo, d, want int
	}{
		{1985, 1, 4, 1}, {2029, 8, 18, 2}, {2013, 3, 24, 3}, {2015, 2, 5, 4},
		{1970, 0, 2, 5}, {2001, 4, 12, 6}, {2000, 0, 2, 7},
	}
	for _, c := range cases {
		if got := fromArr(c.y, c.mo, c.d).ISOWeekday(); got != c.want {
			t.Errorf("isoWeekday(%d-%02d-%02d) = %d, want %d", c.y, c.mo+1, c.d, got, c.want)
		}
	}
}

// TestParityISOWeekYear mirrors week-year.js "iso week year" (locale-
// independent) known-answer cases.
func TestParityISOWeekYear(t *testing.T) {
	cases := []struct {
		y, mo, d, want int
	}{
		{2005, 0, 1, 2004}, {2005, 0, 2, 2004}, {2005, 0, 3, 2005},
		{2005, 11, 31, 2005}, {2006, 0, 1, 2005}, {2006, 0, 2, 2006},
		{2007, 0, 1, 2007}, {2007, 11, 30, 2007}, {2007, 11, 31, 2008},
		{2008, 0, 1, 2008}, {2008, 11, 28, 2008}, {2008, 11, 29, 2009},
		{2008, 11, 30, 2009}, {2008, 11, 31, 2009}, {2009, 0, 1, 2009},
		{2010, 0, 1, 2009}, {2010, 0, 2, 2009}, {2010, 0, 3, 2009},
		{2010, 0, 4, 2010},
	}
	for _, c := range cases {
		if got := fromArr(c.y, c.mo, c.d).ISOWeekYear(); got != c.want {
			t.Errorf("isoWeekYear(%d-%02d-%02d) = %d, want %d", c.y, c.mo+1, c.d, got, c.want)
		}
	}
}

// TestParityWeekYearLocale mirrors week-year.js "week year" under the
// dow:1,doy:4 (ISO-like) week rules the upstream test installs. The bundled
// "fr" locale carries exactly those rules.
func TestParityWeekYearLocale(t *testing.T) {
	cases := []struct {
		y, mo, d, want int
	}{
		{2005, 0, 1, 2004}, {2005, 0, 3, 2005}, {2006, 0, 1, 2005},
		{2007, 11, 31, 2008}, {2008, 11, 29, 2009}, {2010, 0, 1, 2009},
		{2010, 0, 4, 2010},
	}
	for _, c := range cases {
		if got := fromArr(c.y, c.mo, c.d).Locale("fr").WeekYear(); got != c.want {
			t.Errorf("weekYear(fr)(%d-%02d-%02d) = %d, want %d", c.y, c.mo+1, c.d, got, c.want)
		}
	}
}

// TestParityISOWeeksInYear mirrors weeks-in-year.js "isoWeeksInISOWeekYear":
// the port's ISOWeeksInYear is keyed by the ISO week-numbering year.
func TestParityISOWeeksInYear(t *testing.T) {
	cases := []struct {
		date string
		want int
	}{
		{"2003-12-29", 53}, {"2005-01-03", 52}, {"2006-01-02", 52},
		{"2007-01-01", 52}, {"2007-12-31", 52}, {"2008-12-29", 53},
	}
	for _, c := range cases {
		tt, _ := time.Parse("2006-01-02", c.date)
		if got := New(tt).ISOWeeksInYear(); got != c.want {
			t.Errorf("isoWeeksInYear(%s) = %d, want %d", c.date, got, c.want)
		}
	}
}

// TestParityFormatOrdinals mirrors format.js "quarter ordinal formats" (Qo) and
// the English day ordinals (Do).
func TestParityFormatOrdinals(t *testing.T) {
	qo := []struct {
		y, mo, d int
		want     string
	}{
		{1985, 1, 4, "1st"}, {2029, 8, 18, "3rd"}, {2013, 3, 24, "2nd"},
		{2015, 2, 5, "1st"}, {1970, 0, 2, "1st"}, {2001, 11, 12, "4th"},
	}
	for _, c := range qo {
		if got := fromArr(c.y, c.mo, c.d).Format("Qo"); got != c.want {
			t.Errorf("Qo(%d-%02d-%02d) = %q, want %q", c.y, c.mo+1, c.d, got, c.want)
		}
	}
	if got := fromArr(2000, 0, 2).Format("Qo [quarter] YYYY"); got != "1st quarter 2000" {
		t.Errorf("Qo compound = %q, want %q", got, "1st quarter 2000")
	}
	do := []struct {
		day  int
		want string
	}{
		{1, "1st"}, {2, "2nd"}, {3, "3rd"}, {4, "4th"}, {11, "11th"},
		{21, "21st"}, {22, "22nd"}, {23, "23rd"}, {31, "31st"},
	}
	for _, c := range do {
		if got := fromArr(2013, 0, c.day).Format("Do"); got != c.want {
			t.Errorf("Do(day %d) = %q, want %q", c.day, got, c.want)
		}
	}
}

// TestParityConvertRoundTrip mirrors to_type.js toObject/toArray/toJSON.
func TestParityConvertRoundTrip(t *testing.T) {
	obj := map[string]int{
		"years": 2010, "months": 3, "date": 5, "hours": 15,
		"minutes": 10, "seconds": 3, "milliseconds": 123,
	}
	got := FromObject(obj).ToObject()
	for k, v := range obj {
		if got[k] != v {
			t.Errorf("toObject[%q] = %d, want %d", k, got[k], v)
		}
	}

	arr := []int{2014, 11, 26, 11, 46, 58, 17}
	gotArr := FromArray(arr).ToArray()
	if len(gotArr) != len(arr) {
		t.Fatalf("toArray len = %d, want %d", len(gotArr), len(arr))
	}
	for i := range arr {
		if gotArr[i] != arr[i] {
			t.Errorf("toArray[%d] = %d, want %d", i, gotArr[i], arr[i])
		}
	}

	if got := New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)).ToJSON(); got != "2000-01-01T00:00:00.000Z" {
		t.Errorf("toJSON = %q, want %q", got, "2000-01-01T00:00:00.000Z")
	}
}

// TestParityHumanizeSuffix mirrors relative_time.js "default thresholds
// fromNow": each boundary rounds to the expected past phrase. The reference is a
// fixed Moment so the assertion is deterministic.
func TestParityHumanizeSuffix(t *testing.T) {
	base := New(time.Date(2017, 6, 15, 12, 0, 0, 0, time.UTC))
	cases := []struct {
		n    int
		unit Unit
		want string
	}{
		{44, Second, "a few seconds ago"},
		{45, Second, "a minute ago"},
		{44, Minute, "44 minutes ago"},
		{45, Minute, "an hour ago"},
		{21, Hour, "21 hours ago"},
		{22, Hour, "a day ago"},
		{25, Day, "25 days ago"},
		{26, Day, "a month ago"},
		{10, Month, "10 months ago"},
		{11, Month, "a year ago"},
	}
	for _, c := range cases {
		if got := base.Subtract(c.n, c.unit).From(base); got != c.want {
			t.Errorf("%d %s ago = %q, want %q", c.n, c.unit, got, c.want)
		}
	}
	// Future direction (relative_time.js implies the symmetric "in" phrasing).
	if got := base.Add(45, Second).From(base); got != "in a minute" {
		t.Errorf("in a minute = %q", got)
	}
	if got := base.Add(5, Minute).From(base); got != "in 5 minutes" {
		t.Errorf("in 5 minutes = %q", got)
	}
}

// TestParityStartEndOfWeek mirrors start_end_of.js week and iso-week cases from
// Feb 2 2011 (a Wednesday) with the default Sunday-first locale.
func TestParityStartEndOfWeek(t *testing.T) {
	base := New(time.Date(2011, time.February, 2, 3, 4, 5, 6*int(time.Millisecond), time.UTC))

	sw := base.StartOf(Week)
	if month0(sw) != 0 || int(sw.Weekday()) != 0 || sw.Date() != 30 || sw.Hour() != 0 || sw.Millisecond() != 0 {
		t.Errorf("startOf(week) = %s (weekday %d)", sw.Format("YYYY-MM-DD HH:mm:ss"), int(sw.Weekday()))
	}
	ew := base.EndOf(Week)
	if month0(ew) != 1 || int(ew.Weekday()) != 6 || ew.Date() != 5 ||
		ew.Hour() != 23 || ew.Minute() != 59 || ew.Second() != 59 || ew.Millisecond() != 999 {
		t.Errorf("endOf(week) = %s (weekday %d)", ew.Format("YYYY-MM-DD HH:mm:ss"), int(ew.Weekday()))
	}

	siw := base.StartOf(ISOWeek)
	if month0(siw) != 0 || siw.ISOWeekday() != 1 || siw.Date() != 31 || siw.Hour() != 0 {
		t.Errorf("startOf(isoWeek) = %s (isoWeekday %d)", siw.Format("YYYY-MM-DD HH:mm:ss"), siw.ISOWeekday())
	}
	eiw := base.EndOf(ISOWeek)
	if month0(eiw) != 1 || eiw.ISOWeekday() != 7 || eiw.Date() != 6 ||
		eiw.Hour() != 23 || eiw.Minute() != 59 || eiw.Second() != 59 || eiw.Millisecond() != 999 {
		t.Errorf("endOf(isoWeek) = %s (isoWeekday %d)", eiw.Format("YYYY-MM-DD HH:mm:ss"), eiw.ISOWeekday())
	}
}
