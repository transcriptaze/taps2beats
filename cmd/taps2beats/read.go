package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func read(r io.Reader, isJSON bool) (int, [][]float64, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, nil, err
	}

	var data [][]float64
	if isJSON {
		data, err = parseJSON(bytes)
	} else {
		data, err = parseTXT(bytes)
	}

	if err != nil {
		return 0, nil, err
	}

	count := 0
	for _, row := range data {
		count += len(row)
	}

	return count, data, nil
}

func parseJSON(bytes []byte) ([][]float64, error) {
	taps := struct {
		Taps [][]float64 `json:"taps"`
	}{}

	if err := json.Unmarshal(bytes, &taps); err != nil {
		return nil, err
	}

	return taps.Taps, nil
}

func parseTXT(bytes []byte) ([][]float64, error) {
	data := [][]float64{}
	re := regexp.MustCompile(`\s+`)
	for _, line := range strings.Split(string(bytes), "\n") {
		tokens := re.Split(line, -1)
		row := []float64{}
		for _, t := range tokens {
			if t != "" {
				if v, err := strconv.ParseFloat(t, 64); err != nil {
					fmt.Printf("  ** WARN: invalid value (%s)\n", t)
				} else {
					row = append(row, v)
				}
			}
		}

		if len(row) > 0 {
			data = append(data, row)
		}
	}

	return data, nil
}
