## IN PROGRESS

- [ ] Estimate BPM and offset
- [ ] --range [start,end]
- [ ] --no-quantize
- [ ] --cluster-only
- [ ] --shift to adjust so that first beat is 0.0
- [ ] Format output nicely
- [ ] usage()
- [ ] --help
- [ ] JSON output
- [ ] JSON input
- [ ] (TEMPORARILY) copy ckmeans implementation into this project
- [ ] Web page
- [ ] README

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
