package moment

import "testing"

func TestParseZone(t *testing.T) {
	m, err := ParseZone("2017-07-14T02:40:00-05:00")
	if err != nil {
		t.Fatalf("ParseZone error: %v", err)
	}
	if got := m.UTCOffset(); got != -300 {
		t.Errorf("UTCOffset = %d, want -300", got)
	}
	// The offset is preserved in formatting rather than normalised to UTC.
	if got, want := m.Format("HH:mm Z"), "02:40 -05:00"; got != want {
		t.Errorf("Format = %q, want %q", got, want)
	}
	if m.CreationData() == nil || m.CreationData().Format != "parseZone" {
		t.Errorf("CreationData.Format = %+v, want parseZone", m.CreationData())
	}

	// A positive offset.
	p, err := ParseZone("2017-07-14T02:40:00+09:30")
	if err != nil {
		t.Fatalf("ParseZone(+) error: %v", err)
	}
	if got := p.UTCOffset(); got != 570 {
		t.Errorf("UTCOffset(+09:30) = %d, want 570", got)
	}

	if _, err := ParseZone("not a date"); err == nil {
		t.Errorf("ParseZone should reject junk")
	}
}
