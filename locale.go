package moment

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// RelativeTime holds the phrase templates a locale uses to render humanized
// durations such as "in 3 days" or "2 hours ago". Each field corresponds to a
// moment.js relative-time key; templates may contain the verb "%d" placeholder,
// which is replaced by the magnitude. Future and Past wrap the resulting phrase
// (their "%s" placeholder is replaced by the rendered magnitude).
type RelativeTime struct {
	Future  string // wraps a future phrase, e.g. "in %s"
	Past    string // wraps a past phrase, e.g. "%s ago"
	Second  string // key "s"  – a few seconds
	Seconds string // key "ss" – "%d seconds"
	Minute  string // key "m"  – a minute
	Minutes string // key "mm" – "%d minutes"
	Hour    string // key "h"  – an hour
	Hours   string // key "hh" – "%d hours"
	Day     string // key "d"  – a day
	Days    string // key "dd" – "%d days"
	Week    string // key "w"  – a week
	Weeks   string // key "ww" – "%d weeks"
	Month   string // key "M"  – a month
	Months  string // key "MM" – "%d months"
	Year    string // key "y"  – a year
	Years   string // key "yy" – "%d years"

	// Pluralize, when non-nil, overrides the field templates above. It is used
	// by locales with non-trivial plural rules (Slavic languages, Arabic). key
	// is one of "s","ss","m","mm","h","hh","d","dd","w","ww","M","MM","y","yy".
	Pluralize func(number int, withoutSuffix bool, key string, isFuture bool) string
}

// CalendarFormats holds the moment.js calendar templates used by Calendar to
// describe a Moment relative to a reference day. Each value is a format string
// (tokens plus [bracketed] literals), for example "[Today at] LT".
type CalendarFormats struct {
	SameDay  string
	NextDay  string
	NextWeek string
	LastDay  string
	LastWeek string
	SameElse string
}

// LongDateFormats holds the locale's named long date/time patterns, referenced
// from format strings and calendar templates by the tokens LT, LTS, L, LL, LLL
// and LLLL.
type LongDateFormats struct {
	LT   string // time, e.g. "h:mm A"
	LTS  string // time with seconds, e.g. "h:mm:ss A"
	L    string // numeric date, e.g. "MM/DD/YYYY"
	LL   string // long date, e.g. "MMMM D, YYYY"
	LLL  string // long date and time, e.g. "MMMM D, YYYY h:mm A"
	LLLL string // full date and time, e.g. "dddd, MMMM D, YYYY h:mm A"
}

// Locale describes everything language- and region-specific about rendering and
// parsing a Moment: month and weekday names, the meridiem, ordinals, the week
// numbering rules, and the relative-time, calendar and long-date templates. Use
// RegisterLocale to make a Locale available by name and Moment.Locale to select
// one per value.
type Locale struct {
	// Name is the canonical, lower-case locale identifier, e.g. "en" or "pt-br".
	Name string

	// Months and MonthsShort hold the 12 month names, January first.
	Months      []string
	MonthsShort []string

	// Weekdays, WeekdaysShort and WeekdaysMin hold the 7 weekday names, Sunday
	// first (index 0 == Sunday), matching time.Weekday ordering.
	Weekdays      []string
	WeekdaysShort []string
	WeekdaysMin   []string

	// AM and PM are the simple meridiem strings used when Meridiem is nil.
	AM string
	PM string
	// Meridiem, when non-nil, computes the meridiem string for the given hour
	// and minute; isLower requests the lower-case form. Locales such as Chinese
	// and Japanese use this for time-of-day words rather than AM/PM.
	Meridiem func(hour, minute int, isLower bool) string

	// Ordinal renders an ordinal number for the given format token (e.g. "D",
	// "Do", "M", "w"). English yields "1st", "2nd"; German yields "1.".
	Ordinal func(number int, token string) string

	// FirstDayOfWeek is the locale's first weekday (0 == Sunday, 1 == Monday),
	// used by week-of-year tokens and StartOf(Week). FirstWeekContainsDate is
	// moment's "doy": the January date that is always in week 1.
	FirstDayOfWeek        int
	FirstWeekContainsDate int

	// RelativeTime, Calendar and LongDateFormats supply the templates used by
	// FromNow/Humanize, Calendar and the long-date format tokens respectively.
	RelativeTime    RelativeTime
	Calendar        CalendarFormats
	LongDateFormats LongDateFormats
}

// weekRules returns the (dow, doy) pair used by the week-number algorithms,
// defaulting to the moment.js global default (Sunday start, doy 6) when unset.
func (l *Locale) weekRules() (dow, doy int) {
	dow = l.FirstDayOfWeek
	doy = l.FirstWeekContainsDate
	if doy == 0 {
		doy = 6
	}
	return dow, doy
}

// meridiem returns the meridiem string for the given clock time.
func (l *Locale) meridiem(hour, minute int, lower bool) string {
	if l.Meridiem != nil {
		return l.Meridiem(hour, minute, lower)
	}
	if hour < 12 {
		if lower {
			return strings.ToLower(l.AM)
		}
		return l.AM
	}
	if lower {
		return strings.ToLower(l.PM)
	}
	return l.PM
}

// ordinal renders number as an ordinal for the given token, falling back to the
// bare decimal when the locale has no ordinal function.
func (l *Locale) ordinal(number int, token string) string {
	if l.Ordinal != nil {
		return l.Ordinal(number, token)
	}
	return strconv.Itoa(number)
}

var (
	localeMu       sync.RWMutex
	localeRegistry = map[string]*Locale{}
	globalLocale   = "en"
)

// RegisterLocale adds or replaces a Locale in the global registry, keyed by its
// lower-cased Name. It is safe for concurrent use. Registering a locale named
// "en" replaces the built-in English default.
func RegisterLocale(l *Locale) {
	if l == nil || l.Name == "" {
		return
	}
	localeMu.Lock()
	defer localeMu.Unlock()
	key := strings.ToLower(l.Name)
	l.Name = key
	localeRegistry[key] = l
}

// LookupLocale returns the registered Locale for name and whether it was found.
// Lookup is case-insensitive and falls back from a regional variant to its base
// language (for example "en-au" falls back to "en").
func LookupLocale(name string) (*Locale, bool) {
	localeMu.RLock()
	defer localeMu.RUnlock()
	return lookupLocked(name)
}

// lookupLocked performs a registry lookup; callers must hold localeMu.
func lookupLocked(name string) (*Locale, bool) {
	key := strings.ToLower(strings.ReplaceAll(name, "_", "-"))
	if l, ok := localeRegistry[key]; ok {
		return l, true
	}
	if i := strings.IndexByte(key, '-'); i > 0 {
		if l, ok := localeRegistry[key[:i]]; ok {
			return l, true
		}
	}
	return nil, false
}

// AvailableLocales returns the sorted names of every registered locale. The
// standard build ships a representative set of about twenty common locales;
// applications may register additional ones with RegisterLocale.
func AvailableLocales() []string {
	localeMu.RLock()
	defer localeMu.RUnlock()
	out := make([]string, 0, len(localeRegistry))
	for name := range localeRegistry {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}

// SetGlobalLocale sets the default locale used by Moments that do not carry an
// explicit locale. It reports whether the name was found; an unknown name
// leaves the global default unchanged.
func SetGlobalLocale(name string) bool {
	localeMu.Lock()
	defer localeMu.Unlock()
	if l, ok := lookupLocked(name); ok {
		globalLocale = l.Name
		return true
	}
	return false
}

// GlobalLocale returns the name of the current global default locale.
func GlobalLocale() string {
	localeMu.RLock()
	defer localeMu.RUnlock()
	return globalLocale
}

// mustLocale returns the named locale, or the global default, or the built-in
// English locale, never nil.
func mustLocale(name string) *Locale {
	localeMu.RLock()
	defer localeMu.RUnlock()
	if name != "" {
		if l, ok := lookupLocked(name); ok {
			return l
		}
	}
	if l, ok := localeRegistry[globalLocale]; ok {
		return l
	}
	if l, ok := localeRegistry["en"]; ok {
		return l
	}
	return englishLocale
}

// Locale returns a copy of the Moment bound to the named locale, affecting
// Format, Calendar, FromNow and the humanized relative-time helpers. An unknown
// name leaves the Moment on the global default; use LookupLocale to test first.
func (m Moment) Locale(name string) Moment {
	if l, ok := LookupLocale(name); ok {
		m.loc = l
	} else {
		m.loc = nil
	}
	return m
}

// LocaleName returns the name of the locale bound to the Moment, or the global
// default when none is set.
func (m Moment) LocaleName() string {
	if m.loc != nil {
		return m.loc.Name
	}
	return GlobalLocale()
}

// LocaleData returns the effective Locale for the Moment: its bound locale, or
// the global default. The returned pointer must not be mutated.
func (m Moment) LocaleData() *Locale {
	if m.loc != nil {
		return m.loc
	}
	return mustLocale("")
}

// localeOf returns the effective locale for internal rendering.
func (m Moment) localeOf() *Locale {
	if m.loc != nil {
		return m.loc
	}
	return mustLocale("")
}
