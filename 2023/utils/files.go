package utils

import (
	"bufio"
	"os"
	"strings"
)

// ReadLines reads all the lines from a file and returns the space-trimmed lines in a slice.
func ReadLines(name string) ([]string, error) {
	return ReadLinesTransform(name, strings.TrimSpace)
}

// ReadLinesTransform reads all lines from a file and transforms them though the transform function.
func ReadLinesTransform[T any](name string, transformer func(string) T) ([]T, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var ts []T
	s := bufio.NewScanner(f)

	for s.Scan() {
		t := transformer(s.Text())
		ts = append(ts, t)
	}

	return ts, s.Err()
}
