## IN PROGRESS

- [ ] Implement Beats.MarshalJSON/UnmarshalJSON
- [ ] Implement Beats.MarshalText/UnmarshalText
- [ ] godoc
- [ ] Initial version release
- [ ] Web page

- [x] Commonalize BPM + offset
- [x] Make github repo public
- [x] Read from stdin
- [x] Commonalise regression
- [x] JSON input
- [x] JSON output
- [x] README
- [x] Copy ckmeans implementation into this project
- [x] Factor latency out as callable function Sub
- [x] Format output nicely
- [x] commonalize taps.interpolate and taps.remap
- [x] test Quantize with missing beat (shouldn't interpolate!)
- [x] usage()
- [x] --help
- [x] split interpolate out as callable functions and let taps2beats just do clustering
- [x] split quantize out as callable functions and let taps2beats just do clustering
- [x] make BPM and Offset not pointers
- [x] --shift to adjust so that first beat is 0.0
- [x] --range [start,end]
- [x] --interpolate
- [x] --quantize
- [x] Estimate BPM and offset
- [x] Exponential forgetting factor
- [x] --latency
- [x] --precision
- [x] Limit precision to 1ms
- [x] Limit number of interpolation loops to a 'reasonable' number based on max BPM * 16beats/measure (or somesuch)
- [x] Interpolation for pathological data
- [x] Extend beats to min/max
- [x] Interpolate missing beats

## TODO

1. Look into beat detection algorithms
2. Look into gradient descent for interpolation
3. Look into constrained Deming regression for interpolation
3. https://moultano.wordpress.com/2018/11/08/minhashing-3kbzhsxyg4467-6
4. --format option to set the output value format (?)
5. --interactive, to record the 'taps' directly (?)
