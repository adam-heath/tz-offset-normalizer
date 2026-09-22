// Package tzfmt normalizes the many ways people and systems write UTC
// offsets and timezone abbreviations into one canonical shape.
//
// Every exported function here is pure: given the same input it always
// returns the same output, and none of them touch the clock, the
// filesystem, or the local timezone database. That makes them safe to
// call from anywhere (request handlers, CSV importers, form validation)
// without worrying about hidden state.
package tzfmt

import (
	"fmt"
	"strconv"
	"strings"
)

// minOffsetMinutes and maxOffsetMinutes bound what a real UTC offset can
// be. The extremes are Baker Island (UTC-12) and the Line Islands
// (UTC+14); anything outside that range is not a timezone, it's a typo.
const (
	minOffsetMinutes = -12 * 60
	maxOffsetMinutes = 14 * 60
)

// NormalizeOffset parses a UTC offset written in one of the common messy
// forms and returns it in canonical "+HH:MM" / "-HH:MM" form.
//
// Accepted inputs include:
//
//	"Z", "UTC", "GMT"       -> "+00:00"
//	"+5", "-5"              -> "+05:00", "-05:00"
//	"+0530", "-0530"        -> "+05:30", "-05:30"
//	"+5:30", "UTC+5:30"     -> "+05:30"
//	"  +05:00  "            -> "+05:00" (surrounding whitespace ignored)
func NormalizeOffset(raw string) (string, error) {
	minutes, err := OffsetMinutes(raw)
	if err != nil {
		return "", err
	}
	return FormatOffsetMinutes(minutes)
}

// OffsetMinutes parses the same messy inputs as NormalizeOffset but
// returns the signed number of minutes east of UTC instead of a string.
// This is useful when the caller needs to do arithmetic (e.g. convert a
// local time to UTC) rather than just display the offset.
func OffsetMinutes(raw string) (int, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return 0, fmt.Errorf("tzfmt: empty offset")
	}

	upper := strings.ToUpper(s)
	if upper == "Z" || upper == "UTC" || upper == "GMT" {
		return 0, nil
	}

	for _, label := range []string{"UTC", "GMT"} {
		if strings.HasPrefix(upper, label) {
			upper = strings.TrimSpace(strings.TrimPrefix(upper, label))
			break
		}
	}
	if upper == "" {
		return 0, nil
	}

	sign := 1
	switch upper[0] {
	case '+':
		upper = upper[1:]
	case '-':
		sign = -1
		upper = upper[1:]
	default:
		return 0, fmt.Errorf("tzfmt: offset %q is missing a +/- sign", raw)
	}
	upper = strings.TrimSpace(upper)
	if upper == "" {
		return 0, fmt.Errorf("tzfmt: offset %q has no digits after the sign", raw)
	}

	hours, minutes, err := splitHoursMinutes(upper)
	if err != nil {
		return 0, fmt.Errorf("tzfmt: offset %q: %w", raw, err)
	}
	if minutes < 0 || minutes > 59 {
		return 0, fmt.Errorf("tzfmt: offset %q has an invalid minute component %d", raw, minutes)
	}

	total := sign * (hours*60 + minutes)
	if total < minOffsetMinutes || total > maxOffsetMinutes {
		return 0, fmt.Errorf("tzfmt: offset %q (%d minutes) is outside the -12:00..+14:00 range", raw, total)
	}
	return total, nil
}

// splitHoursMinutes handles the digit-only part of an offset, after the
// sign and any UTC/GMT/Z label have been stripped off. It supports a
// colon separator ("5:30") as well as the compact forms people paste
// from ISO 8601 strings ("5", "05", "530", "0530").
func splitHoursMinutes(digits string) (hours, minutes int, err error) {
	// strconv.Atoi accepts its own leading +/- sign, so without this
	// check a second sign character left over from a malformed input
	// like "++5" or "+-5" would sneak through as a valid-looking digit
	// string instead of being rejected here.
	for i := 0; i < len(digits); i++ {
		if digits[i] != ':' && !isDigit(digits[i]) {
			return 0, 0, fmt.Errorf("unrecognized offset digits %q", digits)
		}
	}

	if idx := strings.IndexByte(digits, ':'); idx >= 0 {
		hours, err = strconv.Atoi(digits[:idx])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour %q", digits[:idx])
		}
		minutes, err = strconv.Atoi(digits[idx+1:])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute %q", digits[idx+1:])
		}
		return hours, minutes, nil
	}

	switch len(digits) {
	case 1, 2:
		hours, err = strconv.Atoi(digits)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour %q", digits)
		}
		return hours, 0, nil
	case 3:
		hours, err = strconv.Atoi(digits[:1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour %q", digits[:1])
		}
		minutes, err = strconv.Atoi(digits[1:])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute %q", digits[1:])
		}
		return hours, minutes, nil
	case 4:
		hours, err = strconv.Atoi(digits[:2])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid hour %q", digits[:2])
		}
		minutes, err = strconv.Atoi(digits[2:])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid minute %q", digits[2:])
		}
		return hours, minutes, nil
	default:
		return 0, 0, fmt.Errorf("unrecognized offset digits %q", digits)
	}
}

// FormatOffsetMinutes turns a signed minute count into canonical
// "+HH:MM" / "-HH:MM" form. It is the inverse of OffsetMinutes and is
// exported separately so callers who already have a minute count (say,
// from time.Time.Zone()) can format it without round-tripping through a
// string first.
func FormatOffsetMinutes(total int) (string, error) {
	if total < minOffsetMinutes || total > maxOffsetMinutes {
		return "", fmt.Errorf("tzfmt: %d minutes is outside the -12:00..+14:00 range", total)
	}
	sign := "+"
	if total < 0 {
		sign = "-"
		total = -total
	}
	return fmt.Sprintf("%s%02d:%02d", sign, total/60, total%60), nil
}
