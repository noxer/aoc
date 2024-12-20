package main

import (
	"fmt"
	"os"

	"github.com/noxer/aoc/2024/utils"
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

type Maze struct {
	data  map[utils.Vec]byte
	size  utils.Vec
	start utils.Vec
	end   utils.Vec
	times map[utils.Vec]int
}

func (m *Maze) FindStartAndEnd() {
	for pos, val := range m.data {
		switch val {
		case 'S':
			m.start = pos
			delete(m.data, pos)
		case 'E':
			m.end = pos
			delete(m.data, pos)
		}
	}
}

func (m Maze) Get(pos utils.Vec) byte {
	return m.data[pos]
}

func (m Maze) Wall(pos utils.Vec) bool {
	return m.Get(pos) == '#'
}

func (m Maze) Print() {
	for y := range m.size.Y {
		for x := range m.size.X {
			pos := utils.Vec{X: x, Y: y}

			if b, ok := m.data[pos]; ok {
				fmt.Print(string(b))
			} else {
				if time, ok := m.times[pos]; ok {
					fmt.Print(time % 10)
				} else {
					fmt.Print(".")
				}
			}
		}
		fmt.Println()
	}
}

func (m Maze) PopulateTimes() {
	last := m.start
	pos := m.start
	time := 0
	for {
		m.times[pos] = time

		if pos == m.end {
			break
		}

		for _, dir := range utils.Directions {
			next := pos.Add(dir)
			if next == last || m.Wall(next) {
				continue
			}

			last = pos
			pos = next
			time++
			break
		}
	}
}

type Shortcut struct {
	Start, End utils.Vec
	Saves      int
}

func (m Maze) FindShortcuts() []Shortcut {
	var shortcuts []Shortcut

	for pos, start := range m.times {
		for _, dir := range utils.Directions {
			next := pos.Add(dir)
			if !m.Wall(next) {
				continue
			}

			nextNext := next.Add(dir)
			nextNextTime := m.times[nextNext]

			if nextNextTime > start+2 {
				shortcuts = append(shortcuts, Shortcut{
					Start: next,
					End:   nextNext,
					Saves: nextNextTime - (start + 2),
				})
			}
		}
	}

	return shortcuts
}

func task1(args []string) error {
	data, size, err := utils.ReadMapWithSize(args[0], '.')
	if err != nil {
		return err
	}

	maze := Maze{
		data:  data,
		size:  size,
		times: make(map[utils.Vec]int),
	}
	maze.FindStartAndEnd()

	maze.PopulateTimes()

	shortcuts := maze.FindShortcuts()

	fmt.Println(shortcuts)

	count := 0
	for _, shortcut := range shortcuts {
		if shortcut.Saves >= 100 {
			count++
		}
	}

	fmt.Printf("Shortcuts >= 100ps: %d\n", count)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type FromTo struct {
	From, To utils.Vec
}

func Distance(a, b utils.Vec) int {
	return (max(a.X, b.X) - min(a.X, b.X)) + (max(a.Y, b.Y) - min(a.Y, b.Y))
}

func (m Maze) FindLongShortcuts() map[FromTo]int {
	shortcuts := make(map[FromTo]int)

	for from, start := range m.times {
		for to, end := range m.times {
			if start >= end {
				continue
			}

			dist := Distance(from, to)
			if dist > 20 {
				continue
			}

			if end <= start+dist {
				continue
			}

			shortcuts[FromTo{from, to}] = end - (start + dist)
		}
	}

	return shortcuts
}

func task2(args []string) error {

	data, size, err := utils.ReadMapWithSize(args[0], '.')
	if err != nil {
		return err
	}

	maze := Maze{
		data:  data,
		size:  size,
		times: make(map[utils.Vec]int),
	}
	maze.FindStartAndEnd()

	maze.PopulateTimes()

	shortcuts := maze.FindLongShortcuts()

	fmt.Println(shortcuts)

	count := 0
	for _, saves := range shortcuts {
		if saves >= 100 {
			count++
		}
	}

	fmt.Printf("Shortcuts >= 100ps: %d\n", count)

	return nil
}
