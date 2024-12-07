package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	var err error

	if len(os.Args) <= 1 {
		fmt.Println("Missing argument, please specify the task you want to execute (1 or 2).")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "1":
		err = task1(os.Args[2:])
	case "2":
		err = task2(os.Args[2:])
	default:
		fmt.Println("Invalid argument, please specify the task you want to execute (1 or 2).")
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error executing task %s: %s\n", os.Args[1], err)
		os.Exit(1)
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

var (
	Up    = Vec{Y: -1}
	Down  = Vec{Y: +1}
	Left  = Vec{X: -1}
	Right = Vec{X: +1}
)

var rightTurn = map[Vec]Vec{
	Up:    Right,
	Right: Down,
	Down:  Left,
	Left:  Up,
}

type Vec struct {
	X, Y int
}

func (v Vec) Add(o Vec) Vec {
	return Vec{v.X + o.X, v.Y + o.Y}
}

type Map struct {
	size Vec
	data map[Vec]byte
}

func (m Map) Contains(v Vec) bool {
	return v.X >= 0 && v.X < m.size.X &&
		v.Y >= 0 && v.Y < m.size.Y
}

func (m Map) Obstacle(v Vec) bool {
	return m.data[v] == '#'
}

func (m Map) Mark(v Vec) {
	m.data[v] = 'X'
}

func (m Map) CountMarks() int {
	count := 0
	for _, b := range m.data {
		if b == 'X' {
			count++
		}
	}
	return count
}

func parseMap(name string) (Map, Vec, error) {
	f, err := os.Open(name)
	if err != nil {
		return Map{}, Vec{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	m := Map{
		data: make(map[Vec]byte),
	}
	start := Vec{}
	y := 0
	for s.Scan() {
		for x, b := range s.Bytes() {
			switch b {
			case '#':
				m.data[Vec{x, y}] = '#'
			case '^':
				start = Vec{x, y}
				m.data[start] = 'X'
			}
		}

		y++
		m.size.X = len(s.Bytes())
		m.size.Y = y
	}

	return m, start, nil
}

func task1(args []string) error {
	m, s, err := parseMap(args[0])
	if err != nil {
		return err
	}

	pos := s
	dir := Up
	for {
		newPos := pos.Add(dir)
		if m.Obstacle(newPos) {
			dir = rightTurn[dir]
			continue
		}

		if !m.Contains(newPos) {
			break
		}

		m.Mark(newPos)
		pos = newPos
	}

	fmt.Printf("Marked: %d\n", m.CountMarks())

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

var dirMarker = map[Vec]byte{
	Up:    '^',
	Down:  'v',
	Left:  '<',
	Right: '>',
}

func (m Map) Walk(start Vec) {
	pos := start
	dir := Up
	for {
		newPos := pos.Add(dir)
		if m.Obstacle(newPos) {
			dir = rightTurn[dir]
			continue
		}

		if !m.Contains(newPos) {
			break
		}

		m.Mark(newPos)
		pos = newPos
	}
}

func (m Map) MarkAs(v Vec, b byte) {
	m.data[v] = b
}

func (m Map) Check(v Vec, b byte) bool {
	return m.data[v] == b
}

func (m Map) WalkWithDir(start Vec) bool {
	pos := start
	dir := Up
	for {
		newPos := pos.Add(dir)
		if m.Obstacle(newPos) {
			dir = rightTurn[dir]
			continue
		}

		if !m.Contains(newPos) {
			return false
		}

		if m.Check(newPos, dirMarker[dir]) {
			return true
		}

		m.MarkAs(newPos, dirMarker[dir])
		pos = newPos
	}
}

func (m Map) IterateEqual(e byte) func(func(Vec) bool) {
	return func(yield func(Vec) bool) {
		for v, b := range m.data {
			if b == e {
				if !yield(v) {
					return
				}
			}
		}
	}
}

func (m Map) Copy() Map {
	c := Map{
		size: m.size,
		data: make(map[Vec]byte, len(m.data)),
	}

	for v, b := range m.data {
		c.data[v] = b
	}

	return c
}

func (m Map) Reset() {
	for v, b := range m.data {
		if b != '#' {
			delete(m.data, v)
		}
	}
}

func task2(args []string) error {
	m, s, err := parseMap(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	m.Walk(s)
	c := m.Copy()
	count := 0
	for v := range m.IterateEqual('X') {
		if v == s {
			continue
		}

		c.MarkAs(v, '#')
		if c.WalkWithDir(s) {
			count++
		}

		c.MarkAs(v, '.')
		c.Reset()
	}

	elapsed := time.Since(start)

	fmt.Printf("Count: %d (%s)\n", count, elapsed)

	return nil
}
