## IN PROGRESS

- [ ] --latency
- [ ] Estimate BPM and offset
- [ ] Exponential forgetting factor
- [ ] --no-quantize
- [ ] --range [start,end]
- [ ] --shift
- [ ] Format output nicely
- [ ] JSON output
- [ ] JSON input
- [ ] (TEMPORARILY) copy ckmeans implementation into this project
- [ ] Web page

- [x] --precision
- [x] Limit precision to 1ms
- [x] Limit number of interpolation loops to a 'reasonable' number based on max BPM * 16beats/measure (or somesuch)
- [x] Interpolation for pathological data
- [x] Extend beats to min/max
- [x] Interpolate missing beats

## TODO

1. Look at beat detection algorithms
2. Gradient descent for interpolation (?)
