package taps2beats

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMarshallJSON(t *testing.T) {
	expected := `{"BPM":114,"offset":"0.316","beats":[{"at":"0.316","mean":"0.000","variance":"0.000","taps":[]},{"at":"4.524","mean":"4.524","variance":"0.024","taps":["4.570","4.506","4.530","4.530","4.518","4.495","4.529","4.524","4.518","4.518"]},{"at":"8.733","mean":"0.000","variance":"0.000","taps":[]}]}`

	beats := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			{At: Seconds(0.316)},
			{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.024), Taps: seconds(bins[0]...)},
			{At: Seconds(8.733)},
		},
	}

	bytes, err := json.Marshal(beats)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	json := string(bytes)

	if json != expected {
		t.Errorf("JSON marshal error:\n    expected: %s\n    got:      %s\n", expected, json)
	}
}

func TestUnmarshallJSON(t *testing.T) {
	bytes := []byte(`{"BPM":114,"offset":"0.316","beats":[{"at":"0.316","mean":"0.000","variance":"0.000","taps":[]},{"at":"4.524","mean":"4.524","variance":"0.024","taps":["4.570","4.506","4.530","4.530","4.518","4.495","4.529","4.524","4.518","4.518"]},{"at":"8.733","mean":"0.000","variance":"0.000","taps":[]}]}`)

	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats: []Beat{
			{At: Seconds(0.316)},
			{At: Seconds(4.523694381), Mean: Seconds(4.523694381), Variance: Seconds(0.024), Taps: seconds(bins[0]...)},
			{At: Seconds(8.733)},
		},
	}

	beats := Beats{}
	err := json.Unmarshal(bytes, &beats)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("JSON unmarshalling error:\n    expected: %+v\n    got:      %+v\n", expected, beats)
	}

	if beats.Offset != expected.Offset {
		t.Errorf("JSON unmarshalling error:\n    expected: %+v\n    got:      %+v\n", expected, beats)
	}

	compare(beats.Beats, expected.Beats, t)
}
