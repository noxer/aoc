package utils

import (
	"bufio"
	"bytes"
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

type Vec struct {
	X, Y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y}
}

func (v Vec) Mul(s int) Vec {
	return Vec{v.X * s, v.Y * s}
}

func (v Vec) Zero() bool {
	return v.X == 0 && v.Y == 0
}

func ReadMap(name string, ignore ...byte) (map[Vec]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	y := 0
	m := make(map[Vec]byte)
	for s.Scan() {
		for x, b := range s.Bytes() {
			if bytes.Contains(ignore, []byte{b}) {
				continue
			}

			m[Vec{x, y}] = b
		}

		y++
	}

	return m, s.Err()
}

func ReadMapWithSize(name string, ignore ...byte) (map[Vec]byte, Vec, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, Vec{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	y := 0
	m := make(map[Vec]byte)
	size := Vec{}
	for s.Scan() {
		for x, b := range s.Bytes() {
			if bytes.Contains(ignore, []byte{b}) {
				continue
			}

			m[Vec{x, y}] = b
		}

		y++
		size.X = max(size.X, len(s.Bytes()))
	}

	size.Y = y
	return m, size, s.Err()
}
