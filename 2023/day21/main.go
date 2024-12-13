package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

type empty = struct{}

func main() {
	// f, err := os.Create("profile.pprof")
	// if err != nil {
	// 	fmt.Printf("Error: %s\n", err)
	// 	os.Exit(1)
	// }
	// defer f.Close()

	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

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

var directions = []Pos{
	{X: +1}, // right
	{X: -1}, // left
	{Y: +1}, // down
	{Y: -1}, // up
}

type Pos struct {
	X, Y int
}

func (p Pos) Hash() int {
	buf := make([]byte, 0, 8)
	buf = binary.AppendVarint(buf, int64(p.X))
	buf = binary.AppendVarint(buf, int64(p.Y))

	return int(binary.LittleEndian.Uint64(buf))
}

func (p Pos) Add(o Pos) Pos {
	return Pos{p.X + o.X, p.Y + o.Y}
}

func (p Pos) ID() uint64 {
	return uint64(p.X) + uint64(p.Y)
}

type Map map[Pos]byte

func (m Map) CanGo(p Pos) bool {
	_, ok := m[p]
	return ok
}

func loadMap(name string) (Map, Pos, error) {
	m := make(Map)

	f, err := os.Open(name)
	if err != nil {
		return nil, Pos{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	start := Pos{}
	y := 0
	for s.Scan() {
		for x, b := range s.Bytes() {
			switch b {
			case 'S':
				start = Pos{x, y}
				fallthrough
			case '.':
				m[Pos{x, y}] = '.'
			}
		}

		y++
	}

	return m, start, nil
}

func move(m Map, positions map[Pos]struct{}) map[Pos]struct{} {
	newPositions := make(map[Pos]struct{})

	for pos := range positions {
		for _, dir := range directions {
			newPos := pos.Add(dir)
			if m.CanGo(newPos) {
				newPositions[newPos] = struct{}{}
			}
		}
	}

	return newPositions
}

func task1(args []string) error {
	m, s, err := loadMap(args[0])
	if err != nil {
		return err
	}

	positions := map[Pos]struct{}{
		s: {},
	}

	for range 64 {
		positions = move(m, positions)
	}

	fmt.Printf("Count: %d\n", len(positions))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type Set[T comparable] map[T]struct{}

func (s Set[T]) Put(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Merge(o Set[T]) {
	for v := range o {
		s[v] = struct{}{}
	}
}

func (s Set[T]) MergeWith(o Set[T], f func(T) T) {
	for v := range o {
		s[f(v)] = struct{}{}
	}
}

// type Set[T comparable] []T

// func (s *Set[T]) Has(v T) bool {
// 	for _, t := range *s {
// 		if t == v {
// 			return true
// 		}
// 	}
// 	return false
// }

// func (s *Set[T]) Put(v T) {
// 	if !s.Has(v) {
// 		*s = append(*s, v)
// 	}
// }

// func (s *Set[T]) Merge(o *Set[T]) {
// 	// *s = slices.Grow(*s, len(*o))
// 	for _, t := range *o {
// 		s.Put(t)
// 	}
// }

// func (s *Set[T]) MergeWith(o *Set[T], f func(T) T) {
// 	// *s = slices.Grow(*s, len(*o))
// 	for _, t := range *o {
// 		s.Put(f(t))
// 	}
// }

type Map2 struct {
	m map[Pos]Set[Pos]
	s Pos
}

func (m Map2) Move(p Pos, result Map2) {
	s := m.m[p]
	if len(s) == 0 {
		return
	}

	for _, dir := range directions {
		newPos := p.Add(dir)
		if !m.CanGo(newPos) {
			continue
		}

		offset := m.CheckNextGarden(newPos)
		target := result.m[m.Lep(newPos)]
		target.MergeWith(s, func(p Pos) Pos {
			return p.Add(offset)
		})
	}
}

func (m Map2) Lep(p Pos) Pos {
	p.X %= m.s.X
	p.Y %= m.s.Y

	if p.X < 0 {
		p.X += m.s.X
	}
	if p.Y < 0 {
		p.Y += m.s.Y
	}

	return p
}

func (m Map2) CheckNextGarden(p Pos) Pos {
	dir := Pos{}

	if p.X < 0 {
		dir.X = -1
	} else if p.X >= m.s.X {
		dir.X = +1
	}

	if p.Y < 0 {
		dir.Y = -1
	} else if p.Y >= m.s.Y {
		dir.Y = +1
	}

	return dir
}

func (m Map2) Reset() {
	for p := range m.m {
		m.m[p] = make(Set[Pos])
	}
}

func (m Map2) EmptyCopy() Map2 {
	c := Map2{
		m: make(map[Pos]Set[Pos], len(m.m)),
		s: m.s,
	}

	for p := range m.m {
		c.m[p] = make(Set[Pos])
	}

	return c
}

func (m Map2) Count() uint64 {
	sum := uint64(0)
	for _, s := range m.m {
		sum += uint64(len(s))
	}
	return sum
}

func (m Map2) CanGo(p Pos) bool {
	p.X %= m.s.X
	p.Y %= m.s.Y

	if p.X < 0 {
		p.X += m.s.X
	}
	if p.Y < 0 {
		p.Y += m.s.Y
	}

	_, ok := m.m[p]
	return ok
}

func loadMap2(name string) (m Map2, err error) {
	m.m = make(map[Pos]Set[Pos])

	f, err := os.Open(name)
	if err != nil {
		return Map2{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	y := 0
	for s.Scan() {
		for x, b := range s.Bytes() {
			m.s.X = max(m.s.X, x)

			switch b {
			case 'S':
				m.m[Pos{x, y}] = Set[Pos]{Pos{}: empty{}}
			case '.':
				m.m[Pos{x, y}] = Set[Pos]{}
			}
		}

		y++
	}

	m.s.X++
	m.s.Y = y
	return m, nil
}

func move2(m, c Map2) {
	for p := range m.m {
		m.Move(p, c)
	}
}

func task2(args []string) error {
	m, err := loadMap2(args[0])
	if err != nil {
		return err
	}

	c := m.EmptyCopy()

	start := time.Now()
	for i := range 500 {
		fmt.Printf("Computing iteration %d...", i)
		loopStart := time.Now()
		move2(m, c)
		m, c = c, m
		// c.Reset()
		fmt.Println("ok |", time.Since(loopStart))
	}
	duration := time.Since(start)

	fmt.Printf("Count: %d (%s)\n", m.Count(), duration)

	return nil
}
