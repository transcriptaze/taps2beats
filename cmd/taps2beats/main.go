package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/twystd/taps2beats/taps"
)

const VERSION = "v0.0.0"

var options = struct {
	precision   time.Duration
	latency     time.Duration
	forgetting  float64
	quantize    bool
	interpolate bool
	outfile     string
	debug       bool
}{
	precision:   taps.Default.Precision,
	latency:     taps.Default.Latency,
	forgetting:  taps.Default.Forgetting,
	quantize:    taps.Default.Quantize,
	interpolate: taps.Default.Interpolate,
	outfile:     "",
	debug:       false,
}

func main() {
	flag.DurationVar(&options.precision, "precision", options.precision, "time precision for returned 'beats'")
	flag.DurationVar(&options.latency, "latency", options.latency, "delay for which to compensate")
	flag.Float64Var(&options.forgetting, "forgetting", options.forgetting, "'forgetting factor' for discounting older taps")
	flag.BoolVar(&options.quantize, "quantize", options.quantize, "adjusts the tapped beats to fit a least squares fitted BPM")
	flag.BoolVar(&options.interpolate, "interpolate", options.interpolate, "adds beats in gaps between tapped beats")
	flag.StringVar(&options.outfile, "out", options.outfile, "output file path")
	flag.BoolVar(&options.debug, "debug", options.debug, "enables debugging")
	flag.Parse()

	if options.debug {
		fmt.Printf("\n  taps2beats %s\n\n", VERSION)
	}

	if len(flag.Args()) == 0 {
		usage()
		os.Exit(1)
	}

	file := flag.Args()[0]
	if options.debug {
		fmt.Printf("  ... reading data from %s\n", file)
	}

	data, err := read(file)
	if err != nil {
		fmt.Printf("\n  ** ERROR: unable to read data from file %s (%v)\n\n", file, err)
		os.Exit(1)
	}

	if options.debug {
		fmt.Printf("  ... %v values read from %s\n", len(data), file)
	}

	t2b := taps.T2B{
		Precision:   options.precision,
		Latency:     options.latency,
		Forgetting:  options.forgetting,
		Quantize:    options.quantize,
		Interpolate: options.interpolate,
	}

	if options.debug {
		fmt.Printf("  ... rounding to %v precision\n", t2b.Precision)
		fmt.Printf("  ... compensating for %v latency\n", t2b.Latency)
		fmt.Printf("  ... using forgetting factor %v latency\n", t2b.Forgetting)

		if t2b.Quantize {
			fmt.Printf("  ... quantizing tapped beats to match estimated BPM\n")
		} else {
			fmt.Printf("  ... tapped beats are not quantized to match estimated BPM\n")
		}

		if t2b.Interpolate {
			fmt.Printf("  ... interpolating missing beats\n")
		} else {
			fmt.Printf("  ... ignoring missing beats\n")
		}
	}

	beats, err := t2b.Taps2Beats(taps.Floats2Seconds(data), 0, 8500*time.Millisecond), nil
	if err != nil {
		fmt.Printf("\n  ** ERROR: unable to translate taps to beats (%v)\n\n", err)
		os.Exit(1)
	}

	if options.debug {
		fmt.Printf("  ... %v beats\n", len(beats.Beats))
	}

	var b bytes.Buffer

	print(&b, beats)

	if options.outfile == "" {
		fmt.Println()
		fmt.Printf("%s", string(b.Bytes()))
		fmt.Println()
	} else {
		ioutil.WriteFile(options.outfile, b.Bytes(), 0644)
	}
}

func read(f string) ([][]float64, error) {
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

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

		data = append(data, row)
	}

	return data, nil
}

func print(f io.Writer, beats taps.Beats) {
	if beats.BPM != nil && beats.Offset != nil {
		fmt.Fprintf(f, "BPM:    %v\n", *beats.BPM)
		fmt.Fprintf(f, "Offset: %v\n\n", beats.Offset)
	} else if beats.BPM != nil {
		fmt.Fprintf(f, "BPM: %v\n\n", beats.BPM)
	} else if beats.Offset != nil {
		fmt.Fprintf(f, "Offset: %v\n\n", beats.Offset)
	}

	for i, b := range beats.Beats {
		fmt.Fprintf(f, "%d %v %v %v %v\n", i+1, b.At, b.Mean, b.Variance, b.Taps)
	}
	//	columns := 0
	//	for _, b := range beats {
	//		if len(b.Values) > columns {
	//			columns = len(c.Values)
	//		}
	//	}
	//	columns += 3
	//
	//	table := make([][]string, len(clusters))
	//	for i := range table {
	//		table[i] = make([]string, columns)
	//	}
	//
	//	for i, c := range clusters {
	//		table[i][0] = fmt.Sprintf("%d", i+1)
	//		table[i][1] = fmt.Sprintf("%v", c.Center)
	//		table[i][2] = fmt.Sprintf("%v", c.Variance)
	//		for j, v := range c.Values {
	//			table[i][j+3] = fmt.Sprintf("%v", v)
	//		}
	//	}
	//
	//	widths := make([]int, columns)
	//	for _, row := range table {
	//		for i, s := range row {
	//			if len(s) > widths[i] {
	//				widths[i] = len(s)
	//			}
	//		}
	//	}
	//
	//	formats := make([]string, columns)
	//	for i, w := range widths {
	//		formats[i] = fmt.Sprintf("%%-%dv", w)
	//	}
	//
	//	for i, c := range clusters {
	//		line := ""
	//		line += fmt.Sprintf(formats[0], i+1)
	//		line += "  "
	//		line += fmt.Sprintf(formats[1], c.Center)
	//		line += "  "
	//		line += fmt.Sprintf(formats[2], c.Variance)
	//		line += "  "
	//		for j, v := range c.Values {
	//			line += " "
	//			line += fmt.Sprintf(formats[j+3], v)
	//		}
	//
	//		fmt.Fprintf(f, "%s\n", line)
	//	}
}

func usage() {
	fmt.Println()
	fmt.Println("  Usage: taps2beats [options] <file>")
	fmt.Println()
	fmt.Println("  Arguments:")
	fmt.Println()
	fmt.Println("    file  Path to file containing the whitespace delimited taps to be clustered into beats")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --debug     Displays internal information for diagnosing errors")
	fmt.Println()
}
