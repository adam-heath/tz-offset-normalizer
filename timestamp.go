package tzfmt

import (
	"fmt"
	"strings"
)

// ExtractOffset locates the UTC offset embedded at the end of a full
// timestamp string and returns it in canonical "+HH:MM" form. It
// recognizes the offset shapes produced by RFC 3339 ("...05Z" or
// "...05+07:00", with or without fractional seconds), RFC 1123 / RFC
// 1123Z ("... 15:04:05 MST" or "... 15:04:05 -0700"), and the
// Apache/Nginx common log format ("...:36 -0700").
//
// ExtractOffset does not parse the date or time portion of the string at
// all - it only looks at what trails it. Callers who need a full
// time.Time should use time.Parse with a matching layout; this is for
// pulling a normalizable offset out of a timestamp whose exact layout
// isn't known ahead of time (mixed log sources, user-pasted timestamps,
// etc).
func ExtractOffset(timestamp string) (string, error) {
	s := strings.TrimSpace(timestamp)
	if s == "" {
		return "", fmt.Errorf("tzfmt: empty timestamp")
	}

	if off, ok := trailingZ(s); ok {
		return off, nil
	}
	if off, ok := trailingAttachedOffset(s); ok {
		return NormalizeOffset(off)
	}
	if tok, ok := trailingSpacedToken(s); ok {
		if isSignedDigits(tok) {
			return NormalizeOffset(tok)
		}
		if normalized, ok := NormalizeAbbreviation(tok); ok {
			return normalized, nil
		}
		return "", fmt.Errorf("tzfmt: %q is an unrecognized or ambiguous zone abbreviation", tok)
	}

	return "", fmt.Errorf("tzfmt: no recognizable offset at the end of %q", timestamp)
}

// trailingZ matches the RFC 3339 "Z" suffix, e.g. "...12:00:00Z". It
// requires a digit right before the Z so that it doesn't fire on
// something like a bare word ending in a "z".
func trailingZ(s string) (string, bool) {
	n := len(s)
	if n < 2 {
		return "", false
	}
	last := s[n-1]
	if last != 'Z' && last != 'z' {
		return "", false
	}
	if !isDigit(s[n-2]) {
		return "", false
	}
	return "+00:00", true
}

// trailingAttachedOffset matches a numeric offset with no separator
// before it, the way RFC 3339 appends one directly to the seconds (or
// fractional seconds) field: "...05+07:00" or "...05-0700". It only
// looks at the last five or six characters, so it's naturally immune to
// misreading the rest of the timestamp.
func trailingAttachedOffset(s string) (string, bool) {
	n := len(s)
	if n >= 6 {
		sign := s[n-6]
		if (sign == '+' || sign == '-') && s[n-3] == ':' &&
			isDigit(s[n-5]) && isDigit(s[n-4]) && isDigit(s[n-2]) && isDigit(s[n-1]) {
			return s[n-6:], true
		}
	}
	if n >= 5 {
		sign := s[n-5]
		if (sign == '+' || sign == '-') &&
			isDigit(s[n-4]) && isDigit(s[n-3]) && isDigit(s[n-2]) && isDigit(s[n-1]) {
			return s[n-5:], true
		}
	}
	return "", false
}

// trailingSpacedToken returns the final space-separated token of s, the
// shape a zone abbreviation or a signed numeric offset takes in RFC 1123
// ("... 15:04:05 MST"), RFC 1123Z ("... 15:04:05 -0700"), and common log
// format ("...:36 -0700").
func trailingSpacedToken(s string) (string, bool) {
	idx := strings.LastIndexByte(s, ' ')
	if idx < 0 || idx == len(s)-1 {
		return "", false
	}
	return s[idx+1:], true
}

// isSignedDigits reports whether s looks like a signed numeric offset
// ("-0700", "+05:30") rather than a zone abbreviation ("MST").
func isSignedDigits(s string) bool {
	if len(s) < 2 {
		return false
	}
	if s[0] != '+' && s[0] != '-' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isDigit(s[i]) && s[i] != ':' {
			return false
		}
	}
	return true
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
