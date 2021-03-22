## IN PROGRESS

- [x] Initial version release
- [x] Error if beats == 1 or variance is too high
- [x] Discard outlier beats i.e. beats with too few taps
- [x] godoc examples
- [x] godoc
- [x] Implement Beats.UnmarshalJSON
- [x] Implement Beats.MarshalJSON
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
6. https://towardsdatascience.com/deep-learning-in-geomtry-arclentgh-learning-119d347231ce
7. https://dsp.stackexchange.com/questions/60528/how-to-compute-key-of-a-song
8. Improve BPM estimation (or at least make it a bit more robust)
9. Improve clustering when tapping double time
10. More sanity checks
    - minimum gap between beats (e.g. when data only has one beat but the clustering produces 5)

## NOTES


