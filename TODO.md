## IN PROGRESS

- [ ] commonalize taps.interpolate and taps.remap
- [ ] (TEMPORARILY) copy ckmeans implementation into this project
- [ ] Format output nicely
- [ ] JSON output
- [ ] JSON input
- [ ] Web page
- [ ] README
- [ ] Read from stdin
- [ ] godoc

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

1. Look at beat detection algorithms
2. Gradient descent for interpolation (?)
3. https://moultano.wordpress.com/2018/11/08/minhashing-3kbzhsxyg4467-6
