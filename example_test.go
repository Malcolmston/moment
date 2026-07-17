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
