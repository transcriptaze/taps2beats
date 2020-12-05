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

type interval struct {
	set   bool
	start *time.Duration
	end   *time.Duration
}

var options = struct {
	precision  time.Duration
	latency    time.Duration
	forgetting float64
	quantize   bool
	interval   interval
	shift      bool
	outfile    string
	debug      bool
}{
	precision:  taps.Default.Precision,
	latency:    taps.Default.Latency,
	forgetting: taps.Default.Forgetting,
	quantize:   false,
	interval:   interval{},
	outfile:    "",
	debug:      false,
}

func main() {
	flag.DurationVar(&options.precision, "precision", options.precision, "time precision for returned 'beats'")
	flag.DurationVar(&options.latency, "latency", options.latency, "delay for which to compensate")
	flag.Float64Var(&options.forgetting, "forgetting", options.forgetting, "'forgetting factor' for discounting older taps")
	flag.BoolVar(&options.quantize, "quantize", options.quantize, "adjusts the tapped beats to fit a least squares fitted BPM")
	flag.Var(&options.interval, "interval", "start and end times (in seconds) for which to return beats e.g. 0.8:10.0")
	flag.BoolVar(&options.shift, "shift", options.shift, "shifts all times so that the first beat is on 0")
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

	N, data, err := read(file)
	if err != nil {
		fmt.Printf("\n  ** ERROR: unable to read data from file %s (%v)\n\n", file, err)
		os.Exit(1)
	} else if N == 0 {
		fmt.Printf("\n  ** ERROR: no data in file %s (%v)\n\n", file, err)
		os.Exit(1)
	}

	if options.debug {
		fmt.Printf("  ... %v values read from %s\n", len(data), file)
	}

	t2b := taps.T2B{
		Precision:  options.precision,
		Latency:    options.latency,
		Forgetting: options.forgetting,
	}

	if options.debug {
		fmt.Printf("  ... rounding to %v precision\n", t2b.Precision)
		fmt.Printf("  ... compensating for %v latency\n", t2b.Latency)
		fmt.Printf("  ... using forgetting factor %v latency\n", t2b.Forgetting)
	}

	beats := t2b.Taps2Beats(taps.Floats2Seconds(data))

	// ... quantize
	if options.quantize {
		if options.debug {
			fmt.Printf("  ... quantizing tapped beats to match estimated BPM\n")
		}

		beats, err = t2b.Quantize(beats)
		if err != nil {
			fmt.Printf("\n  ** ERROR: unable to quantize beats (%v)\n\n", err)
			os.Exit(1)
		}
	} else if options.debug {
		fmt.Printf("  ... tapped beats are not quantized to match estimated BPM\n")
	}

	// ... interpolate
	if options.interval.set && len(beats.Beats) > 0 {
		var start time.Duration
		var end time.Duration

		if options.interval.start == nil {
			start = choose(beats.Beats, func(p, q time.Duration) bool { return p < q })
		} else {
			start = *options.interval.start
		}

		if options.interval.end == nil {
			end = choose(beats.Beats, func(p, q time.Duration) bool { return p > q })
		} else {
			end = *options.interval.end
		}

		if options.debug {
			fmt.Printf("  ... interpolating missing beats over interval %v..%v \n", start, end)
		}

		beats, err = t2b.Interpolate(beats, start, end)
		if err != nil {
			fmt.Printf("\n  ** ERROR: unable to interpolate beats (%v)\n\n", err)
			os.Exit(1)
		}

	} else if options.debug && len(beats.Beats) > 0 {
		fmt.Printf("  ... ignoring missing beats\n")
	}

	if options.debug {
		fmt.Printf("  ... %v beats\n", len(beats.Beats))
	}

	if options.shift {
		if options.debug {
			fmt.Printf("  ... shifting beats to start at 0\n")
		}

		beats = t2b.Shift(beats)
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

func read(f string) (int, [][]float64, error) {
	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		return 0, nil, err
	}

	count := 0
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
					count++
				}
			}
		}

		data = append(data, row)
	}

	return count, data, nil
}

func print(f io.Writer, beats taps.Beats) {
	fmt.Fprintf(f, "BPM:    %v\n", beats.BPM)
	fmt.Fprintf(f, "Offset: %v\n\n", beats.Offset)

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

func choose(beats []taps.Beat, f func(p, q time.Duration) bool) time.Duration {
	if len(beats) < 1 {
		panic("Insufficient data")
	}

	v := beats[0].At
	for _, b := range beats {
		if f(b.At, v) {
			v = b.At
		}

		for _, t := range b.Taps {
			if f(t, v) {
				v = t
			}
		}
	}
	return v
}

func (v *interval) String() string {
	if v.start != nil && v.end != nil {
		return fmt.Sprintf("%v:%v", v.start, v.end)
	} else if v.start != nil {
		return fmt.Sprintf("%v", v.start)
	} else if v.end != nil {
		return fmt.Sprintf(":%v", v.end)
	}

	return "*"
}

func (v *interval) Set(s string) error {
	re := regexp.MustCompile(`[0-9]+(\.[0-9]*)?`)
	tokens := strings.Split(s, ":")

	v.set = true

	if len(tokens) > 1 {
		if re.MatchString(tokens[0]) {
			start, err := time.ParseDuration(tokens[0])
			if err != nil {
				return err
			}

			v.start = &start
		}

		if re.MatchString(tokens[1]) {
			end, err := time.ParseDuration(tokens[1])
			if err != nil {
				return err
			}

			if v.start == nil || end > *v.start {
				v.end = &end
			}
		}

		return nil
	}

	if len(tokens) > 0 {
		if re.MatchString(tokens[0]) {
			start, err := time.ParseDuration(tokens[0])
			if err != nil {
				return err
			}

			v.start = &start

			return nil
		}
	}

	return nil
}
