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

type Vec struct {
	X, Y int
}

func (v Vec) To(o Vec) Vec {
	return Vec{
		X: o.X - v.X,
		Y: o.Y - v.Y,
	}
}

func (v Vec) Add(o Vec) Vec {
	return Vec{
		X: o.X + v.X,
		Y: o.Y + v.Y,
	}
}

func calcAntinodes(a, b Vec) (Vec, Vec) {
	return a.To(b).Add(b), b.To(a).Add(a)
}

type Map struct {
	Size     Vec
	Antennas map[byte][]Vec
}

func (m Map) Contains(pos Vec) bool {
	return pos.X >= 0 && pos.X < m.Size.X && pos.Y >= 0 && pos.Y < m.Size.Y
}

func (m Map) CountAntinodes() int {
	antinodes := make(map[Vec]struct{})

	for _, antennas := range m.Antennas {
		for i, a := range antennas {
			for _, b := range antennas[i+1:] {
				n1, n2 := calcAntinodes(a, b)
				if m.Contains(n1) {
					antinodes[n1] = struct{}{}
				}
				if m.Contains(n2) {
					antinodes[n2] = struct{}{}
				}
			}
		}
	}

	return len(antinodes)
}

func parseMap(name string) (Map, error) {
	f, err := os.Open(name)
	if err != nil {
		return Map{}, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	y := 0
	m := Map{
		Antennas: make(map[byte][]Vec),
	}
	for s.Scan() {
		for x, b := range s.Bytes() {
			if b == '.' {
				continue
			}

			m.Antennas[b] = append(m.Antennas[b], Vec{x, y})
		}

		y++
		m.Size.X = max(m.Size.X, len(s.Bytes()))
	}

	m.Size.Y = y

	return m, nil
}

func task1(args []string) error {
	m, err := parseMap(args[0])
	if err != nil {
		return err
	}

	count := m.CountAntinodes()

	fmt.Printf("Antinodes: %d\n", count)
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (m Map) CountAntinodes2() int {
	antinodes := make(map[Vec]struct{})

	for _, antennas := range m.Antennas {
		for i, a := range antennas {
			for _, b := range antennas[i+1:] {
				diff := a.To(b)
				for anti := a.Add(diff); m.Contains(anti); anti = anti.Add(diff) {
					antinodes[anti] = struct{}{}
				}
				diff = b.To(a)
				for anti := b.Add(diff); m.Contains(anti); anti = anti.Add(diff) {
					antinodes[anti] = struct{}{}
				}
			}
		}
	}

	return len(antinodes)
}

func task2(args []string) error {
	m, err := parseMap(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	count := m.CountAntinodes2()

	elapsed := time.Since(start)

	fmt.Printf("Antinodes: %d (%s)\n", count, elapsed)
	return nil
}
