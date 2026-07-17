package moment

import (
	"testing"
	"time"
)

func TestDurationConvenienceGetters(t *testing.T) {
	d := NewDuration(1, Year).
		Add(NewDuration(2, Month)).
		Add(NewDuration(3, Day)).
		Add(NewDuration(4, Hour)).
		Add(NewDuration(5, Minute)).
		Add(NewDuration(6, Second)).
		Add(NewDuration(700, Millisecond))
	if d.Hours() != 4 || d.Minutes() != 5 || d.Seconds() != 6 || d.Milliseconds() != 700 {
		t.Fatalf("components = %d %d %d %d", d.Hours(), d.Minutes(), d.Seconds(), d.Milliseconds())
	}
	if d.AsYears() <= 1 || d.AsMonths() <= 14 {
		t.Fatalf("as years/months = %v %v", d.AsYears(), d.AsMonths())
	}
	_ = d.AsWeeks()
	_ = d.AsSeconds()
	_ = d.AsMilliseconds()
	if d.String() == "" {
		t.Fatalf("String empty")
	}
}

func TestOrdinalsByLocale(t *testing.T) {
	m := New(time.Date(2017, 7, 1, 0, 0, 0, 0, time.UTC))
	cases := []struct{ locale, want string }{
		{"de", "1."},
		{"es", "1º"},
		{"it", "1º"},
		{"fr", "1er"},
		{"nl", "1ste"},
		{"sv", "1:a"},
		{"tr", "1."},
	}
	for _, c := range cases {
		if got := m.Locale(c.locale).Format("Do"); got != c.want {
			t.Errorf("[%s] Do = %q, want %q", c.locale, got, c.want)
		}
	}
	// French non-1 uses "e".
	if got := New(time.Date(2017, 7, 4, 0, 0, 0, 0, time.UTC)).Locale("fr").Format("Do"); got != "4e" {
		t.Fatalf("fr 4th = %q", got)
	}
	if got := New(time.Date(2017, 7, 2, 0, 0, 0, 0, time.UTC)).Locale("nl").Format("Do"); got != "2de" {
		t.Fatalf("nl 2nd = %q", got)
	}
}

func TestSlavicPluralEdges(t *testing.T) {
	cases := []struct {
		locale string
		dur    Duration
		want   string
	}{
		{"ru", NewDuration(2, Day), "через 2 дня"},
		{"ru", NewDuration(3, Hour), "через 3 часа"},
		{"ru", NewDuration(5, Hour), "через 5 часов"},
		{"ru", NewDuration(1, Year), "через год"},
		{"pl", NewDuration(2, Day), "za 2 dni"},
		{"pl", NewDuration(5, Month), "za 5 miesięcy"},
		{"pl", NewDuration(1, Year), "za rok"},
		{"cs", NewDuration(2, Hour), "za 2 hodiny"},
		{"cs", NewDuration(5, Day), "za 5 dní"},
		{"cs", NewDuration(1, Month), "za měsíc"},
	}
	for _, c := range cases {
		if got := c.dur.Locale(c.locale).Humanize(true); got != c.want {
			t.Errorf("[%s] = %q, want %q", c.locale, got, c.want)
		}
	}
}

func TestWeekGetters(t *testing.T) {
	m := New(time.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC))
	// en: Sunday start, doy 6 -> Jan 1 2017 (Sunday) is week 1.
	if m.Week() != 1 {
		t.Fatalf("en week = %d", m.Week())
	}
	if m.WeekYear() != 2017 {
		t.Fatalf("week year = %d", m.WeekYear())
	}
	// fr: Monday start, doy 4 -> Jan 1 2017 belongs to week 52 of 2016.
	fr := m.Locale("fr")
	if fr.Week() != 52 || fr.WeekYear() != 2016 {
		t.Fatalf("fr week = %d/%d", fr.Week(), fr.WeekYear())
	}
	if New(fixed).WeekYear() != 2017 {
		t.Fatalf("weekyear")
	}
}

func TestParseExtraTokens(t *testing.T) {
	// Ordinal day, short weekday, and 24-hour k token.
	m, err := ParseFormat("Tue 4th July 2017", "ddd Do MMMM YYYY")
	if err != nil || m.Format("YYYY-MM-DD") != "2017-07-04" {
		t.Fatalf("ordinal/weekday parse: %v %q", err, m.ISO())
	}
	mo, err := ParseFormat("7th 2017", "Mo YYYY")
	if err != nil || mo.Month() != time.July {
		t.Fatalf("Mo parse: %v %v", err, mo.Month())
	}
	k, err := ParseFormat("24:00:00 2017-07-04", "kk:mm:ss YYYY-MM-DD")
	if err != nil || k.Hour() != 0 {
		t.Fatalf("k parse: %v %d", err, k.Hour())
	}
	q, err := ParseFormat("3 2017", "Q YYYY")
	if err != nil || q.Month() != time.July {
		t.Fatalf("Q parse: %v %v", err, q.Month())
	}
	yy, err := ParseFormat("04-07-17", "DD-MM-YY")
	if err != nil || yy.Year() != 2017 {
		t.Fatalf("YY parse: %v %d", err, yy.Year())
	}
	neg, err := ParseFormat("-44 01 01", "Y MM DD")
	if err != nil || neg.Year() != -44 {
		t.Fatalf("Y parse: %v %d", err, neg.Year())
	}
	frac, err := ParseFormat("2017-07-04 14:05:09.123", "YYYY-MM-DD HH:mm:ss.SSS")
	if err != nil || frac.Millisecond() != 123 {
		t.Fatalf("frac parse: %v %d", err, frac.Millisecond())
	}
}

func TestTimeGetter(t *testing.T) {
	m := New(fixed)
	if !m.Time().Equal(fixed) {
		t.Fatalf("Time getter")
	}
}

func TestDurationLocaleFallback(t *testing.T) {
	if got := NewDuration(3, Day).Locale("does-not-exist").Humanize(false); got != "3 days" {
		t.Fatalf("fallback = %q", got)
	}
}

func TestArabicHindi(t *testing.T) {
	m := New(fixed)
	if got := m.Locale("ar").Format("MMMM"); got != "يوليو" {
		t.Fatalf("ar month = %q", got)
	}
	if got := m.Locale("hi").Format("MMMM"); got != "जुलाई" {
		t.Fatalf("hi month = %q", got)
	}
	if got := NewDuration(3, Day).Locale("ar").Humanize(true); got != "بعد 3 أيام" {
		t.Fatalf("ar humanize = %q", got)
	}
}
