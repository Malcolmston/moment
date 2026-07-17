package moment

import (
	"fmt"
	"strconv"
)

// This file bundles a representative set of about twenty common locales. It is
// deliberately not the full ~140-locale moment.js catalogue: the Locale type
// and RegisterLocale mechanism let applications add any locale they need, and
// the locales here demonstrate every feature (custom meridiems, ordinals,
// plural functions and week rules). See AvailableLocales for the shipped set.

// enOrdinal implements the English ordinal suffix (1st, 2nd, 3rd, 4th, 11th …).
func enOrdinal(number int, _ string) string {
	b := number % 10
	suffix := "th"
	if number%100 < 10 || number%100 > 20 {
		switch b {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}
	return strconv.Itoa(number) + suffix
}

// dotOrdinal renders the German/Nordic style "1." ordinal.
func dotOrdinal(number int, _ string) string { return strconv.Itoa(number) + "." }

// mascOrdinal renders the Romance-language masculine ordinal "1º".
func mascOrdinal(number int, _ string) string { return strconv.Itoa(number) + "º" }

// englishLocale is the built-in default; it is also registered as "en".
var englishLocale = &Locale{
	Name:                  "en",
	Months:                []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
	MonthsShort:           []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
	Weekdays:              []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	WeekdaysShort:         []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
	WeekdaysMin:           []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"},
	AM:                    "AM",
	PM:                    "PM",
	Ordinal:               enOrdinal,
	FirstDayOfWeek:        0,
	FirstWeekContainsDate: 6,
	RelativeTime: RelativeTime{
		Future: "in %s", Past: "%s ago",
		Second: "a few seconds", Seconds: "%d seconds",
		Minute: "a minute", Minutes: "%d minutes",
		Hour: "an hour", Hours: "%d hours",
		Day: "a day", Days: "%d days",
		Week: "a week", Weeks: "%d weeks",
		Month: "a month", Months: "%d months",
		Year: "a year", Years: "%d years",
	},
	Calendar: CalendarFormats{
		SameDay:  "[Today at] LT",
		NextDay:  "[Tomorrow at] LT",
		NextWeek: "dddd [at] LT",
		LastDay:  "[Yesterday at] LT",
		LastWeek: "[Last] dddd [at] LT",
		SameElse: "L",
	},
	LongDateFormats: LongDateFormats{
		LT: "h:mm A", LTS: "h:mm:ss A",
		L: "MM/DD/YYYY", LL: "MMMM D, YYYY",
		LLL: "MMMM D, YYYY h:mm A", LLLL: "dddd, MMMM D, YYYY h:mm A",
	},
}

// registerBundledLocales installs englishLocale and the other bundled locales.
func init() {
	RegisterLocale(englishLocale)

	RegisterLocale(&Locale{
		Name:                  "en-gb",
		Months:                englishLocale.Months,
		MonthsShort:           englishLocale.MonthsShort,
		Weekdays:              englishLocale.Weekdays,
		WeekdaysShort:         englishLocale.WeekdaysShort,
		WeekdaysMin:           englishLocale.WeekdaysMin,
		AM:                    "AM",
		PM:                    "PM",
		Ordinal:               enOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime:          englishLocale.RelativeTime,
		Calendar:              englishLocale.Calendar,
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD/MM/YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd, D MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:          "fr",
		Months:        []string{"janvier", "février", "mars", "avril", "mai", "juin", "juillet", "août", "septembre", "octobre", "novembre", "décembre"},
		MonthsShort:   []string{"janv.", "févr.", "mars", "avr.", "mai", "juin", "juil.", "août", "sept.", "oct.", "nov.", "déc."},
		Weekdays:      []string{"dimanche", "lundi", "mardi", "mercredi", "jeudi", "vendredi", "samedi"},
		WeekdaysShort: []string{"dim.", "lun.", "mar.", "mer.", "jeu.", "ven.", "sam."},
		WeekdaysMin:   []string{"di", "lu", "ma", "me", "je", "ve", "sa"},
		AM:            "AM",
		PM:            "PM",
		Ordinal: func(n int, _ string) string {
			if n == 1 {
				return "1er"
			}
			return strconv.Itoa(n) + "e"
		},
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "dans %s", Past: "il y a %s",
			Second: "quelques secondes", Seconds: "%d secondes",
			Minute: "une minute", Minutes: "%d minutes",
			Hour: "une heure", Hours: "%d heures",
			Day: "un jour", Days: "%d jours",
			Week: "une semaine", Weeks: "%d semaines",
			Month: "un mois", Months: "%d mois",
			Year: "un an", Years: "%d ans",
		},
		Calendar: CalendarFormats{
			SameDay: "[Aujourd’hui à] LT", NextDay: "[Demain à] LT",
			NextWeek: "dddd [à] LT", LastDay: "[Hier à] LT",
			LastWeek: "dddd [dernier à] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD/MM/YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd D MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "de",
		Months:                []string{"Januar", "Februar", "März", "April", "Mai", "Juni", "Juli", "August", "September", "Oktober", "November", "Dezember"},
		MonthsShort:           []string{"Jan.", "Feb.", "März", "Apr.", "Mai", "Juni", "Juli", "Aug.", "Sep.", "Okt.", "Nov.", "Dez."},
		Weekdays:              []string{"Sonntag", "Montag", "Dienstag", "Mittwoch", "Donnerstag", "Freitag", "Samstag"},
		WeekdaysShort:         []string{"So.", "Mo.", "Di.", "Mi.", "Do.", "Fr.", "Sa."},
		WeekdaysMin:           []string{"So", "Mo", "Di", "Mi", "Do", "Fr", "Sa"},
		AM:                    "vorm.",
		PM:                    "nachm.",
		Ordinal:               dotOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "in %s", Past: "vor %s",
			Pluralize: germanPlural,
		},
		Calendar: CalendarFormats{
			SameDay: "[heute um] LT [Uhr]", NextDay: "[morgen um] LT [Uhr]",
			NextWeek: "dddd [um] LT [Uhr]", LastDay: "[gestern um] LT [Uhr]",
			LastWeek: "[letzten] dddd [um] LT [Uhr]", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD.MM.YYYY", LL: "D. MMMM YYYY",
			LLL: "D. MMMM YYYY HH:mm", LLLL: "dddd, D. MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "es",
		Months:                []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"},
		MonthsShort:           []string{"ene.", "feb.", "mar.", "abr.", "may.", "jun.", "jul.", "ago.", "sep.", "oct.", "nov.", "dic."},
		Weekdays:              []string{"domingo", "lunes", "martes", "miércoles", "jueves", "viernes", "sábado"},
		WeekdaysShort:         []string{"dom.", "lun.", "mar.", "mié.", "jue.", "vie.", "sáb."},
		WeekdaysMin:           []string{"do", "lu", "ma", "mi", "ju", "vi", "sá"},
		AM:                    "a. m.",
		PM:                    "p. m.",
		Ordinal:               mascOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "en %s", Past: "hace %s",
			Second: "unos segundos", Seconds: "%d segundos",
			Minute: "un minuto", Minutes: "%d minutos",
			Hour: "una hora", Hours: "%d horas",
			Day: "un día", Days: "%d días",
			Week: "una semana", Weeks: "%d semanas",
			Month: "un mes", Months: "%d meses",
			Year: "un año", Years: "%d años",
		},
		Calendar: CalendarFormats{
			SameDay: "[hoy a las] LT", NextDay: "[mañana a las] LT",
			NextWeek: "dddd [a las] LT", LastDay: "[ayer a las] LT",
			LastWeek: "[el] dddd [pasado a las] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "H:mm", LTS: "H:mm:ss",
			L: "DD/MM/YYYY", LL: "D [de] MMMM [de] YYYY",
			LLL: "D [de] MMMM [de] YYYY H:mm", LLLL: "dddd, D [de] MMMM [de] YYYY H:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "it",
		Months:                []string{"gennaio", "febbraio", "marzo", "aprile", "maggio", "giugno", "luglio", "agosto", "settembre", "ottobre", "novembre", "dicembre"},
		MonthsShort:           []string{"gen", "feb", "mar", "apr", "mag", "giu", "lug", "ago", "set", "ott", "nov", "dic"},
		Weekdays:              []string{"domenica", "lunedì", "martedì", "mercoledì", "giovedì", "venerdì", "sabato"},
		WeekdaysShort:         []string{"dom", "lun", "mar", "mer", "gio", "ven", "sab"},
		WeekdaysMin:           []string{"do", "lu", "ma", "me", "gi", "ve", "sa"},
		AM:                    "AM",
		PM:                    "PM",
		Ordinal:               mascOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "tra %s", Past: "%s fa",
			Second: "alcuni secondi", Seconds: "%d secondi",
			Minute: "un minuto", Minutes: "%d minuti",
			Hour: "un'ora", Hours: "%d ore",
			Day: "un giorno", Days: "%d giorni",
			Week: "una settimana", Weeks: "%d settimane",
			Month: "un mese", Months: "%d mesi",
			Year: "un anno", Years: "%d anni",
		},
		Calendar: CalendarFormats{
			SameDay: "[Oggi alle] LT", NextDay: "[Domani alle] LT",
			NextWeek: "dddd [alle] LT", LastDay: "[Ieri alle] LT",
			LastWeek: "[lo scorso] dddd [alle] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD/MM/YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd D MMMM YYYY HH:mm",
		},
	})

	ptMonths := []string{"janeiro", "fevereiro", "março", "abril", "maio", "junho", "julho", "agosto", "setembro", "outubro", "novembro", "dezembro"}
	ptMonthsShort := []string{"jan", "fev", "mar", "abr", "mai", "jun", "jul", "ago", "set", "out", "nov", "dez"}
	ptWeekdays := []string{"domingo", "segunda-feira", "terça-feira", "quarta-feira", "quinta-feira", "sexta-feira", "sábado"}
	ptWeekdaysShort := []string{"dom", "seg", "ter", "qua", "qui", "sex", "sáb"}
	ptWeekdaysMin := []string{"do", "2ª", "3ª", "4ª", "5ª", "6ª", "sá"}
	ptRel := RelativeTime{
		Future: "em %s", Past: "há %s",
		Second: "poucos segundos", Seconds: "%d segundos",
		Minute: "um minuto", Minutes: "%d minutos",
		Hour: "uma hora", Hours: "%d horas",
		Day: "um dia", Days: "%d dias",
		Week: "uma semana", Weeks: "%d semanas",
		Month: "um mês", Months: "%d meses",
		Year: "um ano", Years: "%d anos",
	}
	RegisterLocale(&Locale{
		Name: "pt", Months: ptMonths, MonthsShort: ptMonthsShort,
		Weekdays: ptWeekdays, WeekdaysShort: ptWeekdaysShort, WeekdaysMin: ptWeekdaysMin,
		AM: "AM", PM: "PM", Ordinal: mascOrdinal,
		FirstDayOfWeek: 1, FirstWeekContainsDate: 4, RelativeTime: ptRel,
		Calendar: CalendarFormats{
			SameDay: "[Hoje às] LT", NextDay: "[Amanhã às] LT",
			NextWeek: "dddd [às] LT", LastDay: "[Ontem às] LT",
			LastWeek: "[Última] dddd [às] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD/MM/YYYY", LL: "D [de] MMMM [de] YYYY",
			LLL: "D [de] MMMM [de] YYYY HH:mm", LLLL: "dddd, D [de] MMMM [de] YYYY HH:mm",
		},
	})
	RegisterLocale(&Locale{
		Name: "pt-br", Months: ptMonths, MonthsShort: ptMonthsShort,
		Weekdays: ptWeekdays, WeekdaysShort: ptWeekdaysShort, WeekdaysMin: ptWeekdaysMin,
		AM: "AM", PM: "PM", Ordinal: mascOrdinal,
		FirstDayOfWeek: 0, FirstWeekContainsDate: 6, RelativeTime: ptRel,
		Calendar: CalendarFormats{
			SameDay: "[Hoje às] LT", NextDay: "[Amanhã às] LT",
			NextWeek: "dddd [às] LT", LastDay: "[Ontem às] LT",
			LastWeek: "[Última] dddd [às] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD/MM/YYYY", LL: "D [de] MMMM [de] YYYY",
			LLL: "D [de] MMMM [de] YYYY [às] HH:mm", LLLL: "dddd, D [de] MMMM [de] YYYY [às] HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:          "nl",
		Months:        []string{"januari", "februari", "maart", "april", "mei", "juni", "juli", "augustus", "september", "oktober", "november", "december"},
		MonthsShort:   []string{"jan.", "feb.", "mrt.", "apr.", "mei", "jun.", "jul.", "aug.", "sep.", "okt.", "nov.", "dec."},
		Weekdays:      []string{"zondag", "maandag", "dinsdag", "woensdag", "donderdag", "vrijdag", "zaterdag"},
		WeekdaysShort: []string{"zo.", "ma.", "di.", "wo.", "do.", "vr.", "za."},
		WeekdaysMin:   []string{"zo", "ma", "di", "wo", "do", "vr", "za"},
		AM:            "AM",
		PM:            "PM",
		Ordinal: func(n int, _ string) string {
			if n == 1 || n == 8 || n >= 20 {
				return strconv.Itoa(n) + "ste"
			}
			return strconv.Itoa(n) + "de"
		},
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "over %s", Past: "%s geleden",
			Second: "een paar seconden", Seconds: "%d seconden",
			Minute: "één minuut", Minutes: "%d minuten",
			Hour: "één uur", Hours: "%d uur",
			Day: "één dag", Days: "%d dagen",
			Week: "één week", Weeks: "%d weken",
			Month: "één maand", Months: "%d maanden",
			Year: "één jaar", Years: "%d jaar",
		},
		Calendar: CalendarFormats{
			SameDay: "[vandaag om] LT", NextDay: "[morgen om] LT",
			NextWeek: "dddd [om] LT", LastDay: "[gisteren om] LT",
			LastWeek: "[afgelopen] dddd [om] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD-MM-YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd D MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "ru",
		Months:                []string{"января", "февраля", "марта", "апреля", "мая", "июня", "июля", "августа", "сентября", "октября", "ноября", "декабря"},
		MonthsShort:           []string{"янв.", "февр.", "мар.", "апр.", "мая", "июня", "июля", "авг.", "сент.", "окт.", "нояб.", "дек."},
		Weekdays:              []string{"воскресенье", "понедельник", "вторник", "среда", "четверг", "пятница", "суббота"},
		WeekdaysShort:         []string{"вс", "пн", "вт", "ср", "чт", "пт", "сб"},
		WeekdaysMin:           []string{"вс", "пн", "вт", "ср", "чт", "пт", "сб"},
		AM:                    "ночи",
		PM:                    "дня",
		Ordinal:               dotOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 7,
		RelativeTime: RelativeTime{
			Future: "через %s", Past: "%s назад",
			Pluralize: russianPlural,
		},
		Calendar: CalendarFormats{
			SameDay: "[Сегодня, в] LT", NextDay: "[Завтра, в] LT",
			NextWeek: "dddd, [в] LT", LastDay: "[Вчера, в] LT",
			LastWeek: "[В прошлый] dddd, [в] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "H:mm", LTS: "H:mm:ss",
			L: "DD.MM.YYYY", LL: "D MMMM YYYY г.",
			LLL: "D MMMM YYYY г., H:mm", LLLL: "dddd, D MMMM YYYY г., H:mm",
		},
	})

	zhCnMeridiem := func(hour, minute int, _ bool) string {
		hm := hour*100 + minute
		switch {
		case hm < 600:
			return "凌晨"
		case hm < 900:
			return "早上"
		case hm < 1130:
			return "上午"
		case hm < 1230:
			return "中午"
		case hm < 1800:
			return "下午"
		default:
			return "晚上"
		}
	}
	RegisterLocale(&Locale{
		Name:          "zh-cn",
		Months:        []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"},
		MonthsShort:   []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		Weekdays:      []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
		WeekdaysShort: []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"},
		WeekdaysMin:   []string{"日", "一", "二", "三", "四", "五", "六"},
		Meridiem:      zhCnMeridiem,
		Ordinal: func(n int, token string) string {
			switch token {
			case "d", "D", "DDD":
				return strconv.Itoa(n) + "日"
			case "M":
				return strconv.Itoa(n) + "月"
			case "w", "W":
				return strconv.Itoa(n) + "周"
			default:
				return strconv.Itoa(n)
			}
		},
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "%s后", Past: "%s前",
			Second: "几秒", Seconds: "%d 秒",
			Minute: "1 分钟", Minutes: "%d 分钟",
			Hour: "1 小时", Hours: "%d 小时",
			Day: "1 天", Days: "%d 天",
			Week: "1 周", Weeks: "%d 周",
			Month: "1 个月", Months: "%d 个月",
			Year: "1 年", Years: "%d 年",
		},
		Calendar: CalendarFormats{
			SameDay: "[今天]LT", NextDay: "[明天]LT",
			NextWeek: "[下]ddddLT", LastDay: "[昨天]LT",
			LastWeek: "[上]ddddLT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "YYYY/MM/DD", LL: "YYYY年M月D日",
			LLL: "YYYY年M月D日 HH:mm", LLLL: "YYYY年M月D日dddd HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:          "zh-tw",
		Months:        []string{"一月", "二月", "三月", "四月", "五月", "六月", "七月", "八月", "九月", "十月", "十一月", "十二月"},
		MonthsShort:   []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		Weekdays:      []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"},
		WeekdaysShort: []string{"週日", "週一", "週二", "週三", "週四", "週五", "週六"},
		WeekdaysMin:   []string{"日", "一", "二", "三", "四", "五", "六"},
		Meridiem:      zhCnMeridiem,
		Ordinal: func(n int, token string) string {
			switch token {
			case "d", "D", "DDD":
				return strconv.Itoa(n) + "日"
			case "M":
				return strconv.Itoa(n) + "月"
			case "w", "W":
				return strconv.Itoa(n) + "週"
			default:
				return strconv.Itoa(n)
			}
		},
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "%s後", Past: "%s前",
			Second: "幾秒", Seconds: "%d 秒",
			Minute: "1 分鐘", Minutes: "%d 分鐘",
			Hour: "1 小時", Hours: "%d 小時",
			Day: "1 天", Days: "%d 天",
			Week: "1 週", Weeks: "%d 週",
			Month: "1 個月", Months: "%d 個月",
			Year: "1 年", Years: "%d 年",
		},
		Calendar: CalendarFormats{
			SameDay: "[今天] LT", NextDay: "[明天] LT",
			NextWeek: "[下]dddd LT", LastDay: "[昨天] LT",
			LastWeek: "[上]dddd LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "YYYY/MM/DD", LL: "YYYY年M月D日",
			LLL: "YYYY年M月D日 HH:mm", LLLL: "YYYY年M月D日dddd HH:mm",
		},
	})

	jaMeridiem := func(hour, _ int, _ bool) string {
		if hour < 12 {
			return "午前"
		}
		return "午後"
	}
	RegisterLocale(&Locale{
		Name:          "ja",
		Months:        []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		MonthsShort:   []string{"1月", "2月", "3月", "4月", "5月", "6月", "7月", "8月", "9月", "10月", "11月", "12月"},
		Weekdays:      []string{"日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"},
		WeekdaysShort: []string{"日", "月", "火", "水", "木", "金", "土"},
		WeekdaysMin:   []string{"日", "月", "火", "水", "木", "金", "土"},
		Meridiem:      jaMeridiem,
		Ordinal: func(n int, token string) string {
			if token == "d" || token == "D" || token == "DDD" {
				return strconv.Itoa(n) + "日"
			}
			return strconv.Itoa(n)
		},
		FirstDayOfWeek:        0,
		FirstWeekContainsDate: 6,
		RelativeTime: RelativeTime{
			Future: "%s後", Past: "%s前",
			Second: "数秒", Seconds: "%d秒",
			Minute: "1分", Minutes: "%d分",
			Hour: "1時間", Hours: "%d時間",
			Day: "1日", Days: "%d日",
			Week: "1週間", Weeks: "%d週間",
			Month: "1ヶ月", Months: "%dヶ月",
			Year: "1年", Years: "%d年",
		},
		Calendar: CalendarFormats{
			SameDay: "[今日] LT", NextDay: "[明日] LT",
			NextWeek: "[来週]dddd LT", LastDay: "[昨日] LT",
			LastWeek: "[先週]dddd LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "YYYY/MM/DD", LL: "YYYY年M月D日",
			LLL: "YYYY年M月D日 HH:mm", LLLL: "YYYY年M月D日 dddd HH:mm",
		},
	})

	koMeridiem := func(hour, _ int, _ bool) string {
		if hour < 12 {
			return "오전"
		}
		return "오후"
	}
	RegisterLocale(&Locale{
		Name:          "ko",
		Months:        []string{"1월", "2월", "3월", "4월", "5월", "6월", "7월", "8월", "9월", "10월", "11월", "12월"},
		MonthsShort:   []string{"1월", "2월", "3월", "4월", "5월", "6월", "7월", "8월", "9월", "10월", "11월", "12월"},
		Weekdays:      []string{"일요일", "월요일", "화요일", "수요일", "목요일", "금요일", "토요일"},
		WeekdaysShort: []string{"일", "월", "화", "수", "목", "금", "토"},
		WeekdaysMin:   []string{"일", "월", "화", "수", "목", "금", "토"},
		Meridiem:      koMeridiem,
		Ordinal: func(n int, token string) string {
			switch token {
			case "d", "D", "DDD":
				return strconv.Itoa(n) + "일"
			case "M":
				return strconv.Itoa(n) + "월"
			case "w", "W":
				return strconv.Itoa(n) + "주"
			default:
				return strconv.Itoa(n)
			}
		},
		FirstDayOfWeek:        0,
		FirstWeekContainsDate: 6,
		RelativeTime: RelativeTime{
			Future: "%s 후", Past: "%s 전",
			Second: "몇 초", Seconds: "%d초",
			Minute: "1분", Minutes: "%d분",
			Hour: "한 시간", Hours: "%d시간",
			Day: "하루", Days: "%d일",
			Week: "일주일", Weeks: "%d주",
			Month: "한 달", Months: "%d달",
			Year: "일 년", Years: "%d년",
		},
		Calendar: CalendarFormats{
			SameDay: "오늘 LT", NextDay: "내일 LT",
			NextWeek: "dddd LT", LastDay: "어제 LT",
			LastWeek: "지난주 dddd LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "A h:mm", LTS: "A h:mm:ss",
			L: "YYYY.MM.DD.", LL: "YYYY년 MMMM D일",
			LLL: "YYYY년 MMMM D일 A h:mm", LLLL: "YYYY년 MMMM D일 dddd A h:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "ar",
		Months:                []string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"},
		MonthsShort:           []string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"},
		Weekdays:              []string{"الأحد", "الإثنين", "الثلاثاء", "الأربعاء", "الخميس", "الجمعة", "السبت"},
		WeekdaysShort:         []string{"أحد", "إثنين", "ثلاثاء", "أربعاء", "خميس", "جمعة", "سبت"},
		WeekdaysMin:           []string{"ح", "ن", "ث", "ر", "خ", "ج", "س"},
		AM:                    "ص",
		PM:                    "م",
		Ordinal:               func(n int, _ string) string { return strconv.Itoa(n) },
		FirstDayOfWeek:        6,
		FirstWeekContainsDate: 12,
		RelativeTime: RelativeTime{
			Future: "بعد %s", Past: "منذ %s",
			Second: "ثانية واحدة", Seconds: "%d ثانية",
			Minute: "دقيقة واحدة", Minutes: "%d دقائق",
			Hour: "ساعة واحدة", Hours: "%d ساعات",
			Day: "يوم واحد", Days: "%d أيام",
			Week: "أسبوع واحد", Weeks: "%d أسابيع",
			Month: "شهر واحد", Months: "%d أشهر",
			Year: "عام واحد", Years: "%d أعوام",
		},
		Calendar: CalendarFormats{
			SameDay: "[اليوم عند الساعة] LT", NextDay: "[غدًا عند الساعة] LT",
			NextWeek: "dddd [عند الساعة] LT", LastDay: "[أمس عند الساعة] LT",
			LastWeek: "dddd [عند الساعة] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "D/M/YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd D MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "hi",
		Months:                []string{"जनवरी", "फ़रवरी", "मार्च", "अप्रैल", "मई", "जून", "जुलाई", "अगस्त", "सितम्बर", "अक्टूबर", "नवम्बर", "दिसम्बर"},
		MonthsShort:           []string{"जन.", "फ़र.", "मार्च", "अप्रै.", "मई", "जून", "जुल.", "अग.", "सित.", "अक्टू.", "नव.", "दिस."},
		Weekdays:              []string{"रविवार", "सोमवार", "मंगलवार", "बुधवार", "गुरूवार", "शुक्रवार", "शनिवार"},
		WeekdaysShort:         []string{"रवि", "सोम", "मंगल", "बुध", "गुरू", "शुक्र", "शनि"},
		WeekdaysMin:           []string{"र", "सो", "मं", "बु", "गु", "शु", "श"},
		AM:                    "सुबह",
		PM:                    "शाम",
		Ordinal:               func(n int, _ string) string { return strconv.Itoa(n) },
		FirstDayOfWeek:        0,
		FirstWeekContainsDate: 6,
		RelativeTime: RelativeTime{
			Future: "%s में", Past: "%s पहले",
			Second: "कुछ ही क्षण", Seconds: "%d सेकंड",
			Minute: "एक मिनट", Minutes: "%d मिनट",
			Hour: "एक घंटा", Hours: "%d घंटे",
			Day: "एक दिन", Days: "%d दिन",
			Week: "एक सप्ताह", Weeks: "%d सप्ताह",
			Month: "एक महीने", Months: "%d महीने",
			Year: "एक वर्ष", Years: "%d वर्ष",
		},
		Calendar: CalendarFormats{
			SameDay: "[आज] LT", NextDay: "[कल] LT",
			NextWeek: "dddd, LT", LastDay: "[कल] LT",
			LastWeek: "[पिछले] dddd, LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "A h:mm बजे", LTS: "A h:mm:ss बजे",
			L: "DD/MM/YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY, A h:mm बजे", LLLL: "dddd, D MMMM YYYY, A h:mm बजे",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "tr",
		Months:                []string{"Ocak", "Şubat", "Mart", "Nisan", "Mayıs", "Haziran", "Temmuz", "Ağustos", "Eylül", "Ekim", "Kasım", "Aralık"},
		MonthsShort:           []string{"Oca", "Şub", "Mar", "Nis", "May", "Haz", "Tem", "Ağu", "Eyl", "Eki", "Kas", "Ara"},
		Weekdays:              []string{"Pazar", "Pazartesi", "Salı", "Çarşamba", "Perşembe", "Cuma", "Cumartesi"},
		WeekdaysShort:         []string{"Paz", "Pzt", "Sal", "Çar", "Per", "Cum", "Cmt"},
		WeekdaysMin:           []string{"Pz", "Pt", "Sa", "Ça", "Pe", "Cu", "Ct"},
		AM:                    "ÖÖ",
		PM:                    "ÖS",
		Ordinal:               dotOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 7,
		RelativeTime: RelativeTime{
			Future: "%s sonra", Past: "%s önce",
			Second: "birkaç saniye", Seconds: "%d saniye",
			Minute: "bir dakika", Minutes: "%d dakika",
			Hour: "bir saat", Hours: "%d saat",
			Day: "bir gün", Days: "%d gün",
			Week: "bir hafta", Weeks: "%d hafta",
			Month: "bir ay", Months: "%d ay",
			Year: "bir yıl", Years: "%d yıl",
		},
		Calendar: CalendarFormats{
			SameDay: "[bugün saat] LT", NextDay: "[yarın saat] LT",
			NextWeek: "[gelecek] dddd [saat] LT", LastDay: "[dün] LT",
			LastWeek: "[geçen] dddd [saat] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD.MM.YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd, D MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "pl",
		Months:                []string{"styczeń", "luty", "marzec", "kwiecień", "maj", "czerwiec", "lipiec", "sierpień", "wrzesień", "październik", "listopad", "grudzień"},
		MonthsShort:           []string{"sty", "lut", "mar", "kwi", "maj", "cze", "lip", "sie", "wrz", "paź", "lis", "gru"},
		Weekdays:              []string{"niedziela", "poniedziałek", "wtorek", "środa", "czwartek", "piątek", "sobota"},
		WeekdaysShort:         []string{"ndz", "pon", "wt", "śr", "czw", "pt", "sob"},
		WeekdaysMin:           []string{"Nd", "Pn", "Wt", "Śr", "Cz", "Pt", "So"},
		AM:                    "AM",
		PM:                    "PM",
		Ordinal:               dotOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "za %s", Past: "%s temu",
			Pluralize: polishPlural,
		},
		Calendar: CalendarFormats{
			SameDay: "[Dziś o] LT", NextDay: "[Jutro o] LT",
			NextWeek: "dddd [o] LT", LastDay: "[Wczoraj o] LT",
			LastWeek: "[W zeszły] dddd [o] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "DD.MM.YYYY", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY HH:mm", LLLL: "dddd, D MMMM YYYY HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:          "sv",
		Months:        []string{"januari", "februari", "mars", "april", "maj", "juni", "juli", "augusti", "september", "oktober", "november", "december"},
		MonthsShort:   []string{"jan", "feb", "mar", "apr", "maj", "jun", "jul", "aug", "sep", "okt", "nov", "dec"},
		Weekdays:      []string{"söndag", "måndag", "tisdag", "onsdag", "torsdag", "fredag", "lördag"},
		WeekdaysShort: []string{"sön", "mån", "tis", "ons", "tor", "fre", "lör"},
		WeekdaysMin:   []string{"sö", "må", "ti", "on", "to", "fr", "lö"},
		AM:            "fm",
		PM:            "em",
		Ordinal: func(n int, _ string) string {
			r := n % 10
			if (r == 1 || r == 2) && n%100 != 11 && n%100 != 12 {
				return strconv.Itoa(n) + ":a"
			}
			return strconv.Itoa(n) + ":e"
		},
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "om %s", Past: "för %s sedan",
			Second: "några sekunder", Seconds: "%d sekunder",
			Minute: "en minut", Minutes: "%d minuter",
			Hour: "en timme", Hours: "%d timmar",
			Day: "en dag", Days: "%d dagar",
			Week: "en vecka", Weeks: "%d veckor",
			Month: "en månad", Months: "%d månader",
			Year: "ett år", Years: "%d år",
		},
		Calendar: CalendarFormats{
			SameDay: "[Idag] LT", NextDay: "[Imorgon] LT",
			NextWeek: "[På] dddd LT", LastDay: "[Igår] LT",
			LastWeek: "[I] dddd[s] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "HH:mm", LTS: "HH:mm:ss",
			L: "YYYY-MM-DD", LL: "D MMMM YYYY",
			LLL: "D MMMM YYYY [kl.] HH:mm", LLLL: "dddd D MMMM YYYY [kl.] HH:mm",
		},
	})

	RegisterLocale(&Locale{
		Name:                  "cs",
		Months:                []string{"leden", "únor", "březen", "duben", "květen", "červen", "červenec", "srpen", "září", "říjen", "listopad", "prosinec"},
		MonthsShort:           []string{"led", "úno", "bře", "dub", "kvě", "čvn", "čvc", "srp", "zář", "říj", "lis", "pro"},
		Weekdays:              []string{"neděle", "pondělí", "úterý", "středa", "čtvrtek", "pátek", "sobota"},
		WeekdaysShort:         []string{"ne", "po", "út", "st", "čt", "pá", "so"},
		WeekdaysMin:           []string{"ne", "po", "út", "st", "čt", "pá", "so"},
		AM:                    "dop.",
		PM:                    "odp.",
		Ordinal:               dotOrdinal,
		FirstDayOfWeek:        1,
		FirstWeekContainsDate: 4,
		RelativeTime: RelativeTime{
			Future: "za %s", Past: "před %s",
			Pluralize: czechPlural,
		},
		Calendar: CalendarFormats{
			SameDay: "[dnes v] LT", NextDay: "[zítra v] LT",
			NextWeek: "dddd [v] LT", LastDay: "[včera v] LT",
			LastWeek: "[minulý] dddd [v] LT", SameElse: "L",
		},
		LongDateFormats: LongDateFormats{
			LT: "H:mm", LTS: "H:mm:ss",
			L: "DD.MM.YYYY", LL: "D. MMMM YYYY",
			LLL: "D. MMMM YYYY H:mm", LLLL: "dddd D. MMMM YYYY H:mm",
		},
	})
}

// germanPlural implements German relative-time forms, which differ between the
// suffixed (dative) and unsuffixed (nominative) cases.
func germanPlural(number int, withoutSuffix bool, key string, isFuture bool) string {
	forms := map[string][2]string{
		"s":  {"ein paar Sekunden", "ein paar Sekunden"},
		"ss": {fmt.Sprintf("%d Sekunden", number), fmt.Sprintf("%d Sekunden", number)},
		"m":  {"eine Minute", "einer Minute"},
		"mm": {fmt.Sprintf("%d Minuten", number), fmt.Sprintf("%d Minuten", number)},
		"h":  {"eine Stunde", "einer Stunde"},
		"hh": {fmt.Sprintf("%d Stunden", number), fmt.Sprintf("%d Stunden", number)},
		"d":  {"ein Tag", "einem Tag"},
		"dd": {fmt.Sprintf("%d Tage", number), fmt.Sprintf("%d Tagen", number)},
		"w":  {"eine Woche", "einer Woche"},
		"ww": {fmt.Sprintf("%d Wochen", number), fmt.Sprintf("%d Wochen", number)},
		"M":  {"ein Monat", "einem Monat"},
		"MM": {fmt.Sprintf("%d Monate", number), fmt.Sprintf("%d Monaten", number)},
		"y":  {"ein Jahr", "einem Jahr"},
		"yy": {fmt.Sprintf("%d Jahre", number), fmt.Sprintf("%d Jahren", number)},
	}
	pair, ok := forms[key]
	if !ok {
		return ""
	}
	if withoutSuffix {
		return pair[0]
	}
	return pair[1]
}

// russianPlural implements Russian relative-time pluralization.
func russianPlural(number int, withoutSuffix bool, key string, isFuture bool) string {
	plural := func(one, few, many string) string {
		n := number % 100
		if n >= 11 && n <= 14 {
			return many
		}
		switch number % 10 {
		case 1:
			return one
		case 2, 3, 4:
			return few
		default:
			return many
		}
	}
	format := func(one, few, many string) string {
		return fmt.Sprintf(plural(one, few, many), number)
	}
	switch key {
	case "s":
		return "несколько секунд"
	case "ss":
		return format("%d секунду", "%d секунды", "%d секунд")
	case "m":
		if withoutSuffix {
			return "минута"
		}
		return "минуту"
	case "mm":
		return format("%d минуту", "%d минуты", "%d минут")
	case "h":
		return "час"
	case "hh":
		return format("%d час", "%d часа", "%d часов")
	case "d":
		return "день"
	case "dd":
		return format("%d день", "%d дня", "%d дней")
	case "w":
		return "неделя"
	case "ww":
		return format("%d неделю", "%d недели", "%d недель")
	case "M":
		return "месяц"
	case "MM":
		return format("%d месяц", "%d месяца", "%d месяцев")
	case "y":
		return "год"
	case "yy":
		return format("%d год", "%d года", "%d лет")
	}
	return ""
}

// polishPlural implements Polish relative-time pluralization.
func polishPlural(number int, withoutSuffix bool, key string, isFuture bool) string {
	plural := func(one, few, many string) string {
		if number == 1 {
			return one
		}
		n := number % 100
		if (number%10 >= 2 && number%10 <= 4) && (n < 12 || n > 14) {
			return few
		}
		return many
	}
	format := func(one, few, many string) string {
		if number == 1 {
			return one
		}
		return fmt.Sprintf(plural(one, few, many), number)
	}
	switch key {
	case "s":
		return "kilka sekund"
	case "ss":
		return format("sekunda", "%d sekundy", "%d sekund")
	case "m":
		return "minuta"
	case "mm":
		return format("minuta", "%d minuty", "%d minut")
	case "h":
		return "godzina"
	case "hh":
		return format("godzina", "%d godziny", "%d godzin")
	case "d":
		return "dzień"
	case "dd":
		return format("dzień", "%d dni", "%d dni")
	case "w":
		return "tydzień"
	case "ww":
		return format("tydzień", "%d tygodnie", "%d tygodni")
	case "M":
		return "miesiąc"
	case "MM":
		return format("miesiąc", "%d miesiące", "%d miesięcy")
	case "y":
		return "rok"
	case "yy":
		return format("rok", "%d lata", "%d lat")
	}
	return ""
}

// czechPlural implements Czech relative-time pluralization.
func czechPlural(number int, withoutSuffix bool, key string, isFuture bool) string {
	few := number >= 2 && number <= 4
	format := func(one, fewForm, many string) string {
		if few {
			return fmt.Sprintf(fewForm, number)
		}
		return fmt.Sprintf(many, number)
	}
	switch key {
	case "s":
		return "pár sekund"
	case "ss":
		return format("", "%d sekundy", "%d sekund")
	case "m":
		return "minuta"
	case "mm":
		return format("", "%d minuty", "%d minut")
	case "h":
		return "hodina"
	case "hh":
		return format("", "%d hodiny", "%d hodin")
	case "d":
		return "den"
	case "dd":
		return format("", "%d dny", "%d dní")
	case "w":
		return "týden"
	case "ww":
		return format("", "%d týdny", "%d týdnů")
	case "M":
		return "měsíc"
	case "MM":
		return format("", "%d měsíce", "%d měsíců")
	case "y":
		return "rok"
	case "yy":
		return format("", "%d roky", "%d let")
	}
	return ""
}
