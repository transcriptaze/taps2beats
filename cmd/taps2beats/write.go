package main

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/twystd/taps2beats/taps"
)

type instant time.Duration

func (t instant) MarshalText() ([]byte, error) {
	s := fmt.Sprintf("%.3f", time.Duration(t).Seconds())

	return []byte(s), nil
}

func formatJSON(beats taps.Beats, f io.Writer) error {
	type beat struct {
		At       instant   `json:"at"`
		Mean     instant   `json:"mean"`
		Variance instant   `json:"variance"`
		Taps     []instant `json:"taps"`
	}

	b := struct {
		BPM    uint    `json:"BPM"`
		Offset instant `json:"offset"`
		Beats  []beat  `json:"beats"`
	}{
		BPM:    beats.BPM,
		Offset: instant(beats.Offset),
		Beats:  make([]beat, len(beats.Beats)),
	}

	for i, bb := range beats.Beats {
		b.Beats[i] = beat{
			At:       instant(bb.At),
			Mean:     instant(bb.Mean),
			Variance: instant(bb.Mean),
			Taps:     make([]instant, len(bb.Taps)),
		}

		for j, t := range bb.Taps {
			b.Beats[i].Taps[j] = instant(t)
		}
	}

	bytes, err := json.MarshalIndent(beats, "", " ")
	if err != nil {
		return err
	}

	f.Write(bytes)

	return nil
}

func formatTXT(beats taps.Beats, f io.Writer) error {
	grid := [][]string{}
	for i, b := range beats.Beats {
		row := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%v", b.At),
		}

		if len(b.Taps) > 0 {
			row = append(row, fmt.Sprintf("%v", b.Mean))
			row = append(row, fmt.Sprintf("%v", b.Variance))
			for _, t := range b.Taps {
				row = append(row, fmt.Sprintf("%v", t))
			}
		}

		grid = append(grid, row)
	}

	columns := 0
	for _, row := range grid {
		if len(row) > columns {
			columns = len(row)
		}
	}

	cols := make([]int, columns)
	for _, row := range grid {
		for i, v := range row {
			if len(v) > cols[i] {
				cols[i] = len(v)
			}
		}
	}

	fmt.Fprintf(f, "BPM:    %v\n", beats.BPM)
	fmt.Fprintf(f, "Offset: %v\n\n", beats.Offset)
	for _, row := range grid {
		fmt.Fprintf(f, "%-*s", cols[0], row[0])
		for i, v := range row[1:] {
			fmt.Fprintf(f, " %-*s", cols[i+1], v)
		}
		fmt.Fprintln(f)
	}

	return nil
}
