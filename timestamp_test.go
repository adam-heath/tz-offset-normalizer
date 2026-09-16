package tzfmt

import "testing"

func TestExtractOffset(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"2023-06-01T12:00:00Z", "+00:00"},
		{"2023-06-01T12:00:00z", "+00:00"},
		{"2023-06-01T12:00:00+05:30", "+05:30"},
		{"2023-06-01T12:00:00-07:00", "-07:00"},
		{"2023-06-01T12:00:00.123456789-07:00", "-07:00"},
		{"2023-06-01T12:00:00+0530", "+05:30"},
		{"Mon, 02 Jan 2006 15:04:05 MST", "-07:00"},
		{"Mon, 02 Jan 2006 15:04:05 GMT", "+00:00"},
		{"Mon, 02 Jan 2006 15:04:05 -0700", "-07:00"},
		{"10/Oct/2000:13:55:36 -0700", "-07:00"},
		{"2023-06-01 12:00:00+05:30", "+05:30"},
		{"  2023-06-01T12:00:00Z  ", "+00:00"},
	}
	for _, c := range cases {
		got, err := ExtractOffset(c.in)
		if err != nil {
			t.Errorf("ExtractOffset(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ExtractOffset(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestExtractOffsetErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"2023-06-01T12:00:00",           // no offset at all
		"2023-06-01",                    // plain date, no time
		"Mon, 02 Jan 2006 15:04:05 IST", // ambiguous abbreviation
		"Mon, 02 Jan 2006 15:04:05",     // no zone
	}
	for _, in := range cases {
		if got, err := ExtractOffset(in); err == nil {
			t.Errorf("ExtractOffset(%q) = %q, want an error", in, got)
		}
	}
}
