# Changelog

All notable changes to this project are documented in this file. The format is
based on [Keep a Changelog](https://keepachangelog.com/).

## [0.2.0]

Major expansion toward moment.js parity: internationalization and durations.

### Added

- **Locale framework.** New `Locale` type covering month/weekday names (long,
  short, min), meridiem, ordinals, week rules (`FirstDayOfWeek`,
  `FirstWeekContainsDate`), relative-time templates, calendar templates and long
  date formats (LT/LTS/L/LL/LLL/LLLL). Registration and lookup via
  `RegisterLocale`, `LookupLocale`, `AvailableLocales`, `SetGlobalLocale`,
  `GlobalLocale`, plus per-Moment `Locale`, `LocaleName` and `LocaleData`.
- **Bundled locales.** A representative set of ~20 common locales: en, en-gb,
  fr, de, es, it, pt, pt-br, nl, ru, zh-cn, zh-tw, ja, ko, ar, hi, tr, pl, sv,
  cs. (The full ~140-locale moment.js catalogue is not bundled; the mechanism
  plus these are.) Format, Parse, FromNow and Calendar are all locale-aware.
- **Durations.** New `Duration` type: `NewDuration`, `DurationBetween`,
  `DurationFromTime`; `As`/`Get` (years…milliseconds), `Add`/`Subtract`, `Abs`,
  `Clone`, locale-aware `Humanize` with configurable relative-time thresholds
  (`SetRelativeTimeThreshold`), and ISO-8601 duration parse/format
  (`ParseDuration`, `ISOString`).
- **Format.** Full moment token set: `Q`, `Qo`, `Do`, `DDD`, `DDDo`, `DDDD`,
  `w`/`wo`/`ww`, `W`/`Wo`/`WW`, `e`, `E`, `gg`/`gggg`, `GG`/`GGGG`, `k`/`kk`,
  `x`/`X`, `z`/`zz`, runs of `S`, and long-date tokens LT/LTS/L/LL/LLL/LLLL.
- **Query & manipulation.** `Quarter`, `ISOWeekYear`, `ISOWeekNumber`,
  `ISOWeekday`, `Week`, `WeekYear`, `DaysInMonth`, `WeeksInYear`,
  `ISOWeeksInYear`, `UTCOffset`, `IsDST`; `SetAll` (set by object), `SetUTCOffset`,
  new `Quarter`, `ISOWeek` and `DayOfYear` units; package-level `Max`/`Min`.
- **Parsing.** Strict parsing (`ParseFormatStrict`), multiple-format parsing
  (`ParseFormats`), RFC 2822 (`ParseRFC2822`), ISO helper (`ParseISO`),
  array/object construction (`FromArray`, `FromObject`), validity and creation
  data (`Invalid`, `CreationData`), and `ToISOString`, `ValueOf`.

### Changed

- `Format`/`ParseFormat` now use a custom locale-aware token engine rather than
  translating to Go reference layouts, enabling the full token set and
  localized names.
- Relative-time and calendar helpers route through the new `Duration` and
  `Locale` machinery; English output is unchanged.

## [0.1.0]

- Initial release: immutable `Moment` over `time.Time`, moment-style token
  format/parse, manipulation, comparison, diff and English relative time.
