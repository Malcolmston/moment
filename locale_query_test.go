package moment

import (
	"testing"
	"time"
)

func TestMonthAndWeekdayLists(t *testing.T) {
	if got := Months("en"); len(got) != 12 || got[0] != "January" || got[6] != "July" {
		t.Errorf("Months(en) = %v", got)
	}
	if got := MonthsShort("en"); got[6] != "Jul" {
		t.Errorf("MonthsShort(en)[6] = %q", got[6])
	}
	if got := Weekdays("en"); len(got) != 7 || got[0] != "Sunday" || got[5] != "Friday" {
		t.Errorf("Weekdays(en) = %v", got)
	}
	if got := WeekdaysShort("en")[1]; got != "Mon" {
		t.Errorf("WeekdaysShort(en)[1] = %q", got)
	}
	if got := WeekdaysMin("en")[1]; got != "Mo" {
		t.Errorf("WeekdaysMin(en)[1] = %q", got)
	}
	if got := Months("fr")[0]; got != "janvier" {
		t.Errorf("Months(fr)[0] = %q, want janvier", got)
	}
	// Returned slice is a copy: mutating it must not affect the locale.
	Months("en")[0] = "MUT"
	if MonthName("en", time.January) != "January" {
		t.Errorf("Months returned a live reference")
	}
}

func TestNamedLookups(t *testing.T) {
	if got := MonthName("fr", time.July); got != "juillet" {
		t.Errorf("MonthName(fr, July) = %q", got)
	}
	if got := MonthShortName("en", time.September); got != "Sep" {
		t.Errorf("MonthShortName(en, Sep) = %q", got)
	}
	if got := WeekdayName("en", time.Wednesday); got != "Wednesday" {
		t.Errorf("WeekdayName = %q", got)
	}
	if got := WeekdayShortName("en", time.Wednesday); got != "Wed" {
		t.Errorf("WeekdayShortName = %q", got)
	}
	if got := WeekdayMinName("en", time.Wednesday); got != "We" {
		t.Errorf("WeekdayMinName = %q", got)
	}
}

func TestOrdinalMeridiemWeekRules(t *testing.T) {
	if got := Ordinal("en", 2); got != "2nd" {
		t.Errorf("Ordinal(en,2) = %q", got)
	}
	if got := Ordinal("de", 2); got != "2." {
		t.Errorf("Ordinal(de,2) = %q", got)
	}
	if got := Meridiem("en", 13, 0, false); got != "PM" {
		t.Errorf("Meridiem PM = %q", got)
	}
	if got := Meridiem("en", 13, 0, true); got != "pm" {
		t.Errorf("Meridiem pm = %q", got)
	}
	if FirstDayOfWeek("en") != 0 || FirstDayOfWeek("fr") != 1 {
		t.Errorf("FirstDayOfWeek en/fr = %d/%d", FirstDayOfWeek("en"), FirstDayOfWeek("fr"))
	}
	if FirstWeekContainsDate("en") != 6 || FirstWeekContainsDate("fr") != 4 {
		t.Errorf("FirstWeekContainsDate en/fr = %d/%d", FirstWeekContainsDate("en"), FirstWeekContainsDate("fr"))
	}
}

func TestLongDateFormat(t *testing.T) {
	tests := map[string]string{
		"LT": "h:mm A", "L": "MM/DD/YYYY", "LL": "MMMM D, YYYY", "bogus": "",
	}
	for token, want := range tests {
		if got := LongDateFormat("en", token); got != want {
			t.Errorf("LongDateFormat(en,%q) = %q, want %q", token, got, want)
		}
	}
}

func TestNormalizeUnit(t *testing.T) {
	tests := []struct {
		in   string
		want Unit
		ok   bool
	}{
		{"days", Day, true},
		{"d", Day, true},
		{"M", Month, true},
		{"months", Month, true},
		{"y", Year, true},
		{"ms", Millisecond, true},
		{"Q", Quarter, true},
		{"fortnight", Unit("fortnight"), false},
	}
	for _, tc := range tests {
		got, ok := NormalizeUnit(tc.in)
		if got != tc.want || ok != tc.ok {
			t.Errorf("NormalizeUnit(%q) = (%s,%v), want (%s,%v)", tc.in, got, ok, tc.want, tc.ok)
		}
	}
}
