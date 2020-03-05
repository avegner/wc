package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"
)

var mlog = log.New(os.Stderr, "", 0)

var errUsage = errors.New("usage error")

func main() {
	if err := run(); err != nil {
		if err == errUsage {
			mlog.Print(`Usage: wc FILE...
Print newline, word, byte and character counts for each FILE, and a total line if
more than one FILE is specified. A word is a non-zero-length sequence of
characters delimited by white space. Unicode BOM isn't considered a character.
Only ASCII and UTF-8 encodings are supported. Invalid characters are ignored.`)
			os.Exit(2)
		}
		mlog.Fatalf("err: %v", err)
	}
}

const (
	uninitedWordState = iota
	beforeWordState
	insideWordState
)

type stats struct {
	lines int
	words int
	bytes int
	chars int
}

func (st stats) String() string {
	return fmt.Sprintf("%9d %9d %9d %9d", st.lines, st.words, st.chars, st.bytes)
}

func run() error {
	if len(os.Args) < 2 {
		return errUsage
	}

	// TODO: introduce multi-threading for many files
	total := stats{}
	filePaths := os.Args[1:]

	for _, fp := range filePaths {
		st, err := getStats(fp)
		if err != nil {
			return err
		}

		total.lines += st.lines
		total.words += st.words
		total.bytes += st.bytes
		total.chars += st.chars

		mlog.Printf("%v %s", st, fp)
	}

	if len(filePaths) > 1 {
		mlog.Printf("%v total", total)
	}

	return nil
}

func getStats(filePath string) (*stats, error) {
	ws := uninitedWordState
	st := stats{}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)

	for {
		r, sz, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		st.bytes += sz

		// skip invalid chars
		if r == unicode.ReplacementChar && sz == 1 {
			continue
		}
		// skip BOM
		if r == '\uFEFF' {
			continue
		}
		st.chars++

		if r == '\n' {
			st.lines++
		}

		switch ws {
		case beforeWordState:
			if !unicode.IsSpace(r) {
				ws = insideWordState
			}
		case insideWordState:
			if unicode.IsSpace(r) {
				ws = beforeWordState
				st.words++
			}
		case uninitedWordState:
			if unicode.IsSpace(r) {
				ws = beforeWordState
			} else {
				ws = insideWordState
			}
		}
	}

	return &st, nil
}
