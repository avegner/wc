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

var (
	errUsage  = errors.New("usage error")
	errFailed = errors.New("failed")
)

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

type statsResult struct {
	st       *stats
	filePath string
	err      error
}

func run() error {
	if len(os.Args) < 2 {
		return errUsage
	}

	filePaths := os.Args[1:]
	resc := make(chan *statsResult)

	for _, fp := range filePaths {
		p := fp
		go func() {
			st, err := getStats(p)
			resc <- &statsResult{st, p, err}
		}()
	}

	var err error
	total := stats{}
	got := 0

	for res := range resc {
		if res.err == nil {
			total.lines += res.st.lines
			total.words += res.st.words
			total.bytes += res.st.bytes
			total.chars += res.st.chars
			mlog.Printf("%v %s", res.st, res.filePath)
		} else {
			err = errFailed
			mlog.Printf("err: %v", res.err)
		}
		if got++; got == len(filePaths) {
			close(resc)
		}
	}
	if len(filePaths) > 1 {
		mlog.Printf("%v total", total)
	}

	return err
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
