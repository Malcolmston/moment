# moment

Date and time parsing and formatting for Go — a [moment.js](https://momentjs.com)-style
API layered on top of the standard `time` package.

`moment` is implemented with the Go standard library only: no cgo, no third-party
dependencies. A `Moment` is an immutable wrapper around `time.Time`; every
manipulation method returns a new value and never mutates the receiver. The
clock is injectable, so relative-time helpers such as `FromNow` are deterministic
in tests.

## Install

```sh
go get github.com/malcolmston/moment
```

Requires Go 1.24 or newer.

## Quick start

```go
package main

import (
	"fmt"
	"time"

	"github.com/malcolmston/moment"
)

func main() {
	// Parse with moment-style tokens.
	m, _ := moment.ParseFormat("14/07/2017 02:40", "DD/MM/YYYY HH:mm")

	// Format with tokens.
	fmt.Println(m.Format("dddd, MMMM D, YYYY [at] h:mm A"))
	// Friday, July 14, 2017 at 2:40 AM

	// Manipulate (immutable — returns a new Moment).
	fmt.Println(m.Add(3, moment.Day).StartOf(moment.Day).ISO())
	// 2017-07-17T00:00:00Z

	// Compare and diff.
	later := m.Add(2, moment.Hour)
	fmt.Println(m.IsBefore(later), later.DiffInt(m, moment.Minute))
	// true 120

	// Relative time with a deterministic clock.
	clock := moment.FixedClock(time.Date(2017, 7, 14, 2, 40, 0, 0, time.UTC))
	fmt.Println(m.WithClock(clock).Add(2, moment.Hour).FromNow())
	// in 2 hours
}
```

## Format tokens

`Format` and `ParseFormat` translate moment-style tokens to Go layouts. Use
`FormatLayout` / `ParseLayout` when you would rather pass a raw Go reference
layout. Wrap literal text in square brackets, e.g. `"YYYY [year]"`.

| Token  | Meaning                | Example   |
| ------ | ---------------------- | --------- |
| `YYYY` | Four-digit year        | `2017`    |
| `YY`   | Two-digit year         | `17`      |
| `MMMM` | Full month name        | `July`    |
| `MMM`  | Short month name       | `Jul`     |
| `MM`   | Two-digit month        | `07`      |
| `M`    | Month                  | `7`       |
| `DD`   | Two-digit day of month | `04`      |
| `D`    | Day of month           | `4`       |
| `dddd` | Full weekday name      | `Tuesday` |
| `ddd`  | Short weekday name     | `Tue`     |
| `HH`   | Two-digit 24-hour      | `14`      |
| `H`    | 24-hour (padded)       | `14`      |
| `hh`   | Two-digit 12-hour      | `02`      |
| `h`    | 12-hour                | `2`       |
| `mm`   | Two-digit minute       | `05`      |
| `m`    | Minute                 | `5`       |
| `ss`   | Two-digit second       | `09`      |
| `s`    | Second                 | `9`       |
| `SSS`  | Milliseconds           | `000`     |
| `A`    | Upper meridiem         | `PM`      |
| `a`    | Lower meridiem         | `pm`      |
| `Z`    | Numeric zone offset    | `-07:00`  |
| `ZZ`   | Numeric zone offset    | `-0700`   |

> Note: Go's `time` package has no unpadded 24-hour token, so `H` renders the
> same as `HH`.

## API overview

- **Construct:** `New`, `FromTime`, `Now`, `NowWith`, `Unix`, `UnixMilli`,
  `DateTime`, `Parse`, `ParseFormat`, `ParseInLocation`, `ParseLayout`.
- **Format:** `Format`, `FormatLayout`, `ISO`, `String`.
- **Manipulate:** `Add`, `Subtract`, `AddDuration`, `StartOf`, `EndOf`, `Set`,
  `Clone`.
- **Query:** `Year`, `Month`, `Date`/`Day`, `Hour`, `Minute`, `Second`,
  `Millisecond`, `Nanosecond`, `Weekday`, `DayOfYear`, `ISOWeek`, `IsLeapYear`.
- **Compare:** `IsBefore`, `IsAfter`, `IsSame`, `IsSameOrBefore`,
  `IsSameOrAfter`, `IsBetween`, `IsSameUnit`.
- **Diff:** `Diff` (float), `DiffInt`, `DiffDuration`.
- **Relative:** `FromNow`, `From`, `To`, `ToNow`, `Calendar`, `CalendarNow`,
  and the package-level `Humanize`.
- **Time zones:** `In`, `UTC`, `Local`, `Location`.

Units accept the `Unit` constants (`Year`, `Month`, `Week`, `Day`, `Date`,
`Hour`, `Minute`, `Second`, `Millisecond`) as well as common moment.js aliases
such as `"days"` or `"h"`.

## License

See repository.
