# tz-offset-normalizer

Every system that lets a human type in a UTC offset ends up with a
different mess: `+5`, `+0530`, `UTC+5:30`, `GMT-8`, `z`, ` +05:00 `. They
all mean the same thing, but string-equal comparisons and downstream
parsers choke on the variety. This library takes any of those forms and
turns it into one canonical shape: `+HH:MM`.

It does not try to be a full timezone database (no DST rules, no IANA
zone names like `America/Chicago`). It normalizes the raw numeric offset
and a small set of unambiguous abbreviations. If you need "what offset
does Chicago use on this date", reach for the standard library's
`time.LoadLocation` instead — this package is for the messy string that
arrives before you get anywhere near that.

## Usage

```go
package main

import (
	"fmt"
	"log"

	tzfmt "github.com/adam-heath/tz-offset-normalizer"
)

func main() {
	inputs := []string{"+5", "UTC+5:30", "-0800", "z", "  GMT-4  "}
	for _, in := range inputs {
		out, err := tzfmt.NormalizeOffset(in)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10q -> %s\n", in, out)
	}

	// If you need the raw minute count instead of a string (e.g. to shift
	// a time.Time), use OffsetMinutes.
	minutes, _ := tzfmt.OffsetMinutes("+5:30")
	fmt.Println(minutes) // 330

	// A small set of unambiguous abbreviations resolve directly.
	if offset, ok := tzfmt.NormalizeAbbreviation("PST"); ok {
		fmt.Println(offset) // -08:00
	}
}
```

Output:

```
"+5"       -> +05:00
"UTC+5:30" -> +05:30
"-0800"    -> -08:00
"z"        -> +00:00
"  GMT-4  "-> -04:00
```

## Design

Every exported function is pure: same input, same output, no reliance on
the local clock or the OS timezone database. That makes the whole
package testable with plain table-driven tests and safe to call from
anywhere without side effects.

Offsets are validated against the real range a UTC offset can take,
`-12:00` to `+14:00`. Anything outside that is treated as a parse error
rather than silently accepted.

Ambiguous abbreviations (`IST`, `CST` outside North America, `BST`) are
deliberately not in the lookup table — guessing wrong is worse than
returning "unknown" and asking the caller for an explicit offset.

## Status

Early skeleton. See the roadmap in the repo description for what's
planned next.
