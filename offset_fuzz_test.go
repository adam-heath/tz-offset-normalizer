package tzfmt

import "testing"

// FuzzOffsetMinutesRoundTrip checks the invariant that ties the three
// offset functions together: whatever OffsetMinutes accepts must produce
// a value FormatOffsetMinutes can turn back into a string that
// OffsetMinutes parses to the exact same number of minutes, and that
// NormalizeOffset agrees with FormatOffsetMinutes on that string. This
// is where a stray edge case in splitHoursMinutes (an input that parses
// but formats into something that no longer parses the same way) would
// show up.
func FuzzOffsetMinutesRoundTrip(f *testing.F) {
	seeds := []string{
		"Z", "z", "UTC", "GMT", "+5", "-5", "+05", "+0530", "-0530",
		"+5:30", "  +5:30  ", "UTC+5:30", "GMT-8", "utc+05:00",
		"+14:00", "-12:00", "", "garbage", "+15:00", "-13:00", "++5",
		"+5:99", "+5:30:00", "-0000", "+00", "UTC-0",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, raw string) {
		minutes, err := OffsetMinutes(raw)
		if err != nil {
			return
		}
		if minutes < minOffsetMinutes || minutes > maxOffsetMinutes {
			t.Fatalf("OffsetMinutes(%q) = %d, outside the -12:00..+14:00 range it's supposed to enforce", raw, minutes)
		}

		formatted, err := FormatOffsetMinutes(minutes)
		if err != nil {
			t.Fatalf("FormatOffsetMinutes(%d) returned an error for a value OffsetMinutes(%q) just accepted: %v", minutes, raw, err)
		}

		roundTripped, err := OffsetMinutes(formatted)
		if err != nil {
			t.Fatalf("OffsetMinutes(%q), the formatted form of %q, returned an error: %v", formatted, raw, err)
		}
		if roundTripped != minutes {
			t.Fatalf("round trip through %q: got %d minutes, want %d (original input %q)", formatted, roundTripped, minutes, raw)
		}

		normalized, err := NormalizeOffset(raw)
		if err != nil {
			t.Fatalf("NormalizeOffset(%q) returned an error but OffsetMinutes(%q) accepted it: %v", raw, raw, err)
		}
		if normalized != formatted {
			t.Fatalf("NormalizeOffset(%q) = %q, want %q to match FormatOffsetMinutes(%d)", raw, normalized, formatted, minutes)
		}
	})
}
