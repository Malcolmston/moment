package moment

import (
	"reflect"
	"testing"
	"time"
)

var convRef = DateTime(2017, 7, 14, 2, 40, 9, 123*int(time.Millisecond), time.UTC)

func TestToArray(t *testing.T) {
	got := convRef.ToArray()
	want := []int{2017, 6, 14, 2, 40, 9, 123} // month is 0-based
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToArray = %v, want %v", got, want)
	}
	if Invalid().ToArray() != nil {
		t.Fatalf("invalid ToArray should be nil")
	}
	// Round-trip through FromArray.
	if rt := FromArray(got); !rt.IsSame(convRef) {
		t.Fatalf("FromArray(ToArray) = %v, want %v", rt.ISO(), convRef.ISO())
	}
}

func TestToObject(t *testing.T) {
	got := convRef.ToObject()
	want := map[string]int{
		"years": 2017, "months": 6, "date": 14,
		"hours": 2, "minutes": 40, "seconds": 9, "milliseconds": 123,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ToObject = %v, want %v", got, want)
	}
	if Invalid().ToObject() != nil {
		t.Fatalf("invalid ToObject should be nil")
	}
}

func TestToJSONAndISOStringZone(t *testing.T) {
	if got, want := convRef.ToJSON(), "2017-07-14T02:40:09.123Z"; got != want {
		t.Fatalf("ToJSON = %q, want %q", got, want)
	}
	if got := Invalid().ToJSON(); got != "" {
		t.Fatalf("invalid ToJSON = %q, want empty", got)
	}
	zone := time.FixedZone("CDT", -5*3600)
	m := DateTime(2017, 7, 14, 2, 40, 0, 0, zone)
	if got, want := m.ToISOStringZone(), "2017-07-14T02:40:00.000-05:00"; got != want {
		t.Fatalf("ToISOStringZone = %q, want %q", got, want)
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		unit Unit
		want int
	}{
		{Year, 2017},
		{Quarter, 3},
		{Month, 7}, // 1-based, unlike moment
		{Date, 14},
		{Day, 14},
		{DayOfYear, 195},
		{Hour, 2},
		{Minute, 40},
		{Second, 9},
		{Millisecond, 123},
		{ISOWeek, 28},
		{Unit("bogus"), 0},
	}
	for _, tc := range tests {
		if got := convRef.Get(tc.unit); got != tc.want {
			t.Errorf("Get(%s) = %d, want %d", tc.unit, got, tc.want)
		}
	}
}

func TestDaysInYear(t *testing.T) {
	if got := DateTime(2016, 3, 1, 0, 0, 0, 0, time.UTC).DaysInYear(); got != 366 {
		t.Errorf("2016 DaysInYear = %d, want 366", got)
	}
	if got := DateTime(2017, 3, 1, 0, 0, 0, 0, time.UTC).DaysInYear(); got != 365 {
		t.Errorf("2017 DaysInYear = %d, want 365", got)
	}
}

func TestZoneAndUTCPredicates(t *testing.T) {
	u := convRef.UTC()
	if !u.IsUTC() {
		t.Errorf("UTC moment should report IsUTC")
	}
	if u.IsLocal() && time.Local != time.UTC {
		t.Errorf("UTC moment should not report IsLocal")
	}
	if got := u.ZoneAbbr(); got != "UTC" {
		t.Errorf("ZoneAbbr = %q, want UTC", got)
	}
	if got := u.ZoneName(); got != "UTC" {
		t.Errorf("ZoneName = %q, want UTC", got)
	}
}

func TestLocaleWeekday(t *testing.T) {
	// 2017-07-14 is a Friday.
	if got := convRef.Locale("en").LocaleWeekday(); got != 5 { // en week starts Sunday
		t.Errorf("en LocaleWeekday = %d, want 5", got)
	}
	if got := convRef.Locale("fr").LocaleWeekday(); got != 4 { // fr week starts Monday
		t.Errorf("fr LocaleWeekday = %d, want 4", got)
	}
}

func TestIsPredicates(t *testing.T) {
	if !IsMoment(convRef) || IsMoment("x") {
		t.Errorf("IsMoment wrong")
	}
	if !IsDuration(NewDuration(1, Hour)) || IsDuration(convRef) {
		t.Errorf("IsDuration wrong")
	}
	if !IsDate(time.Now()) || IsDate(convRef) {
		t.Errorf("IsDate wrong")
	}
}
