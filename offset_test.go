package tzfmt

import "testing"

func TestNormalizeOffset(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Z", "+00:00"},
		{"z", "+00:00"},
		{"UTC", "+00:00"},
		{"GMT", "+00:00"},
		{"+5", "+05:00"},
		{"-5", "-05:00"},
		{"+05", "+05:00"},
		{"+0530", "+05:30"},
		{"-0530", "-05:30"},
		{"+5:30", "+05:30"},
		{"  +5:30  ", "+05:30"},
		{"UTC+5:30", "+05:30"},
		{"GMT-8", "-08:00"},
		{"utc+05:00", "+05:00"},
		{"+14:00", "+14:00"},
		{"-12:00", "-12:00"},
	}
	for _, c := range cases {
		got, err := NormalizeOffset(c.in)
		if err != nil {
			t.Errorf("NormalizeOffset(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("NormalizeOffset(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalizeOffsetErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"5",       // missing sign
		"+5:99",   // invalid minutes
		"+15:00",  // beyond +14:00
		"-13:00",  // beyond -12:00
		"+5:30:00", // too many parts
		"++5",
	}
	for _, in := range cases {
		if _, err := NormalizeOffset(in); err == nil {
			t.Errorf("NormalizeOffset(%q) expected an error, got none", in)
		}
	}
}

func TestOffsetMinutes(t *testing.T) {
	got, err := OffsetMinutes("-0430")
	if err != nil {
		t.Fatalf("OffsetMinutes returned error: %v", err)
	}
	want := -(4*60 + 30)
	if got != want {
		t.Errorf("OffsetMinutes(-0430) = %d, want %d", got, want)
	}
}

func TestFormatOffsetMinutesRoundTrip(t *testing.T) {
	for _, minutes := range []int{0, 60, -60, 330, -330, 840, -720} {
		s, err := FormatOffsetMinutes(minutes)
		if err != nil {
			t.Fatalf("FormatOffsetMinutes(%d) returned error: %v", minutes, err)
		}
		back, err := OffsetMinutes(s)
		if err != nil {
			t.Fatalf("OffsetMinutes(%q) returned error: %v", s, err)
		}
		if back != minutes {
			t.Errorf("round trip through %q: got %d, want %d", s, back, minutes)
		}
	}
}

func TestNormalizeAbbreviation(t *testing.T) {
	offset, ok := NormalizeAbbreviation(" pst ")
	if !ok {
		t.Fatal("expected PST to be recognized")
	}
	if offset != "-08:00" {
		t.Errorf("NormalizeAbbreviation(PST) = %q, want -08:00", offset)
	}

	if _, ok := NormalizeAbbreviation("IST"); ok {
		t.Error("IST is ambiguous and should not resolve")
	}
}
