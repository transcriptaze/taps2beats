// +build !js !wasm

// Command line utility for the taps2beats module.
//
//   Usage:
//
//   taps2beats [--verbose] [--out <file>] [--interval <interval>] [--quantize] [--forgetting <factor>] [--precision <time>] [--latency <time>] [--shift] <file>
//
//
//   --verbose              Displays operational information
//
//   --out <file>           Writes the estimated beats to the supplied file
//
//   --interval <interval>  Extrapolates (and interpolates) the beats to extend over
//                          the supplied interval. The interval should be specified as
//                          <start>:<end> where <start> and <end> are in Go time format
//                          e.g. --interval 0.3s:1m10.3s.
//                          A '*' interval (--interval '*') will interpolate the beats
//                          over the interval from the earliest to the latest 'tap'.
//
//   --quantize             Adjusts the estimated beats so that they fit to the estimated BPM
//
//   --forgetting <factor>  Discounts earlier taps from earlier loops as being less accurate
//                          than later loops due to the listener learning the music. e.g. a
//                          factor of 0.1 discounts each loop by 10% over the subsequent one.
//                          A negative factor inverts the weighting i.e. earlier loops are
//                          weighted as more accurate than later loops.
//
//   --precision <time>     Rounds the beats and all times to the specified precision (in Go
//                          time format) e.g. --precision 1ms will round all times to the
//                          nearest millisecond. The default precision is 1ms.
//
//   --latency <time>       Adjusts all times to compensate for the latency between the
//                          actual beat and the detected 'tap' e.g. --latency 73ms
//
//   --clean                Discards outlier taps i.e. taps that are assigned to beats with too few taps.
//
//   --shift                Adjusts all beats (and times) so that the first beat in the
//                          interval falls on 0s.
//
//   --json                 Formats the output as prettified JSON, with all the times converted
//                          to seconds (to a precision of 1ms)
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/transcriptaze/taps2beats/taps2beats"
)

const VERSION = "v0.1.0"

type interval struct {
	set   bool
	start *time.Duration
	end   *time.Duration
}

var options = struct {
	outfile    string
	interval   interval
	quantize   bool
	forgetting float64
	precision  time.Duration
	latency    time.Duration
	clean      bool
	shift      bool
	json       bool
	verbose    bool
	help       bool
}{
	outfile:    "",
	interval:   interval{},
	quantize:   false,
	forgetting: 0.0,
	precision:  1 * time.Millisecond,
	latency:    0 * time.Millisecond,
	clean:      false,
	shift:      false,
	json:       false,
	verbose:    false,
	help:       false,
}

func main() {
	flag.StringVar(&options.outfile, "out", options.outfile, "output file path")
	flag.Var(&options.interval, "interval", "start and end times (in seconds) for which to return beats (e.g. 0.8s:10.0s)")
	flag.BoolVar(&options.quantize, "quantize", options.quantize, "adjusts the tapped beats to fit a least squares fitted BPM")
	flag.Float64Var(&options.forgetting, "forgetting", options.forgetting, "'forgetting factor' for discounting older taps")
	flag.DurationVar(&options.precision, "precision", options.precision, "time precision for returned 'beats', in Go 'time' format (e.g. 1ms)")
	flag.DurationVar(&options.latency, "latency", options.latency, "delay for which to compensate, in Go 'time' format (e.g. 70ms)")
	flag.BoolVar(&options.clean, "clean", options.clean, "discards outlier taps i.e. taps assigned to beats with too few taps")
	flag.BoolVar(&options.shift, "shift", options.shift, "shifts all times so that the first beat is on 0")
	flag.BoolVar(&options.json, "json", options.json, "Sets the output format to prettified JSON")
	flag.BoolVar(&options.verbose, "verbose", options.verbose, "enables verbose progress messages")
	flag.BoolVar(&options.help, "help", options.help, "displays the 'help' information")
	flag.Parse()

	if options.help {
		help()
		os.Exit(0)
	}

	if options.verbose {
		fmt.Printf("\n  taps2beats %s\n\n", VERSION)
	}

	var file string
	var r io.Reader
	var parser = parseTXT
	if len(flag.Args()) == 0 {
		file = "<stdin>"
		r = os.Stdin
	} else {
		file = flag.Args()[0]
		f, err := os.Open(file)
		if err != nil {
			fmt.Printf("\n  ** ERROR: unable to open file %s (%v)\n\n", file, err)
			os.Exit(1)
		}

		defer f.Close()

		r = f

		if strings.HasSuffix(strings.ToLower(file), ".json") {
			parser = parseJSON
		}
	}

	if options.verbose {
		fmt.Printf("  ... reading data from %s\n", file)
	}

	N, data, err := read(r, parser)
	if err != nil {
		fmt.Printf("\n  ** ERROR: unable to read data from %s (%v)\n\n", file, err)
		os.Exit(1)
	} else if N == 0 {
		fmt.Printf("\n  ** ERROR: no data \n\n")
		os.Exit(1)
	}

	if options.verbose {
		fmt.Printf("  ... %v values read from %s\n", N, file)
	}

	if options.verbose {
		fmt.Printf("  ... using forgetting factor %0.1f\n", options.forgetting)
	}

	beats := taps2beats.Taps2Beats(taps2beats.Floats2Seconds(data), options.forgetting)

	// ... sanity check
	if len(beats.Beats) <= 1 || (beats.Variance != nil && *beats.Variance > 0.1) {
		fmt.Printf("\n  ** ERROR: insufficient data\n\n")
		os.Exit(1)
	}

	// ... clean
	if options.clean {
		if options.verbose {
			fmt.Printf("  ... discarding outlier taps\n")
		}

		if b, err := beats.Clean(); err != nil {
			fmt.Printf("\n  ** ERROR: unable to clean beats (%v)\n\n", err)
			os.Exit(1)
		} else {
			beats = b
		}
	}

	// ... quantize
	if options.quantize {
		if options.verbose {
			fmt.Printf("  ... quantizing tapped beats to match estimated BPM\n")
		}

		if err := beats.Quantize(); err != nil {
			fmt.Printf("\n  ** ERROR: unable to quantize beats (%v)\n\n", err)
			os.Exit(1)
		}
	}

	// ... interpolate
	if options.interval.set && len(beats.Beats) > 1 {
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

		if options.verbose {
			fmt.Printf("  ... interpolating missing beats over interval %v..%v \n", start, end)
		}

		if err := beats.Interpolate(start, end); err != nil {
			fmt.Printf("\n  ** ERROR: unable to interpolate beats (%v)\n\n", err)
			os.Exit(1)
		}
	}

	if options.verbose {
		fmt.Printf("  ... %v beats\n", len(beats.Beats))
	}

	if options.latency != 0 {
		if options.verbose {
			fmt.Printf("  ... compensating for %v latency\n", options.latency)
		}

		beats.Sub(options.latency)
	}

	if options.shift {
		if options.verbose {
			fmt.Printf("  ... shifting beats to start at 0\n")
		}

		beats.Sub(beats.Offset)
	}

	// ... round
	if options.verbose {
		fmt.Printf("  ... rounding to %v\n", options.precision)
	}

	beats.Round(options.precision)

	// ... format and print
	var b bytes.Buffer

	if options.json || strings.HasSuffix(strings.ToLower(options.outfile), ".json") {
		if err := formatJSON(beats, &b); err != nil {
			fmt.Printf("\n  ** ERROR: unable to format output as JSON (%v)\n\n", err)
			os.Exit(1)
		}
	} else {
		if err := formatTXT(beats, &b); err != nil {
			fmt.Printf("\n  ** ERROR: unable to format output (%v)\n\n", err)
			os.Exit(1)
		}
	}

	if options.outfile == "" {
		fmt.Println()
		fmt.Printf("%s", string(b.Bytes()))
		fmt.Println()
	} else {
		ioutil.WriteFile(options.outfile, b.Bytes(), 0644)
	}
}

func help() {
	fmt.Println()
	fmt.Printf("  taps2beats %s\n", VERSION)
	fmt.Println()
	fmt.Println("  taps2beats is a utility to estimate the beats of a song from a file containing a list of")
	fmt.Println("  the times at which person (or other musical entity of whatever sort) 'taps' to the beat.")
	fmt.Println()
	fmt.Println("  The 'taps' in the input file are expected to be in seconds, and arranged as lines where each")
	fmt.Println("  line is a single loop of the song being 'tapped' e.g.:")
	fmt.Println()
	fmt.Println("     4.570 5.0635 5.603 6.102 6.642 7.141 7.7106 8.192")
	fmt.Println("     5.045 5.5917 6.114 6.619 7.135 7.693 8.2038")
	fmt.Println("     4.529 5.0576 5.591 6.137 6.630 7.176 7.6989 8.227 9.87353")
	fmt.Println("     ....")
	fmt.Println()
	fmt.Println("  The only requirement is that the taps are separated by whitespace - the lines do not have to")
	fmt.Println("  contain the same number of values, the values do not have to be in time order, nor are they")
	fmt.Println("  required to have the same precision.")
	fmt.Println()
	fmt.Println("  Usage: taps2beats [--interval <interval>] [--quantize] [--forgetting <factor>] [--latency <delay>] [--precision <time>] [--shift] [--out <file>] [--json] [--verbose] <file>")
	fmt.Println()
	fmt.Println("  Arguments:")
	fmt.Println()
	fmt.Println("    file  Path to file containing the whitespace delimited taps to be clustered into beats. Reads")
	fmt.Println("          from <stdin> if the file is not specified.")
	fmt.Println()
	fmt.Println("  Options:")
	fmt.Println()
	fmt.Println("    --interval <interval>  start and end times (in seconds) for which to return beats (e.g. 0.8s:10.0s)")
	fmt.Println("                           If an interval is specified, taps2beats will attempt to interpolate missing")
	fmt.Println("                           beats using a least squares fit which assumes the BPM is more or less constant.")
	fmt.Println("                           (this may not be true unless it was played to a click track). If no interval is")
	fmt.Println("                           specified, taps2beats returns only the beats for which a taps was detected i.e.")
	fmt.Println("                           there is no interpolation for missing beats")
	fmt.Println()
	fmt.Println("    --quantize             linearizes the estimated beats to a least squares fitted BPM")
	fmt.Println("                           Without the 'quantize' option, the beats are estimated to be the mean of the")
	fmt.Println("                           clustered taps for that beat. --quantize adjusts the estimated beats so that")
	fmt.Println("                           they fit a straight line i.e. constant BPM")
	fmt.Println()
	fmt.Println("    --forgetting <factor>  'forgetting factor' for discounting older taps, on the basis that the later")
	fmt.Println("                           taps are probably more accurate since the person is more familiar with the song.")
	fmt.Println("                           The factor is applied on a per-line basis i.e. all the taps in a line are")
	fmt.Println("                           discounted by the same amount. The default factor of 0 weights all taps the same,")
	fmt.Println("                           while a factor of 0.1 cumulatively discounts each subsequent line by 10% (a factor")
	fmt.Println("                           of -0.1 inverts that and discounts later taps rather than earlier taps")
	fmt.Println()
	fmt.Println("    --latency <delay>      delay for which to compensate, in Go 'time' format (e.g. 70ms)")
	fmt.Println()
	fmt.Println("    --precision <time>    time precision for returned 'beats', in Go 'time' format (e.g. 1ms)")
	fmt.Println("    --out                 output file path")
	fmt.Println("    --shift               shifts all times so that the first beat is on 0 and the offset is 0")
	fmt.Println("    --json                formats the output as prettified JSON")
	fmt.Println("    --verbose             enables verbose progress messages")
	fmt.Println("    --help                displays the this information")

	fmt.Println()
}

func choose(beats []taps2beats.Beat, f func(p, q time.Duration) bool) time.Duration {
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
