package moment

import (
	"testing"
	"time"
)

// fixed is the deterministic reference instant used by the feature tests:
// Tuesday, 4 July 2017 14:05:09.123 UTC.
var fixed = time.Date(2017, time.July, 4, 14, 5, 9, 123000000, time.UTC)

func fixedMoment() Moment { return New(fixed) }

func TestExtendedFormatTokens(t *testing.T) {
	m := fixedMoment()
	cases := map[string]string{
		"Q":      "3",
		"Qo":     "3rd",
		"Do":     "4th",
		"DDD":    "185",
		"DDDo":   "185th",
		"DDDD":   "185",
		"w":      "27",
		"wo":     "27th",
		"ww":     "27",
		"W":      "27",
		"Wo":     "27th",
		"WW":     "27",
		"E":      "2",
		"e":      "2",
		"k":      "14",
		"kk":     "14",
		"gggg":   "2017",
		"GGGG":   "2017",
		"gg":     "17",
		"GG":     "17",
		"SSS":    "123",
		"S":      "1",
		"SSSSSS": "123000",
		"X":      "1499177109",
		"x":      "1499177109123",
		"Mo":     "7th",
		"do":     "2nd",
	}
	for format, want := range cases {
		if got := m.Format(format); got != want {
			t.Errorf("Format(%q) = %q, want %q", format, got, want)
		}
	}
}

func TestFormatMidnightHourTokens(t *testing.T) {
	m := New(time.Date(2017, 7, 4, 0, 30, 0, 0, time.UTC))
	if got := m.Format("h hh k kk H HH"); got != "12 12 24 24 0 00" {
		t.Fatalf("hour tokens = %q", got)
	}
}

func TestFormatOffsetTokens(t *testing.T) {
	loc := time.FixedZone("X", 5*3600+30*60)
	m := New(time.Date(2017, 7, 4, 14, 0, 0, 0, loc))
	if got := m.Format("Z ZZ z"); got != "+05:30 +0530 X" {
		t.Fatalf("offset tokens = %q", got)
	}
}

func TestLocaleFormats(t *testing.T) {
	m := fixedMoment()
	cases := []struct {
		locale, format, want string
	}{
		{"en", "LLLL", "Tuesday, July 4, 2017 2:05 PM"},
		{"en-gb", "LLLL", "Tuesday, 4 July 2017 14:05"},
		{"fr", "LLLL", "mardi 4 juillet 2017 14:05"},
		{"fr", "dddd", "mardi"},
		{"de", "LL", "4. Juli 2017"},
		{"es", "LLL", "4 de julio de 2017 14:05"},
		{"it", "MMMM", "luglio"},
		{"pt-br", "dddd", "terça-feira"},
		{"nl", "MMM", "jul."},
		{"ru", "D MMMM YYYY", "4 июля 2017"},
		{"zh-cn", "LLLL", "2017年7月4日星期二 14:05"},
		{"zh-tw", "dddd", "星期二"},
		{"ja", "LL", "2017年7月4日"},
		{"ko", "dddd", "화요일"},
		{"tr", "MMMM", "Temmuz"},
		{"pl", "MMMM", "lipiec"},
		{"sv", "MMMM", "juli"},
		{"cs", "MMMM", "červenec"},
	}
	for _, c := range cases {
		if got := m.Locale(c.locale).Format(c.format); got != c.want {
			t.Errorf("[%s] Format(%q) = %q, want %q", c.locale, c.format, got, c.want)
		}
	}
}

func TestLocaleMeridiem(t *testing.T) {
	morning := New(time.Date(2017, 7, 4, 5, 0, 0, 0, time.UTC))
	if got := morning.Locale("zh-cn").Format("A"); got != "凌晨" {
		t.Fatalf("zh-cn meridiem = %q", got)
	}
	if got := New(time.Date(2017, 7, 4, 15, 0, 0, 0, time.UTC)).Locale("ja").Format("A"); got != "午後" {
		t.Fatalf("ja meridiem = %q", got)
	}
	if got := fixedMoment().Locale("en").Format("a"); got != "pm" {
		t.Fatalf("en lower meridiem = %q", got)
	}
}

func TestLocaleRelativeTime(t *testing.T) {
	cases := []struct {
		locale string
		dur    Duration
		want   string
	}{
		{"en", NewDuration(3, Day), "in 3 days"},
		{"fr", NewDuration(3, Day), "dans 3 jours"},
		{"de", NewDuration(2, Hour), "in 2 Stunden"},
		{"es", NewDuration(5, Minute), "en 5 minutos"},
		{"ru", NewDuration(5, Day), "через 5 дней"},
		{"ru", NewDuration(2, Minute), "через 2 минуты"},
		{"ru", NewDuration(21, Day), "через 21 день"},
		{"pl", NewDuration(5, Day), "za 5 dni"},
		{"cs", NewDuration(3, Hour), "za 3 hodiny"},
		{"tr", NewDuration(1, Year), "bir yıl sonra"},
		{"ja", NewDuration(4, Month), "4ヶ月後"},
	}
	for _, c := range cases {
		if got := c.dur.Locale(c.locale).Humanize(true); got != c.want {
			t.Errorf("[%s] Humanize = %q, want %q", c.locale, got, c.want)
		}
	}
}

func TestRelativeTimePast(t *testing.T) {
	if got := NewDuration(-2, Hour).Humanize(true); got != "2 hours ago" {
		t.Fatalf("past = %q", got)
	}
	if got := NewDuration(-3, Day).Locale("fr").Humanize(true); got != "il y a 3 jours" {
		t.Fatalf("fr past = %q", got)
	}
}

func TestFromNowLocale(t *testing.T) {
	clock := FixedClock(fixed)
	base := New(fixed).WithClock(clock).Locale("de")
	if got := base.Add(3, Day).FromNow(); got != "in 3 Tagen" {
		t.Fatalf("de FromNow = %q", got)
	}
	if got := base.Subtract(1, Hour).FromNow(); got != "vor einer Stunde" {
		t.Fatalf("de FromNow past = %q", got)
	}
}

func TestCalendarLocale(t *testing.T) {
	clock := FixedClock(fixed)
	now := New(fixed)
	base := New(fixed).WithClock(clock).Locale("fr")
	if got := base.Calendar(now); got != "Aujourd’hui à 14:05" {
		t.Fatalf("fr calendar today = %q", got)
	}
	if got := base.Add(1, Day).Calendar(now); got != "Demain à 14:05" {
		t.Fatalf("fr calendar tomorrow = %q", got)
	}
}

func TestDurationAsAndGet(t *testing.T) {
	d := NewDuration(1, Year).Add(NewDuration(2, Month)).Add(NewDuration(10, Day))
	if d.Years() != 1 || d.Months() != 2 || d.Days() != 10 {
		t.Fatalf("components = %d %d %d", d.Years(), d.Months(), d.Days())
	}
	d2 := NewDuration(2, Day).Add(NewDuration(3, Hour)).Add(NewDuration(30, Minute))
	if got := d2.AsHours(); got != 51.5 {
		t.Fatalf("AsHours = %v", got)
	}
	if got := d2.AsMinutes(); got != 3090 {
		t.Fatalf("AsMinutes = %v", got)
	}
	if got := d2.Get(Hour); got != 3 {
		t.Fatalf("Get hour = %d", got)
	}
	if got := d2.Get(Minute); got != 30 {
		t.Fatalf("Get minute = %d", got)
	}
	week := NewDuration(10, Day)
	if week.Weeks() != 1 {
		t.Fatalf("weeks = %d", week.Weeks())
	}
}

func TestDurationISORoundTrip(t *testing.T) {
	cases := []string{"P1Y2M10DT2H30M", "P0D", "PT1H", "P3W", "PT0.5S", "-P1Y"}
	wants := []string{"P1Y2M10DT2H30M", "P0D", "PT1H", "P21D", "PT0.5S", "-P1Y"}
	for i, in := range cases {
		d, err := ParseDuration(in)
		if err != nil {
			t.Fatalf("ParseDuration(%q): %v", in, err)
		}
		if got := d.ISOString(); got != wants[i] {
			t.Errorf("ISOString(%q) = %q, want %q", in, got, wants[i])
		}
	}
	if _, err := ParseDuration("not a duration"); err == nil {
		t.Fatalf("expected error for junk")
	}
	if _, err := ParseDuration("P"); err == nil {
		t.Fatalf("expected error for bare P")
	}
}

func TestDurationBetweenAndArithmetic(t *testing.T) {
	a := New(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC))
	b := New(time.Date(2017, 1, 3, 12, 0, 0, 0, time.UTC))
	d := DurationBetween(a, b)
	if got := d.AsDays(); got != 2.5 {
		t.Fatalf("between days = %v", got)
	}
	if got := d.Humanize(false); got != "3 days" {
		t.Fatalf("between humanize = %q", got)
	}
	sum := NewDuration(1, Hour).Add(NewDuration(30, Minute))
	if got := sum.AsMinutes(); got != 90 {
		t.Fatalf("sum minutes = %v", got)
	}
	diff := NewDuration(2, Hour).Subtract(NewDuration(30, Minute))
	if got := diff.AsMinutes(); got != 90 {
		t.Fatalf("diff minutes = %v", got)
	}
	if got := NewDuration(-5, Hour).Abs().AsHours(); got != 5 {
		t.Fatalf("abs = %v", got)
	}
	if NewDuration(3, Day).Clone().AsDays() != 3 {
		t.Fatalf("clone")
	}
	if got := NewDuration(90, Minute).ToDuration(); got != 90*time.Minute {
		t.Fatalf("ToDuration = %v", got)
	}
}

func TestHumanizeThresholdConfig(t *testing.T) {
	if !SetRelativeTimeThreshold("ss", 3) {
		t.Fatalf("set threshold failed")
	}
	defer SetRelativeTimeThreshold("ss", 44)
	if RelativeTimeThreshold("ss") != 3 {
		t.Fatalf("threshold get")
	}
	if got := NewDuration(10, Second).Humanize(false); got != "10 seconds" {
		t.Fatalf("threshold humanize = %q", got)
	}
	if SetRelativeTimeThreshold("bogus", 1) {
		t.Fatalf("bogus unit accepted")
	}
	if RelativeTimeThreshold("bogus") != -1 {
		t.Fatalf("bogus threshold")
	}
}

func TestQuarterAndISOWeek(t *testing.T) {
	if q := fixedMoment().Quarter(); q != 3 {
		t.Fatalf("quarter = %d", q)
	}
	if q := New(time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC)).Quarter(); q != 1 {
		t.Fatalf("q1 = %d", q)
	}
	m := fixedMoment()
	if m.ISOWeekYear() != 2017 || m.ISOWeekNumber() != 27 {
		t.Fatalf("iso week = %d/%d", m.ISOWeekYear(), m.ISOWeekNumber())
	}
	if m.ISOWeekday() != 2 {
		t.Fatalf("iso weekday = %d", m.ISOWeekday())
	}
	// 1 Jan 2017 is a Sunday: ISO week 52 of 2016.
	jan1 := New(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC))
	if jan1.ISOWeekYear() != 2016 || jan1.ISOWeekNumber() != 52 {
		t.Fatalf("jan1 iso = %d/%d", jan1.ISOWeekYear(), jan1.ISOWeekNumber())
	}
	if jan1.ISOWeekday() != 7 {
		t.Fatalf("jan1 iso weekday = %d", jan1.ISOWeekday())
	}
}

func TestStartOfQuarterAndISOWeek(t *testing.T) {
	m := fixedMoment()
	if got := m.StartOf(Quarter).ISO(); got != "2017-07-01T00:00:00Z" {
		t.Fatalf("start of quarter = %q", got)
	}
	if got := m.EndOf(Quarter).Format("YYYY-MM-DD"); got != "2017-09-30" {
		t.Fatalf("end of quarter = %q", got)
	}
	// Tuesday 4 July -> ISO week starts Monday 3 July.
	if got := m.StartOf(ISOWeek).ISO(); got != "2017-07-03T00:00:00Z" {
		t.Fatalf("start of iso week = %q", got)
	}
	// Locale week for en (Sunday start) -> Sunday 2 July.
	if got := m.StartOf(Week).ISO(); got != "2017-07-02T00:00:00Z" {
		t.Fatalf("start of week en = %q", got)
	}
	// fr locale (Monday start) -> Monday 3 July.
	if got := m.Locale("fr").StartOf(Week).ISO(); got != "2017-07-03T00:00:00Z" {
		t.Fatalf("start of week fr = %q", got)
	}
}

func TestDaysInMonthAndSetters(t *testing.T) {
	if got := New(time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC)).DaysInMonth(); got != 29 {
		t.Fatalf("feb 2016 = %d", got)
	}
	if got := fixedMoment().DaysInMonth(); got != 31 {
		t.Fatalf("july = %d", got)
	}
	if got := fixedMoment().Set(Quarter, 1).Format("YYYY-MM-DD"); got != "2017-01-04" {
		t.Fatalf("set quarter = %q", got)
	}
	if got := fixedMoment().Set(DayOfYear, 1).Format("YYYY-MM-DD"); got != "2017-01-01" {
		t.Fatalf("set dayofyear = %q", got)
	}
	y, mo, d := 2020, 3, 15
	got := fixedMoment().SetAll(DateSpec{Year: &y, Month: &mo, Day: &d})
	if got.Format("YYYY-MM-DD HH:mm:ss") != "2020-03-15 14:05:09" {
		t.Fatalf("SetAll = %q", got.Format("YYYY-MM-DD HH:mm:ss"))
	}
}

func TestConstructFromArrayObject(t *testing.T) {
	// moment array: month is 0-based; 6 == July.
	m := FromArray([]int{2017, 6, 4, 14, 5, 9})
	if m.Format("YYYY-MM-DD HH:mm:ss") != "2017-07-04 14:05:09" {
		t.Fatalf("FromArray = %q", m.ISO())
	}
	if !FromArray(nil).IsValid() == false {
		// nil yields invalid
		if FromArray(nil).IsValid() {
			t.Fatalf("nil array should be invalid")
		}
	}
	o := FromObject(map[string]int{"year": 2017, "month": 6, "day": 4, "hour": 9})
	if o.Format("YYYY-MM-DD HH:mm") != "2017-07-04 09:00" {
		t.Fatalf("FromObject = %q", o.ISO())
	}
}

func TestMaxMin(t *testing.T) {
	a := New(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC))
	b := New(time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC))
	c := New(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC))
	if Max(a, b, c).Year() != 2018 {
		t.Fatalf("max")
	}
	if Min(a, b, c).Year() != 2016 {
		t.Fatalf("min")
	}
	if Max().IsValid() {
		t.Fatalf("empty max should be invalid")
	}
	if Max(a, Invalid()).Year() != 2017 {
		t.Fatalf("max ignoring invalid")
	}
}

func TestStrictAndMultiFormatParse(t *testing.T) {
	if _, err := ParseFormatStrict("2017-7-4", "YYYY-MM-DD"); err == nil {
		t.Fatalf("strict should reject single-digit month")
	}
	if m, err := ParseFormatStrict("2017-07-04", "YYYY-MM-DD"); err != nil || m.Format("YYYY-MM-DD") != "2017-07-04" {
		t.Fatalf("strict valid: %v %q", err, m.ISO())
	}
	if _, err := ParseFormatStrict("2017-07-04 extra", "YYYY-MM-DD"); err == nil {
		t.Fatalf("strict should reject trailing input")
	}
	m, err := ParseFormats("04.07.2017", []string{"YYYY-MM-DD", "DD/MM/YYYY", "DD.MM.YYYY"})
	if err != nil || m.Format("YYYY-MM-DD") != "2017-07-04" {
		t.Fatalf("multi-format: %v %q", err, m.ISO())
	}
	if _, err := ParseFormats("nope", []string{"YYYY-MM-DD"}); err == nil {
		t.Fatalf("multi-format should fail")
	}
}

func TestParseNamesAndMeridiem(t *testing.T) {
	m, err := ParseFormat("July 4 2017 2:05 pm", "MMMM D YYYY h:mm a")
	if err != nil {
		t.Fatalf("parse names: %v", err)
	}
	if m.Format("YYYY-MM-DD HH:mm") != "2017-07-04 14:05" {
		t.Fatalf("parse names result = %q", m.ISO())
	}
	fr, err := ParseFormatLocale("4 juillet 2017", "D MMMM YYYY", "fr")
	if err != nil || fr.Format("YYYY-MM-DD") != "2017-07-04" {
		t.Fatalf("fr parse: %v %q", err, fr.ISO())
	}
	iso, err := ParseFormat("2017-07-04T14:05:09+05:30", "YYYY-MM-DDTHH:mm:ssZ")
	if err != nil {
		t.Fatalf("offset parse: %v", err)
	}
	if iso.UTCOffset() != 330 {
		t.Fatalf("offset = %d", iso.UTCOffset())
	}
	if iso.UTC().Format("HH:mm") != "08:35" {
		t.Fatalf("offset instant = %q", iso.UTC().Format("HH:mm"))
	}
}

func TestParseUnixTokens(t *testing.T) {
	m, err := ParseFormat("1499177109", "X")
	if err != nil || m.Unix() != 1499177109 {
		t.Fatalf("parse X: %v %d", err, m.Unix())
	}
	mx, err := ParseFormat("1499177109123", "x")
	if err != nil || mx.UnixMilli() != 1499177109123 {
		t.Fatalf("parse x: %v %d", err, mx.UnixMilli())
	}
}

func TestParseDayOfYear(t *testing.T) {
	m, err := ParseFormat("2017 185", "YYYY DDD")
	if err != nil || m.Format("YYYY-MM-DD") != "2017-07-04" {
		t.Fatalf("parse DDD: %v %q", err, m.ISO())
	}
}

func TestRFC2822AndISO(t *testing.T) {
	m, err := ParseRFC2822("Tue, 04 Jul 2017 14:05:09 +0000")
	if err != nil || m.Format("YYYY-MM-DD HH:mm:ss") != "2017-07-04 14:05:09" {
		t.Fatalf("rfc2822: %v %q", err, m.ISO())
	}
	if _, err := ParseRFC2822("garbage"); err == nil {
		t.Fatalf("rfc2822 should fail")
	}
	iso, err := ParseISO("2017-07-04T14:05:09Z")
	if err != nil || iso.Year() != 2017 {
		t.Fatalf("iso: %v", err)
	}
}

func TestToISOStringAndValueOf(t *testing.T) {
	m := fixedMoment()
	if got := m.ToISOString(); got != "2017-07-04T14:05:09.123Z" {
		t.Fatalf("ToISOString = %q", got)
	}
	if m.ValueOf() != m.UnixMilli() {
		t.Fatalf("ValueOf")
	}
}

func TestInvalidAndCreationData(t *testing.T) {
	if Invalid().IsValid() {
		t.Fatalf("invalid should not be valid")
	}
	if got := Invalid().Format("YYYY"); got != "Invalid date" {
		t.Fatalf("invalid format = %q", got)
	}
	m, _ := ParseFormat("2017-07-04", "YYYY-MM-DD")
	cd := m.CreationData()
	if cd == nil || cd.Format != "YYYY-MM-DD" || !cd.Valid {
		t.Fatalf("creation data = %+v", cd)
	}
	if _, err := ParseFormat("bad", "YYYY-MM-DD"); err == nil {
		t.Fatalf("bad parse should error")
	}
}

func TestLocaleRegistry(t *testing.T) {
	if _, ok := LookupLocale("EN-US"); !ok {
		t.Fatalf("en-us should fall back to en")
	}
	if _, ok := LookupLocale("does-not-exist"); ok {
		t.Fatalf("unknown locale found")
	}
	if len(AvailableLocales()) < 20 {
		t.Fatalf("expected >=20 locales, got %d", len(AvailableLocales()))
	}
	prev := GlobalLocale()
	defer SetGlobalLocale(prev)
	if !SetGlobalLocale("fr") {
		t.Fatalf("set global fr failed")
	}
	if GlobalLocale() != "fr" {
		t.Fatalf("global not fr")
	}
	if got := fixedMoment().Format("MMMM"); got != "juillet" {
		t.Fatalf("global locale format = %q", got)
	}
	if SetGlobalLocale("nope") {
		t.Fatalf("set unknown global should fail")
	}
	m := fixedMoment().Locale("de")
	if m.LocaleName() != "de" {
		t.Fatalf("locale name = %q", m.LocaleName())
	}
	if m.LocaleData().Name != "de" {
		t.Fatalf("locale data")
	}
	if got := fixedMoment().Locale("unknown"); got.LocaleName() != "fr" {
		t.Fatalf("unknown locale should fall back to global, got %q", got.LocaleName())
	}
}

func TestUTCOffsetAndDST(t *testing.T) {
	m := fixedMoment().SetUTCOffset(120)
	if m.UTCOffset() != 120 {
		t.Fatalf("utc offset = %d", m.UTCOffset())
	}
	nyc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("tzdata unavailable")
	}
	summer := New(time.Date(2017, 7, 4, 12, 0, 0, 0, nyc))
	winter := New(time.Date(2017, 1, 4, 12, 0, 0, 0, nyc))
	if !summer.IsDST() {
		t.Fatalf("summer should be DST")
	}
	if winter.IsDST() {
		t.Fatalf("winter should not be DST")
	}
}

func TestWeeksInYear(t *testing.T) {
	// 2015 has 53 ISO weeks.
	m := New(time.Date(2015, 6, 1, 0, 0, 0, 0, time.UTC))
	if m.ISOWeeksInYear() != 53 {
		t.Fatalf("2015 iso weeks = %d", m.ISOWeeksInYear())
	}
	if fixedMoment().WeeksInYear() < 52 {
		t.Fatalf("weeks in year")
	}
}
