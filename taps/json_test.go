package taps

import (
	"testing"
	"time"
)

func TestMarshallJSON(t *testing.T) {
	expected := `{"BPM":114,"offset":"0.316","beats":[{"at":"0.316","mean":"0.000","variance":"0.000","taps":[]},{"at":"4.524","mean":"4.524","variance":"4.524","taps":["4.570","4.506","4.530","4.530","4.518","4.495","4.529","4.524","4.518","4.518"]},{"at":"8.733","mean":"0.000","variance":"0.000","taps":[]}]}`

	beats := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{beats[0], beats[8], beats[16]},
	}

	bytes, err := beats.MarshalJSON()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	json := string(bytes)

	if json != expected {
		t.Errorf("Invalid JSON:\n    expected: %s\n    got:      %s\n", expected, json)
	}
}
