# wc
Well-known wc utility written in Go. Simplified version.

# Installation
```bash
GOFLAGS='-ldflags=-s -trimpath' go get github.com/avegner/wc
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
$(which time) -v $(go env GOPATH)/bin/wc <file>...
```
* compare mem and cpu usage stats

# UPX for Go bins
Static Go binaries are significantly larger than C ones. Use UPX to reduce binary size:
```bash
sudo apt update
sudo apt install -y upx
upx -9 -o $(go env GOPATH)/bin/wc.upx $(go env GOPATH)/bin/wc
```

# Output Format
```
<lines> <words> <chars> <bytes> <file-path>
...
```

# Stats
Wall clock

Test | wc -mclw | Go wc | Go wc.upx
--- | --- | --- | ---
1 x 320.4 MiB UTF-8 file | 2.14 s | 3.69 s | 3.77 s
10 x 320.4 MiB UTF-8 files | 21.51 s | 20.59 s | 20.65 s

Maximum RSS

Test | wc -mclw | Go wc | Go wc.upx
--- | --- | --- | ---
1 x 320.4 MiB UTF-8 file | 2220 KiB | 2132 KiB | 2168 KiB
10 x 320.4 MiB UTF-8 files | 2228 KiB | 2384 KiB | 2424 KiB
