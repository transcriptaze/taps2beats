package taps2beats

import (
	"math"
	"testing"
	"time"
)

func TestQuantize(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{quantized[8], quantized[9], quantized[10], quantized[11], quantized[12], quantized[13], quantized[14], quantized[15]},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[10], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	err := beats.Quantize()
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestQuantizeWithNoData(t *testing.T) {
	expected := Beats{
		BPM:    0,
		Offset: 0 * time.Millisecond,
		Beats:  []Beat{},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats:  []Beat{},
	}

	err := beats.Quantize()
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if beats.Offset != expected.Offset {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestQuantizeWithOneDataPoint(t *testing.T) {
	expected := Beats{
		BPM:    0,
		Offset: 4524 * time.Millisecond,
		Beats:  []Beat{beats[8]},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats:  []Beat{beats[8]},
	}

	err := beats.Quantize()
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) >= 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}

func TestQuantizeWithMissingBeat(t *testing.T) {
	expected := Beats{
		BPM:    114,
		Offset: 316 * time.Millisecond,
		Beats:  []Beat{quantized[8], quantized[9], quantized[11], quantized[12], quantized[13], quantized[14], quantized[15]},
	}

	beats := Beats{
		BPM:    123,
		Offset: 117 * time.Millisecond,
		Beats:  []Beat{beats[8], beats[9], beats[11], beats[12], beats[13], beats[14], beats[15]},
	}

	err := beats.Quantize()
	if err != nil {
		t.Fatalf("Unexpected error (%v)", err)
	}

	if beats.BPM != expected.BPM {
		t.Errorf("BPM unexpectedly modified - expected:%v, got:%v", expected.BPM, beats.BPM)
	}

	if math.Abs(beats.Offset.Seconds()-expected.Offset.Seconds()) > 0.0011 {
		t.Errorf("Incorrect offset - expected:%v, got:%v", expected.Offset, beats.Offset)
	}

	compare(beats.Beats, expected.Beats, t)
}
