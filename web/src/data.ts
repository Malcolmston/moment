// Library content for the moment documentation site. Mirrors the shape used by
// the malcolmston/go landing site's data.ts so the sibling sites stay in sync.
export interface Lib {
  id: string; name: string; icon: string; accent: string; pkg: string; node: string;
  repo: string; docs: string; tagline: string; blurb: string; tags: string[];
  features: string[]; node_code: string; go_code: string; integrate: string;
}

export const NODE_ACCENT = '#8cc84b';

export const MOMENT: Lib = {
  id:"moment", name:"Moment", icon:'<i class="fa-solid fa-clock"></i>', accent:"#f59e0b",
  pkg:"github.com/malcolmston/moment", node:"moment/moment",
  repo:"https://github.com/malcolmston/moment", docs:"https://malcolmston.github.io/moment/",
  tagline:"moment.js-style dates and times in Go.",
  blurb:"A from-scratch, standard-library-only Go take on moment.js, layered directly on the time package with "+
    "no cgo and no third-party dependencies. A <code>Moment</code> is an immutable wrapper around "+
    "<code>time.Time</code>: every manipulation method returns a new value and never mutates the receiver, so "+
    "moments are safe to share. You get moment-token parsing and formatting (YYYY, MMMM, dddd, HH…), unit-based "+
    "Add/Subtract/StartOf/EndOf/Set arithmetic, float and integer Diff, comparison and query helpers, timezone "+
    "reinterpretation, and humanized relative time (FromNow, Calendar, Humanize) driven by an injectable Clock so "+
    "tests stay deterministic. The import path and package name are both moment.",
  tags:["immutable Moment","time.Time","moment tokens","Format/ParseFormat","Add/StartOf","Diff","FromNow","injectable Clock","zero deps"],
  features:[
    "Immutable <code>Moment</code> over <code>time.Time</code> — every method returns a new value; construct with <code>New</code>, <code>FromTime</code>, <code>Now</code>, <code>Unix</code>, <code>UnixMilli</code> or <code>DateTime</code>",
    "moment-token <code>Format</code> and <code>ParseFormat</code> (<code>YYYY</code>/<code>MMMM</code>/<code>dddd</code>/<code>HH</code>…) with <code>[literal]</code> escaping, plus raw-layout <code>FormatLayout</code>/<code>ParseLayout</code> and forgiving <code>Parse</code>",
    "Unit arithmetic — <code>Add</code>, <code>Subtract</code>, <code>AddDuration</code>, <code>StartOf</code>, <code>EndOf</code> and <code>Set</code> keyed by the <code>Unit</code> constants (<code>Year</code>…<code>Millisecond</code>) with moment.js aliases",
    "Comparison &amp; query — <code>IsBefore</code>, <code>IsAfter</code>, <code>IsSame</code>, <code>IsBetween</code>, <code>IsSameUnit</code> and getters like <code>Year</code>, <code>Weekday</code>, <code>DayOfYear</code>, <code>ISOWeek</code>",
    "Difference in any unit — <code>Diff</code> (float64), <code>DiffInt</code> (truncated) and <code>DiffDuration</code> (<code>time.Duration</code>), with moment.js-style fractional month math",
    "Humanized relative time — <code>FromNow</code>, <code>From</code>, <code>To</code>, <code>ToNow</code>, <code>Calendar</code> and package-level <code>Humanize</code> (\"in 3 days\", \"Today at 2:30 PM\")",
    "Deterministic clock — the injectable <code>Clock</code> interface with <code>FixedClock</code> and <code>WithClock</code> makes relative-time output reproducible in tests",
    "Time zones — <code>In</code>, <code>UTC</code> and <code>Local</code> reinterpret a moment in another <code>time.Location</code> without changing the instant",
    "Zero dependencies — pure Go standard library, no cgo, nothing to audit but the toolchain"
  ],
  node_code:
`const moment = require('moment');

const m = moment('14/07/2017 02:40', 'DD/MM/YYYY HH:mm');
console.log(m.format('dddd, MMMM D, YYYY [at] h:mm A'));
// Friday, July 14, 2017 at 2:40 AM

const later = m.clone().add(2, 'hours');
console.log(m.isBefore(later), later.diff(m, 'minutes'));
// true 120`,
  go_code:
`import "github.com/malcolmston/moment"

// The import path and the package name are both moment.
m, _ := moment.ParseFormat("14/07/2017 02:40", "DD/MM/YYYY HH:mm")
fmt.Println(m.Format("dddd, MMMM D, YYYY [at] h:mm A"))
// Friday, July 14, 2017 at 2:40 AM

later := m.Add(2, moment.Hour)
fmt.Println(m.IsBefore(later), later.DiffInt(m, moment.Minute))
// true 120`,
  integrate:
`<span class="tok-c">// Parse moment-style tokens, then manipulate immutably — every</span>
<span class="tok-c">// call returns a new Moment, so m is untouched.</span>
m, _ := moment.ParseFormat("14/07/2017 02:40", "DD/MM/YYYY HH:mm")
start := m.StartOf(moment.Day)
end := m.Add(3, moment.Day).EndOf(moment.Day)
fmt.Println(end.DiffInt(start, moment.Day)) // 3

<span class="tok-c">// Reinterpret the same instant in another time zone.</span>
ny, _ := time.LoadLocation("America/New_York")
fmt.Println(m.In(ny).Format("h:mm A Z"))

<span class="tok-c">// Inject a fixed clock so relative time is deterministic in tests.</span>
clock := moment.FixedClock(time.Date(2017, 7, 14, 2, 40, 0, 0, time.UTC))
soon := m.WithClock(clock).Add(2, moment.Hour)
fmt.Println(soon.FromNow())                       // in 2 hours
fmt.Println(soon.Calendar(moment.NowWith(clock))) // Today at 4:40 AM`
};
