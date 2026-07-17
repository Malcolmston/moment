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
  `DateTime`, `FromArray`, `FromObject`, `Invalid`.
- **Parse:** `Parse`, `ParseISO`, `ParseRFC2822`, `ParseFormat`,
  `ParseFormatStrict`, `ParseFormatLocale`, `ParseFormats`, `ParseInLocation`,
  `ParseLayout`, `ParseDuration`.
- **Format:** `Format`, `FormatLayout`, `ISO`, `ToISOString`, `String`.
- **Manipulate:** `Add`, `Subtract`, `AddDuration`, `StartOf`, `EndOf`, `Set`,
  `SetAll`, `SetUTCOffset`, `Clone`.
- **Query:** `Year`, `Quarter`, `Month`, `Date`/`Day`, `Hour`, `Minute`,
  `Second`, `Millisecond`, `Nanosecond`, `Weekday`, `DayOfYear`, `ISOWeek`,
  `ISOWeekYear`, `ISOWeekNumber`, `ISOWeekday`, `Week`, `WeekYear`,
  `DaysInMonth`, `WeeksInYear`, `ISOWeeksInYear`, `UTCOffset`, `IsDST`,
  `IsLeapYear`, `IsValid`, `CreationData`.
- **Compare:** `IsBefore`, `IsAfter`, `IsSame`, `IsSameOrBefore`,
  `IsSameOrAfter`, `IsBetween`, `IsSameUnit`, and package-level `Max` / `Min`.
- **Diff:** `Diff` (float), `DiffInt`, `DiffDuration`.
- **Relative:** `FromNow`, `From`, `To`, `ToNow`, `Calendar`, `CalendarNow`,
  and the package-level `Humanize`.
- **Durations:** `NewDuration`, `DurationBetween`, `DurationFromTime`,
  `ParseDuration`; methods `As`, `Get`, `Add`, `Subtract`, `Abs`, `Clone`,
  `Humanize`, `ISOString`, `Locale`.
- **Locales:** `Locale`, `LocaleName`, `LocaleData`, `RegisterLocale`,
  `LookupLocale`, `AvailableLocales`, `SetGlobalLocale`, `GlobalLocale`,
  `SetRelativeTimeThreshold`.
- **Time zones:** `In`, `UTC`, `Local`, `Location`.

Units accept the `Unit` constants (`Year`, `Quarter`, `Month`, `Week`,
`ISOWeek`, `Day`, `Date`, `DayOfYear`, `Hour`, `Minute`, `Second`,
`Millisecond`) as well as common moment.js aliases such as `"days"` or `"h"`.

## Locales

Formatting and relative time are locale-aware. Select a locale per value with
`Locale`, or set the process default with `SetGlobalLocale`:

```go
m.Locale("fr").Format("LLLL") // mardi 4 juillet 2017 14:05
m.Locale("de").FromNow()      // in 3 Tagen
```

About twenty common locales are bundled: `en`, `en-gb`, `fr`, `de`, `es`, `it`,
`pt`, `pt-br`, `nl`, `ru`, `zh-cn`, `zh-tw`, `ja`, `ko`, `ar`, `hi`, `tr`, `pl`,
`sv`, `cs`. The full moment.js catalogue of ~140 locales is **not** bundled, but
the `Locale` type and `RegisterLocale` let you add any locale you need.

## Durations

```go
d := moment.NewDuration(1, moment.Year).
	Add(moment.NewDuration(2, moment.Month))
d.AsMonths()      // 14
d.ISOString()     // P1Y2M
d.Humanize(true)  // in a year

moment.ParseDuration("P1Y2M10DT2H30M") // round-trips via ISOString
```

## Format tokens (extended)

Beyond the table above, `Format` and `ParseFormat` support the full moment
token set: `Q`/`Qo`, `Do`, `DDD`/`DDDo`/`DDDD`, `w`/`wo`/`ww`, `W`/`Wo`/`WW`,
`e`/`E`, `gg`/`gggg`, `GG`/`GGGG`, `k`/`kk`, `x`/`X`, `z`/`zz`, runs of `S` for
fractional seconds, and the long-date tokens `LT`, `LTS`, `L`, `LL`, `LLL`,
`LLLL`.

## License

See repository.
