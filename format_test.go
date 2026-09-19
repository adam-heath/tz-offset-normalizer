package tzfmt

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name string
		when time.Time
		want string
	}{
		{
			name: "UTC",
			when: time.Date(2023, 6, 1, 12, 0, 0, 0, time.UTC),
			want: "2023-06-01T12:00:00+00:00",
		},
		{
			name: "positive fixed offset",
			when: time.Date(2023, 6, 1, 12, 0, 0, 0, time.FixedZone("IST", 5*3600+30*60)),
			want: "2023-06-01T12:00:00+05:30",
		},
		{
			name: "negative fixed offset",
			when: time.Date(2023, 6, 1, 12, 0, 0, 0, time.FixedZone("PDT", -7*3600)),
			want: "2023-06-01T12:00:00-07:00",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := Format(c.when, "2006-01-02T15:04:05")
			if err != nil {
				t.Fatalf("Format returned error: %v", err)
			}
			if got != c.want {
				t.Errorf("Format() = %q, want %q", got, c.want)
			}
		})
	}
}

func TestFormatOutOfRangeOffset(t *testing.T) {
	when := time.Date(2023, 6, 1, 12, 0, 0, 0, time.FixedZone("bogus", 15*3600))
	if _, err := Format(when, "2006-01-02T15:04:05"); err == nil {
		t.Error("Format with an out-of-range offset should return an error")
	}
}
