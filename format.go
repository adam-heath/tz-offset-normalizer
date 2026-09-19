package tzfmt

import "time"

// Format renders t using layout for the date and time portion, then
// appends the UTC offset in canonical "+HH:MM" form.
//
// This exists because Go's own zone verbs are inconsistent about the
// shape of the offset: "Z0700" prints a bare "Z" for UTC but "+0530"
// everywhere else, "-0700" never has a colon, and "MST" prints an
// abbreviation instead of an offset at all. Format sidesteps all of
// that by handling the offset itself, so callers who want one
// predictable suffix (say, for a value that gets string-compared or
// re-parsed downstream) don't have to normalize Go's output afterward.
//
// layout should describe only the date and time portion of the
// timestamp; it must not contain any of Go's zone verbs (Z0700,
// Z07:00, -0700, -07:00, MST), since Format appends the offset on its
// own.
func Format(t time.Time, layout string) (string, error) {
	_, offsetSeconds := t.Zone()
	offset, err := FormatOffsetMinutes(offsetSeconds / 60)
	if err != nil {
		return "", err
	}
	return t.Format(layout) + offset, nil
}
