# wc
Well-known wc utility written in Go. Simplified version.

# Installation
```bash
go get github.com/avegner/wc
```

# Usage
```
Usage: wc FILE...
Print newline, word, byte and character counts for each FILE, and a total line if
more than one FILE is specified. A word is a non-zero-length sequence of
characters delimited by white space. Unicode BOM isn't considered a character.
Only ASCII and UTF-8 encodings are supported. Invalid characters are ignored.
```

# Benchmarking
* install time utility:
```bash
sudo apt update
sudo apt install -y time
```
* run original wc:
```bash
$(which time) -v wc -mclw <file>...
```
* run go wc:
```bash
$(which time) -v <path-to-go-version>/wc <file>...
```
* compare mem and cpu usage stats
