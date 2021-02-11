// +build !js !wasm

package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/transcriptaze/taps2beats/taps2beats"
)

func formatJSON(beats taps2beats.Beats, f io.Writer) error {
	bytes, err := json.MarshalIndent(beats, "", " ")
	if err != nil {
		return err
	}

	f.Write(bytes)

	return nil
}

func formatTXT(beats taps2beats.Beats, f io.Writer) error {
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
