### taps2beats

`taps2beats` is a somewhat _(?)_ oddball command-line utility and Go module that estimates the beats of a piece
of music from a file of the beats as _tapped_ by a person (or other musical entity of whatever sort). 

The internal algorithm uses an implementation of _Ckmeans.1d.dp_ to cluster the supplied _'taps'_ into the optimal 
equivalent beats, followed by least squares regression to estimate the BPM and offset and (optionally) quantize 
and (optionally) interpolate the beats over the interval.

#### Requirements

- `Go` version 1.15+ 

#### Installation

To install the stable version, download the `tar.gz` or `zip` release for your platform from the [Releases](https://github.com/twystd/taps2beats/releases) 
page, and unpack into the directory of your choice.  

The build the development version, clone the [taps2beats](https://github.com/twystd/taps2beats) _github_ repository:

    git clone https://github.com/twystd/taps2beats.git
    cd taps2beats
    make build


#### Usage

`taps2beats [options] <file>`

The `file` is a list of lines of the 'taps' (in seconds), separated by whitespace e.g.
```
4.57027 5.063594027 5.6035973 6.1026998 6.642708943 7.14179696 7.7107 8.1470916
4.52956007 5.069284016 5.603428973 6.102591998 6.613455 7.14764 7.69912088  8.215609871
4.517865093 5.022782107 5.580101018 6.096715 6.654118921 7.1763719 7.681405914 8.215537871
5.1392891 5.545395086 6.067721066 6.578564068 7.1300991 7.65971 8.134273
4.4981138 5.040234073 5.562732052 6.079333043 6.624973977 7.14165 7.66 8.19905
4.511093 5.069403016 5.586174007 6.108568986 6.57864 7.147523 7.681614 8.262078
```

There is no requirement that the number of taps on each line be the same or that the taps are specified to the
same precision. Ideally, each line represents a _loop_ of the music being tapped, but this is only significant
if a _forgetting_ factor (described below) is used to weight later 'taps' as being more accurate than e.g. the
first few attempts.

If the input filename ends with '.json', the file is parsed as a JSON object that is expected to contain:
```
{ 
  "taps": [][]float 
}
```

Invoking `taps2beats` without an input file reads the _taps_ from stdin.

The output format is a fixed column width list of beats, with each beat represented by a line that contains

1. Beat number
2. Estimated beat 
3. Average of the 'taps' for the beat
4. Variance of the 'taps' for the beat
5. The 'taps' from the input file used to estimate the beat

| Beat | At | Mean | Variance | Taps |
|------|----|------|----------|------|
| ...  |
| 6    | 7.153s | 7.153s | 1.3ms | 7.142s 7.136s 7.177s 7.148s 7.176s 7.13s 7.13s  7.165s 7.148s |
| 7    | 7.686s | 7.686s | 2.1ms | 7.711s 7.693s 7.699s 7.681s 7.652s 7.722s 7.653s 7.687s 7.682s |
| ...  |

Options:

`taps2beats [--verbose] [--out <file>] [--interval <interval>] [--quantize] [--forgetting <factor>] [--precision <time>] [--latency <time>] [--shift] <file>`

```
--verbose              Displays operational information

--out <file>           Writes the estimated beats to the supplied file

--interval <interval>  Extrapolates (and interpolates) the beats to extend over
                       the supplied interval. The interval should be specified as 
                       <start>:<end> where <start> and <end> are in Go time format
                       e.g. --interval 0.3s:1m10.3s.
                       A '*' interval (--interval '*') will interpolate the beats
                       over the interval from the earliest to the latest 'tap'.

--quantize             Adjusts the estimated beats so that they fit to the estimated BPM 

--forgetting <factor>  Discounts earlier taps from earlier loops as being less accurate
                       than later loops due to the listener learning the music. e.g. a
                       factor of 0.1 discounts each loop by 10% over the subsequent one.
                       A negative factor inverts the weighting i.e. earlier loops are
                       weighted as more accurate than later loops.

--precision <time>     Rounds the beats and all times to the specified precision (in Go
                       time format) e.g. --precision 1ms will round all times to the 
                       nearest millisecond. The default precision is 1ms.
   
--latency <time>       Adjusts all times to compensate for the latency between the 
                       actual beat and the detected 'tap' e.g. --latency 73ms
                       
--shift                Adjusts all beats (and times) so that the first beat in the 
                       interval falls on 0s.
                       
--json                 Formats the output as prettified JSON, with all the times converted
                       to seconds (to a precision of 1ms)
```

#### Examples

`taps2beats examples/taps.txt`

```
BPM:    114
Offset: 316ms

1 4.524s 4.524s 0s  4.57s  4.506s 4.53s  4.53s  4.518s 4.495s 4.529s 4.524s 4.518s 4.518s
2 5.058s 5.058s 1ms 5.064s 5.046s 5.058s 5.069s 5.023s 5.133s 5.04s  5.04s  5.046s 5.046s 5.069s
3 5.578s 5.578s 0s  5.604s 5.592s 5.592s 5.603s 5.58s  5.545s 5.563s 5.557s 5.586s 5.551s 5.586s
4 6.1s   6.1s   0s  6.103s 6.114s 6.137s 6.103s 6.097s 6.068s 6.079s 6.132s 6.091s 6.074s 6.109s
5 6.618s 6.618s 1ms 6.643s 6.619s 6.631s 6.613s 6.654s 6.579s 6.625s 6.654s 6.596s 6.608s 6.579s
6 7.153s 7.153s 0s  7.142s 7.136s 7.177s 7.148s 7.176s 7.13s  7.142s 7.194s 7.13s  7.165s 7.148s
7 7.686s 7.686s 1ms 7.711s 7.693s 7.699s 7.699s 7.681s 7.652s 7.664s 7.722s 7.653s 7.687s 7.682s
8 8.21s  8.21s  1ms 8.192s 8.204s 8.227s 8.216s 8.216s 8.134s 8.198s 8.245s 8.181s 8.239s 8.262s
```

`./bin/taps2beats --verbose --precision 1ms --latency 7ms --quantize --interval 1s:12s  examples/taps.txt`

```
  taps2beats v0.1.0

  ... reading data from ./runtime/taps.txt
  ... 12 values read from ./runtime/taps.txt
  ... using forgetting factor 0
  ... quantizing tapped beats to match estimated BPM
  ... interpolating missing beats over interval 1s..12s 
  ... 21 beats
  ... compensating for 7ms latency
  ... rounding to 10ms

BPM:    114
Offset: 310ms

1  1.36s 
2  1.89s 
3  2.41s 
4  2.94s 
5  3.47s 
6  3.99s 
7  4.52s  4.52s 0s 4.56s 4.5s  4.52s 4.52s 4.51s 4.49s 4.52s 4.52s 4.51s 4.51s
8  5.04s  5.05s 0s 5.06s 5.04s 5.05s 5.06s 5.02s 5.13s 5.03s 5.03s 5.04s 5.04s 5.06s
9  5.57s  5.57s 0s 5.6s  5.58s 5.58s 5.6s  5.57s 5.54s 5.56s 5.55s 5.58s 5.54s 5.58s
10 6.1s   6.09s 0s 6.1s  6.11s 6.13s 6.1s  6.09s 6.06s 6.07s 6.12s 6.08s 6.07s 6.1s 
11 6.62s  6.61s 0s 6.64s 6.61s 6.62s 6.61s 6.65s 6.57s 6.62s 6.65s 6.59s 6.6s  6.57s
12 7.15s  7.15s 0s 7.13s 7.13s 7.17s 7.14s 7.17s 7.12s 7.13s 7.19s 7.12s 7.16s 7.14s
13 7.67s  7.68s 0s 7.7s  7.69s 7.69s 7.69s 7.67s 7.65s 7.66s 7.72s 7.65s 7.68s 7.67s
14 8.2s   8.2s  0s 8.19s 8.2s  8.22s 8.21s 8.21s 8.13s 8.19s 8.24s 8.17s 8.23s 8.26s
15 8.73s 
16 9.25s 
17 9.78s 
18 10.3s 
19 10.83s
20 11.36s
21 11.88s
```

`taps2beats --json examples/taps.txt`

{
 "BPM": 114,
 "offset": 316000000,
 "beats": [
  {
   "at": 4524000000,
   "mean": 4524000000,
   "variance": 0,
   "taps": [
    4570000000,
    4506000000,
    ...,
    4518000000
   ]
  },
  {
   "at": 5058000000,
   "mean": 5058000000,
   "variance": 1000000,
   "taps": [
    5064000000,
    ...,
    5069000000
   ]
  },
  ...
  {
   "at": 8210000000,
   "mean": 8210000000,
   "variance": 1000000,
   "taps": [
    8192000000,
    ...,
    8262000000
   ]
  }
 ]
}


```
BPM:    114
Offset: 316ms

1 4.524s 4.524s 0s  4.57s  4.506s 4.53s  4.53s  4.518s 4.495s 4.529s 4.524s 4.518s 4.518s
2 5.058s 5.058s 1ms 5.064s 5.046s 5.058s 5.069s 5.023s 5.133s 5.04s  5.04s  5.046s 5.046s 5.069s
3 5.578s 5.578s 0s  5.604s 5.592s 5.592s 5.603s 5.58s  5.545s 5.563s 5.557s 5.586s 5.551s 5.586s
4 6.1s   6.1s   0s  6.103s 6.114s 6.137s 6.103s 6.097s 6.068s 6.079s 6.132s 6.091s 6.074s 6.109s
5 6.618s 6.618s 1ms 6.643s 6.619s 6.631s 6.613s 6.654s 6.579s 6.625s 6.654s 6.596s 6.608s 6.579s
6 7.153s 7.153s 0s  7.142s 7.136s 7.177s 7.148s 7.176s 7.13s  7.142s 7.194s 7.13s  7.165s 7.148s
7 7.686s 7.686s 1ms 7.711s 7.693s 7.699s 7.699s 7.681s 7.652s 7.664s 7.722s 7.653s 7.687s 7.682s
8 8.21s  8.21s  1ms 8.192s 8.204s 8.227s 8.216s 8.216s 8.134s 8.198s 8.245s 8.181s 8.239s 8.262s
```

#### Issues and Feature Requests

Please create an issue in the [taps2beats](https://github.com/twystd/taps2beats) _github_ repository.

#### References

1. Wang H, Song M (2011). _“Ckmeans.1d.dp: Optimal $k$-means clustering in one dimension by dynamic programming.”_ 
The R Journal, 3(2), 29–33. doi: 10.32614/RJ-2011-015.

2. Song M, Zhong H (2020). _“Efficient weighted univariate clustering maps outstanding dysregulated genomic zones in human cancers.”_ Bioinformatics. doi: 10.1093/bioinformatics/btaa613, [Published online ahead of print, 2020 Jul 3].

3. Zhong H (2019). _Model-free Gene-to-zone Network Inference of Molecular Mechanisms in Biology._ Ph.D. thesis, Department of Computer Science, New Mexico State University, Las Cruces, NM, USA.  

4. [Ckmeans.1d.dp](https://cran.r-project.org/web/packages/Ckmeans.1d.dp/index.html)

5. [W3C multi-thumb slider](https://w3c.github.io/aria-practices/examples/slider/multithumb-slider.html)

6. [AB Loop Player](https://agrahn.gitlab.io/ABLoopPlayer)

#### License

[MIT](https://github.com/twystd/taps2beats/blob/master/LICENSE)

#### Attributions

1. Icons by [FontAwesome](https://fontawesome.com):
   - _[FontAwesome](https://fontawesome.com/license)_

