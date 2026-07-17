package moment

import (
	"testing"
	"time"
)

// ref is the fixed reference instant used across the deterministic tests:
// Friday, 14 July 2017 02:40:00 UTC.
var ref = time.Date(2017, time.July, 14, 2, 40, 0, 0, time.UTC)

func refMoment() Moment { return New(ref) }

func TestConstructors(t *testing.T) {
	if got := New(ref).Year(); got != 2017 {
		t.Fatalf("New year = %d", got)
	}
	if got := FromTime(ref).Month(); got != time.July {
		t.Fatalf("FromTime month = %v", got)
	}
	if got := Unix(1500000000).Unix(); got != 1500000000 {
		t.Fatalf("Unix = %d", got)
	}
	if got := UnixMilli(1500000000123).UnixMilli(); got != 1500000000123 {
		t.Fatalf("UnixMilli = %d", got)
	}
	dt := DateTime(2017, time.July, 14, 2, 40, 0, 0, nil)
	if !dt.IsSame(refMoment()) {
		t.Fatalf("DateTime = %s", dt)
	}
	if DateTime(2017, time.July, 14, 0, 0, 0, 0, time.UTC).Location() != time.UTC {
		t.Fatalf("DateTime location")
	}
}

func TestNowInjectedClock(t *testing.T) {
	m := Now()
	if m.IsZero() {
		t.Fatalf("Now returned zero")
	}
	c := FixedClock(ref)
	nw := NowWith(c)
	if !nw.IsSame(refMoment()) {
		t.Fatalf("NowWith = %s", nw)
	}
	if !NowWith(nil).IsValid() {
		t.Fatalf("NowWith(nil) invalid")
	}
}

func TestZeroValidClone(t *testing.T) {
	var z Moment
	if z.IsValid() || !z.IsZero() {
		t.Fatalf("zero moment validity")
	}
	m := refMoment()
	if !m.IsValid() {
		t.Fatalf("ref invalid")
	}
	clone := m.Clone()
	if !clone.IsSame(m) {
		t.Fatalf("clone mismatch")
	}
}

func TestFormatTokens(t *testing.T) {
	m := New(time.Date(2017, time.July, 4, 14, 5, 9, 0, time.UTC))
	cases := map[string]string{
		"YYYY":                     "2017",
		"YY":                       "17",
		"MMMM":                     "July",
		"MMM":                      "Jul",
		"MM":                       "07",
		"M":                        "7",
		"DD":                       "04",
		"D":                        "4",
		"dddd":                     "Tuesday",
		"ddd":                      "Tue",
		"HH":                       "14",
		"hh":                       "02",
		"h":                        "2",
		"mm":                       "05",
		"m":                        "5",
		"ss":                       "09",
		"s":                        "9",
		"A":                        "PM",
		"a":                        "pm",
		"YYYY-MM-DD":               "2017-07-04",
		"YYYY-MM-DDTHH:mm:ss":      "2017-07-04T14:05:09",
		"dddd, MMMM D, YYYY":       "Tuesday, July 4, 2017",
		"h:mm A":                   "2:05 PM",
		"YYYY [year] MM [month]":   "2017 year 07 month",
		"[literal only no tokens]": "literal only no tokens",
	}
	for format, want := range cases {
		if got := m.Format(format); got != want {
			t.Errorf("Format(%q) = %q, want %q", format, got, want)
		}
	}
}

func TestFormatLayoutAndISO(t *testing.T) {
	m := refMoment()
	if got := m.FormatLayout("2006-01-02"); got != "2017-07-14" {
		t.Fatalf("FormatLayout = %q", got)
	}
	if got := m.ISO(); got != "2017-07-14T02:40:00Z" {
		t.Fatalf("ISO = %q", got)
	}
	if got := m.String(); got != m.ISO() {
		t.Fatalf("String != ISO")
	}
}

func TestParseRoundTrip(t *testing.T) {
	const format = "YYYY-MM-DD HH:mm:ss"
	const value = "2017-07-14 02:40:05"
	m, err := ParseFormat(value, format)
	if err != nil {
		t.Fatalf("ParseFormat: %v", err)
	}
	if got := m.Format(format); got != value {
		t.Fatalf("round trip = %q", got)
	}
	if _, err := ParseFormat("nonsense", format); err == nil {
		t.Fatalf("expected parse error")
	}
}

func TestParseAuto(t *testing.T) {
	for _, v := range []string{
		"2017-07-14T02:40:00Z",
		"2017-07-14 02:40:00",
		"2017-07-14",
		"2017/07/14",
	} {
		if _, err := Parse(v); err != nil {
			t.Errorf("Parse(%q): %v", v, err)
		}
	}
	if _, err := Parse("definitely not a date"); err == nil {
		t.Fatalf("expected error for junk")
	}
}

func TestParseLayoutAndInLocation(t *testing.T) {
	m, err := ParseLayout("2017-07-14", "2006-01-02")
	if err != nil || m.Year() != 2017 {
		t.Fatalf("ParseLayout: %v %s", err, m)
	}
	if _, err := ParseLayout("x", "2006"); err == nil {
		t.Fatalf("expected layout error")
	}
	nyc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	m2, err := ParseInLocation("2017-07-14 02:40:00", "YYYY-MM-DD HH:mm:ss", nyc)
	if err != nil {
		t.Fatalf("ParseInLocation: %v", err)
	}
	if m2.Hour() != 2 {
		t.Fatalf("ParseInLocation hour = %d", m2.Hour())
	}
	if _, err := ParseInLocation("bad", "YYYY", nil); err == nil {
		t.Fatalf("expected error")
	}
}

func TestSupportedTokens(t *testing.T) {
	toks := SupportedTokens()
	found := false
	for _, tk := range toks {
		if tk == "YYYY" {
			found = true
			break
		}
	}
	if len(toks) == 0 || !found {
		t.Fatalf("SupportedTokens = %v", toks)
	}
}

func TestAddSubtract(t *testing.T) {
	m := refMoment()
	cases := []struct {
		n    int
		unit Unit
		want string
	}{
		{1, Year, "2018-07-14T02:40:00Z"},
		{2, Month, "2017-09-14T02:40:00Z"},
		{1, Week, "2017-07-21T02:40:00Z"},
		{3, Day, "2017-07-17T02:40:00Z"},
		{5, Hour, "2017-07-14T07:40:00Z"},
		{20, Minute, "2017-07-14T03:00:00Z"},
		{30, Second, "2017-07-14T02:40:30Z"},
		{1, "days", "2017-07-15T02:40:00Z"},
	}
	for _, c := range cases {
		if got := m.Add(c.n, c.unit).ISO(); got != c.want {
			t.Errorf("Add(%d,%s) = %s, want %s", c.n, c.unit, got, c.want)
		}
	}
	if got := m.Subtract(1, Day).ISO(); got != "2017-07-13T02:40:00Z" {
		t.Fatalf("Subtract = %s", got)
	}
	if got := m.AddDuration(90 * time.Minute).ISO(); got != "2017-07-14T04:10:00Z" {
		t.Fatalf("AddDuration = %s", got)
	}
	// Immutability.
	if m.ISO() != "2017-07-14T02:40:00Z" {
		t.Fatalf("receiver mutated: %s", m)
	}
	if got := m.Add(500, Millisecond).UnixMilli(); got != m.UnixMilli()+500 {
		t.Fatalf("Add millis")
	}
}

func TestStartEndOf(t *testing.T) {
	m := refMoment()
	cases := []struct {
		unit  Unit
		start string
		end   string
	}{
		{Year, "2017-01-01T00:00:00Z", "2017-12-31T23:59:59.999999999Z"},
		{Month, "2017-07-01T00:00:00Z", "2017-07-31T23:59:59.999999999Z"},
		{Day, "2017-07-14T00:00:00Z", "2017-07-14T23:59:59.999999999Z"},
		{Hour, "2017-07-14T02:00:00Z", "2017-07-14T02:59:59.999999999Z"},
		{Minute, "2017-07-14T02:40:00Z", "2017-07-14T02:40:59.999999999Z"},
	}
	for _, c := range cases {
		if got := m.StartOf(c.unit).ISO(); got != c.start {
			t.Errorf("StartOf(%s) = %s, want %s", c.unit, got, c.start)
		}
		if got := m.EndOf(c.unit).ISO(); got != c.end {
			t.Errorf("EndOf(%s) = %s, want %s", c.unit, got, c.end)
		}
	}
	// Week starts Sunday: 2017-07-14 is Friday -> Sunday 2017-07-09.
	if got := m.StartOf(Week).ISO(); got != "2017-07-09T00:00:00Z" {
		t.Fatalf("StartOf week = %s", got)
	}
	if got := m.EndOf(Week).ISO(); got != "2017-07-15T23:59:59.999999999Z" {
		t.Fatalf("EndOf week = %s", got)
	}
	if got := m.StartOf(Second).Nanosecond(); got != 0 {
		t.Fatalf("StartOf second ns = %d", got)
	}
	mm := New(time.Date(2017, 7, 14, 2, 40, 5, 123456789, time.UTC))
	if got := mm.StartOf(Millisecond).Nanosecond(); got != 123000000 {
		t.Fatalf("StartOf millisecond = %d", got)
	}
	if got := mm.EndOf(Millisecond).Nanosecond(); got != 123999999 {
		t.Fatalf("EndOf millisecond = %d", got)
	}
	if got := m.EndOf("unknown"); !got.IsSame(m) {
		t.Fatalf("EndOf unknown should be no-op")
	}
}

func TestSet(t *testing.T) {
	m := refMoment()
	cases := []struct {
		unit  Unit
		value int
		want  string
	}{
		{Year, 2020, "2020-07-14T02:40:00Z"},
		{Month, 12, "2017-12-14T02:40:00Z"},
		{Date, 1, "2017-07-01T02:40:00Z"},
		{Hour, 9, "2017-07-14T09:40:00Z"},
		{Minute, 15, "2017-07-14T02:15:00Z"},
		{Second, 30, "2017-07-14T02:40:30Z"},
	}
	for _, c := range cases {
		if got := m.Set(c.unit, c.value).ISO(); got != c.want {
			t.Errorf("Set(%s,%d) = %s, want %s", c.unit, c.value, got, c.want)
		}
	}
	if got := m.Set(Millisecond, 500).Millisecond(); got != 500 {
		t.Fatalf("Set millisecond = %d", got)
	}
}

func TestGetters(t *testing.T) {
	m := New(time.Date(2017, time.July, 14, 2, 40, 5, 123000000, time.UTC))
	if m.Year() != 2017 || m.Month() != time.July || m.Date() != 14 || m.Day() != 14 {
		t.Fatalf("date getters")
	}
	if m.Hour() != 2 || m.Minute() != 40 || m.Second() != 5 || m.Millisecond() != 123 {
		t.Fatalf("clock getters")
	}
	if m.Nanosecond() != 123000000 {
		t.Fatalf("nanosecond")
	}
	if m.Weekday() != time.Friday {
		t.Fatalf("weekday = %v", m.Weekday())
	}
	if m.DayOfYear() != 195 {
		t.Fatalf("day of year = %d", m.DayOfYear())
	}
	if y, w := m.ISOWeek(); y != 2017 || w != 28 {
		t.Fatalf("ISOWeek = %d %d", y, w)
	}
	if !New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)).IsLeapYear() {
		t.Fatalf("2020 should be leap")
	}
	if m.IsLeapYear() {
		t.Fatalf("2017 not leap")
	}
}

func TestComparisons(t *testing.T) {
	a := refMoment()
	b := a.Add(1, Day)
	if !a.IsBefore(b) || !b.IsAfter(a) {
		t.Fatalf("before/after")
	}
	if !a.IsSame(a.Clone()) {
		t.Fatalf("same")
	}
	if !a.IsSameOrBefore(a) || !a.IsSameOrBefore(b) {
		t.Fatalf("same-or-before")
	}
	if !a.IsSameOrAfter(a) || !b.IsSameOrAfter(a) {
		t.Fatalf("same-or-after")
	}
	mid := a.Add(12, Hour)
	if !mid.IsBetween(a, b, false) {
		t.Fatalf("between exclusive")
	}
	if mid.IsBetween(b, a, false) != true {
		t.Fatalf("between reversed bounds")
	}
	if a.IsBetween(a, b, false) {
		t.Fatalf("exclusive bound should be false")
	}
	if !a.IsBetween(a, b, true) {
		t.Fatalf("inclusive bound should be true")
	}
	if !a.IsSameUnit(a.Add(3, Hour), Day) {
		t.Fatalf("same day unit")
	}
	if a.IsSameUnit(b, Day) {
		t.Fatalf("different day unit")
	}
}

func TestDiff(t *testing.T) {
	a := New(time.Date(2017, 7, 14, 0, 0, 0, 0, time.UTC))
	b := New(time.Date(2017, 7, 11, 0, 0, 0, 0, time.UTC))
	if got := a.Diff(b, Day); got != 3 {
		t.Fatalf("diff days = %v", got)
	}
	if got := a.DiffInt(b, Hour); got != 72 {
		t.Fatalf("diff hours = %d", got)
	}
	if got := a.DiffInt(b, Week); got != 0 {
		t.Fatalf("diff weeks = %d", got)
	}
	if got := b.DiffInt(a, Day); got != -3 {
		t.Fatalf("diff negative = %d", got)
	}
	mar := New(time.Date(2017, 3, 14, 0, 0, 0, 0, time.UTC))
	jan := New(time.Date(2017, 1, 14, 0, 0, 0, 0, time.UTC))
	if got := mar.DiffInt(jan, Month); got != 2 {
		t.Fatalf("diff months = %d", got)
	}
	y2 := New(time.Date(2019, 7, 14, 0, 0, 0, 0, time.UTC))
	if got := y2.DiffInt(a, Year); got != 2 {
		t.Fatalf("diff years = %d", got)
	}
	// Fractional month diff should be between 1 and 2.
	partial := New(time.Date(2017, 2, 28, 0, 0, 0, 0, time.UTC))
	if got := partial.Diff(jan, Month); got <= 1 || got >= 2 {
		t.Fatalf("fractional month diff = %v", got)
	}
	if got := a.DiffInt(b, Minute); got != 3*24*60 {
		t.Fatalf("diff minutes = %d", got)
	}
	if got := a.DiffInt(b, Second); got != 3*24*3600 {
		t.Fatalf("diff seconds = %d", got)
	}
	if got := a.Diff(b, Millisecond); got != 3*24*3600*1000 {
		t.Fatalf("diff millis = %v", got)
	}
	if got := a.Diff(b, "bogus"); got != 0 {
		t.Fatalf("diff bogus = %v", got)
	}
	if got := a.DiffDuration(b); got != 72*time.Hour {
		t.Fatalf("diff duration = %v", got)
	}
}

func TestRelativeFromNow(t *testing.T) {
	clock := FixedClock(ref)
	base := New(ref).WithClock(clock)
	if got := base.Add(3, Day).FromNow(); got != "in 3 days" {
		t.Fatalf("FromNow future = %q", got)
	}
	if got := base.Subtract(2, Hour).FromNow(); got != "2 hours ago" {
		t.Fatalf("FromNow past = %q", got)
	}
	if got := base.Add(1, Minute).FromNow(); got != "in a minute" {
		t.Fatalf("FromNow minute = %q", got)
	}
	if got := base.ToNow(); got != "in a few seconds" {
		t.Fatalf("ToNow = %q", got)
	}
	other := New(ref)
	if got := base.Add(5, Day).From(other); got != "in 5 days" {
		t.Fatalf("From = %q", got)
	}
	if got := base.To(other.Add(5, Day)); got != "in 5 days" {
		t.Fatalf("To = %q", got)
	}
}

func TestHumanizeMagnitudes(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{10 * time.Second, "a few seconds"},
		{60 * time.Second, "a minute"},
		{5 * time.Minute, "5 minutes"},
		{60 * time.Minute, "an hour"},
		{5 * time.Hour, "5 hours"},
		{24 * time.Hour, "a day"},
		{3 * 24 * time.Hour, "3 days"},
		{35 * 24 * time.Hour, "a month"},
		{90 * 24 * time.Hour, "3 months"},
		{400 * 24 * time.Hour, "a year"},
		{800 * 24 * time.Hour, "2 years"},
	}
	for _, c := range cases {
		if got := Humanize(c.d); got != c.want {
			t.Errorf("Humanize(%v) = %q, want %q", c.d, got, c.want)
		}
	}
}

func TestCalendar(t *testing.T) {
	clock := FixedClock(ref)
	base := New(ref).WithClock(clock)
	now := New(ref)
	cases := []struct {
		m    Moment
		want string
	}{
		{base, "Today at 2:40 AM"},
		{base.Add(1, Day), "Tomorrow at 2:40 AM"},
		{base.Subtract(1, Day), "Yesterday at 2:40 AM"},
		{base.Add(2, Day), "Sunday at 2:40 AM"},
		{base.Subtract(2, Day), "Last Wednesday at 2:40 AM"},
		{base.Add(10, Day), "07/24/2017"},
	}
	for _, c := range cases {
		if got := c.m.Calendar(now); got != c.want {
			t.Errorf("Calendar = %q, want %q", got, c.want)
		}
	}
	if got := base.CalendarNow(); got != "Today at 2:40 AM" {
		t.Fatalf("CalendarNow = %q", got)
	}
}

func TestTimezones(t *testing.T) {
	m := refMoment()
	if m.UTC().Location() != time.UTC {
		t.Fatalf("UTC location")
	}
	if m.Local().Location() != time.Local {
		t.Fatalf("Local location")
	}
	if m.In(nil).Location() != time.UTC {
		t.Fatalf("In(nil) should be UTC")
	}
	nyc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	shifted := m.In(nyc)
	if shifted.Unix() != m.Unix() {
		t.Fatalf("In changed instant")
	}
	if shifted.Hour() == m.Hour() {
		t.Fatalf("In did not shift wall clock")
	}
}

func TestWithClockNilResets(t *testing.T) {
	m := refMoment().WithClock(nil)
	if m.clockOf() == nil {
		t.Fatalf("clock should default")
	}
}
