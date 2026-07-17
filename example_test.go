package moment_test

import (
	"fmt"
	"time"

	"github.com/malcolmston/moment"
)

// Example demonstrates formatting, manipulation and deterministic relative time
// using a fixed clock so the output is stable.
func Example() {
	// A fixed instant: Friday, 14 July 2017 02:40:00 UTC.
	ref := time.Date(2017, time.July, 14, 2, 40, 0, 0, time.UTC)
	m := moment.New(ref).WithClock(moment.FixedClock(ref))

	fmt.Println(m.Format("dddd, MMMM D, YYYY [at] h:mm A"))
	fmt.Println(m.Add(3, moment.Day).Format("YYYY-MM-DD"))
	fmt.Println(m.StartOf(moment.Month).ISO())
	fmt.Println(m.Add(2, moment.Hour).FromNow())
	fmt.Println(m.Subtract(1, moment.Day).Calendar(m))

	// Output:
	// Friday, July 14, 2017 at 2:40 AM
	// 2017-07-17
	// 2017-07-01T00:00:00Z
	// in 2 hours
	// Yesterday at 2:40 AM
}

// ExampleMoment_Locale shows locale-aware formatting and relative time.
func ExampleMoment_Locale() {
	ref := time.Date(2017, time.July, 4, 14, 5, 0, 0, time.UTC)
	m := moment.New(ref).WithClock(moment.FixedClock(ref))

	fmt.Println(m.Locale("fr").Format("LLLL"))
	fmt.Println(m.Locale("de").Format("LL"))
	fmt.Println(m.Locale("fr").Add(3, moment.Day).FromNow())

	// Output:
	// mardi 4 juillet 2017 14:05
	// 4. Juli 2017
	// dans 3 jours
}

// ExampleDuration demonstrates durations, conversion and ISO-8601 round-trips.
func ExampleDuration() {
	d := moment.NewDuration(1, moment.Year).
		Add(moment.NewDuration(2, moment.Month)).
		Add(moment.NewDuration(10, moment.Day))

	fmt.Println(d.ISOString())
	fmt.Println(d.AsMonths() > 14)
	fmt.Println(moment.NewDuration(2, moment.Hour).Humanize(true))

	parsed, _ := moment.ParseDuration("P1Y2M10DT2H30M")
	fmt.Println(parsed.ISOString())

	// Output:
	// P1Y2M10D
	// true
	// in 2 hours
	// P1Y2M10DT2H30M
}
