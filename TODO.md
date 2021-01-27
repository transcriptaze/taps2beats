## IN PROGRESS

- [ ] Web page
      - [ ] Fix slider thumbs so that they actually point to the timeline
      - [ ] 'TAP' pad
      - [ ] Explanatory text
      - [ ] 'Loading' windmill
      - [ ] Fade out loading overlay
      - [ ] Show taps
      - [ ] Load audio from file
      - [x] Panels a la transcriptase
      - [x] Loop controls
      - [x] FontAwesome credit in README
      - [x] Start/end controls
      - [x] Parse YouTube URL
      - [x] Splashscreen

- [ ] Initial version release

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

- [ ] Better range slider
      https://css-tricks.com/multi-thumb-sliders-general-case

1. Look into beat detection algorithms
2. Look into gradient descent for interpolation
3. Look into constrained Deming regression for interpolation
3. https://moultano.wordpress.com/2018/11/08/minhashing-3kbzhsxyg4467-6
4. --format option to set the output value format (?)
5. --interactive, to record the 'taps' directly (?)
6. https://gitlab.com/agrahn/ABLoopPlayer
7. https://agrahn.gitlab.io/ABLoopPlayer
8. https://towardsdatascience.com/deep-learning-in-geomtry-arclentgh-learning-119d347231ce

## NOTES

1. https://w3c.github.io/aria-practices/examples/slider/multithumb-slider.html


