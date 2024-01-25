<h1 align="center">
  <img src="static/ghtracker-logo.png" alt="subfinder" width="200px">
  <br>
</h1>

<h4 align="center">CLI tool for tracking dependents repositories and sorting result by Stars ⭐</h4>

<p align="center">
<a href="https://goreportcard.com/report/github.com/zer0yu/ghtracker"><img src="https://goreportcard.com/badge/github.com/zer0yu/ghtracker"></a>
<a href="https://github.com/zer0yu/ghtracker/issues"><img src="https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat"></a>
<a href="https://github.com/zer0yu/ghtracker/releases"><img src="https://img.shields.io/github/release/zer0yu/ghtracker"></a>
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Install</a> •
  <a href="#usage">Usage</a> •
  <a href="#why-ghtracker">Why ghtracker</a> •
  <a href="#references">References</a> •
</p>

---

`ghtracker` is a tool for tracking dependents repositories and sorting result by Stars ⭐. It has a simple architecture and is optimized for speed.

## Features

- Analysis of 'used by' or 'Dependency graph' in Github repositories and sorting result by Stars
- In-depth Search of Key Code in Used Objects 
- Support for Output in JSON Format 
- Provision of Aesthetically Pleasing and Easy-to-Read Terminal Table Outputs 
- Comprehensive Cross-Platform Compatibility 
- Efficient HTTP Request Anti-Ban Mechanism

## Installation

### From Source

`ghtracker` requires **go1.20** to install successfully. Run the following command to install the latest version:

```sh
go install github.com/zer0yu/ghtracker@latest
```

### From Release

Download from [releases](http://github.com/zer0yu/ghtracker/releases/)

## Usage

### Help

```sh
ghtracker -h
```

This will display help for the tool. Here are all the switches it supports.

```sh
CLI tool for tracking dependents repositories and sorting result by Stars

Usage:
  ghtracker [flags]

Flags:
  -d, --description     Show description of packages or repositories (performs additional request per repository)
  -h, --help            help for ghtracker
  -m, --minstar int     Minimum number of stars (default=5) (default 5)
  -f, --output string   File to write output to (optional)
  -r, --repositories    Sort repositories or packages (default packages)
  -o, --rows int        Number of showing repositories (default=10) (default 10)
  -s, --search string   search code at dependents (repositories or packages)
  -t, --table           View mode
  -k, --token string    GitHub token
  -u, --url string      URL to process
```

### Basic Usage

#### 1. Retrieves Packages in 'Used by' or 'Dependency graph' by default, and saves the addresses of projects with more than 5 output stars.

```sh
ghtracker --url https://github.com/AFLplusplus/LibAFL -t 
/
+---------------------------------------+-------+-------------+
|                  URL                  | STARS | DESCRIPTION |
+---------------------------------------+-------+-------------+
| https://github.com/AFLplusplus/LibAFL | 1.7K  |             |
| https://github.com/epi052/feroxfuzz   |   183 |             |
| https://github.com/fkie-cad/butterfly |    40 |             |
| https://github.com/z2-2z/peacock      |     0 |             |
+---------------------------------------+-------+-------------+
found 8 packages others packages are private
Complete!
```

#### 2. Get the Repositories in the 'Used by' or 'Dependency graph', by default it saves the addresses of projects with the first 10.

```sh
ghtracker --url https://github.com/AFLplusplus/LibAFL -t -r
/
+--------------------------------------------------------+-------+-------------+
|                          URL                           | STARS | DESCRIPTION |
+--------------------------------------------------------+-------+-------------+
| https://github.com/AFLplusplus/LibAFL                  | 1.7K  |             |
| https://github.com/hardik05/Damn_Vulnerable_C_Program  |   604 |             |
| https://github.com/fuzzland/ityfuzz                    |   524 |             |
| https://github.com/epi052/feroxfuzz                    |   183 |             |
| https://github.com/epi052/fuzzing-101-solutions        |   119 |             |
| https://github.com/tlspuffin/tlspuffin                 |   117 |             |
| https://github.com/Agnoctopus/Tartiflette              |    90 |             |
| https://github.com/fkie-cad/butterfly                  |    40 |             |
| https://github.com/IntelLabs/PreSiFuzz                 |    38 |             |
| https://github.com/RickdeJager/TrackmaniaFuzzer        |    32 |             |
| https://github.com/AFLplusplus/libafl_paper_artifacts  |    17 |             |
| https://github.com/novafacing/libc-fuzzer              |    12 |             |
| https://github.com/vusec/triereme                      |    11 |             |
| https://github.com/bitterbit/fuzzer-qemu               |     7 |             |
| https://github.com/atredis-jordan/libafl-workshop-blog |     7 |             |
| https://github.com/rezer0dai/LibAFL                    |     6 |             |
| https://github.com/jjjutla/Fuzz                        |     5 |             |
+--------------------------------------------------------+-------+-------------+
found 190 repositories others repositories are private
Complete!
```

#### 3. Save the result in a file in Json format

```sh
ghtracker  --url https://github.com/AFLplusplus/LibAFL -t -r --output ./test.json
```

#### 4. Search Code Pattern in 'Dependency graph' (Need Github Token)

```sh
ghtracker --url https://github.com/AFLplusplus/LibAFL --token your_token_value -t --search AFL --output ./test.json
\
[INF] https://github.com/AFLplusplus/LibAFL/blob/e117b7199ca902d462edc1de1bc0b3cb71c27aff/scripts/afl-persistent-config with 1747 stars
[INF] https://github.com/AFLplusplus/LibAFL/blob/e117b7199ca902d462edc1de1bc0b3cb71c27aff/libafl_cc/src/afl-coverage-pass.cc with 1747 stars
[INF] https://github.com/AFLplusplus/LibAFL/blob/e117b7199ca902d462edc1de1bc0b3cb71c27aff/libafl_targets/src/cmps/observers/aflpp.rs with 1747 stars
[INF] https://github.com/AFLplusplus/LibAFL/blob/e117b7199ca902d462edc1de1bc0b3cb71c27aff/fuzzers/forkserver_libafl_cc/src/bin/libafl_cc.rs with 1747 stars
[INF] https://github.com/AFLplusplus/LibAFL/blob/e117b7199ca902d462edc1de1bc0b3cb71c27aff/fuzzers/libfuzzer_libpng_aflpp_ui/README.md with 1747 stars
...
...
```

#### 5. Get all repositories in the 'Dependency graph'

```sh
ghtracker  --url https://github.com/AFLplusplus/LibAFL -t -r -m 0 --output ./test.json
```

# Why ghtracker

1. Gtihub does not support sorting the 'Dependency graph' by the number of Stars from 2019 until now, so ghtracker 
   was born. Detials in issue in [1537](https://github.com/isaacs/github/issues/1537)
2. GHTOPDEP does not continue to update support, and there are some issues in the community, e.g., [issue](https://github.com/github-tooling/ghtopdep). ghtracker has optimized and improved on this, and has addressed issues in the community.

# References
- [GHTOPDEP](https://github.com/github-tooling/ghtopdep)
