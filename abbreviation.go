package tzfmt

import "strings"

// commonAbbreviations maps timezone abbreviations to a UTC offset, but
// only for abbreviations that mean the same thing everywhere they're
// used. Plenty of real-world abbreviations (IST, CST, BST) are
// ambiguous across countries and are deliberately left out rather than
// guessed at; a wrong guess is worse than a clear "unknown".
var commonAbbreviations = map[string]string{
	"UTC":  "+00:00",
	"GMT":  "+00:00",
	"EST":  "-05:00",
	"EDT":  "-04:00",
	"CST":  "-06:00",
	"CDT":  "-05:00",
	"MST":  "-07:00",
	"MDT":  "-06:00",
	"PST":  "-08:00",
	"PDT":  "-07:00",
	"JST":  "+09:00",
	"KST":  "+09:00",
	"CET":  "+01:00",
	"CEST": "+02:00",
}

// NormalizeAbbreviation looks up a timezone abbreviation and returns its
// UTC offset in canonical "+HH:MM" form. The bool is false when the
// abbreviation is unknown or ambiguous (see commonAbbreviations above),
// in which case the caller should ask for an explicit offset instead of
// guessing.
func NormalizeAbbreviation(abbr string) (string, bool) {
	key := strings.ToUpper(strings.TrimSpace(abbr))
	offset, ok := commonAbbreviations[key]
	return offset, ok
}
