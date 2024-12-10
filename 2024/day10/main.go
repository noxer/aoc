package main

import (
	"fmt"
	"os"
	"time"

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

var directions = []utils.Vec{
	{X: +1}, // right
	{X: -1}, // left
	{Y: +1}, // down
	{Y: -1}, // up
}

type Map struct {
	data []string
}

func (m Map) IterateTrailheads() func(func(utils.Vec) bool) {
	return func(yield func(utils.Vec) bool) {
		for y, row := range m.data {
			for x, r := range []byte(row) {
				if r != '0' {
					continue
				}

				if !yield(utils.Vec{X: x, Y: y}) {
					return
				}
			}
		}
	}
}

func (m Map) Get(pos utils.Vec) int {
	if pos.Y < 0 || pos.Y >= len(m.data) {
		return -1
	}

	row := m.data[pos.Y]
	if pos.X < 0 || pos.X >= len(row) {
		return -1
	}

	b := row[pos.X] - '0'
	return int(b)
}

func (m Map) FindScore(start utils.Vec) int {
	positions := make([]utils.Vec, 1, 256)
	positions[0] = start
	peeks := make(map[utils.Vec]struct{})
	into := make([]utils.Vec, 4)

	for len(positions) > 0 {
		pos := positions[0]
		copy(positions, positions[1:])
		positions = positions[:len(positions)-1]

		if m.Get(pos) == 9 {
			peeks[pos] = struct{}{}
			continue
		}

		into = m.findNextSteps(pos, into[:0])
		positions = append(positions, into...)
	}

	return len(peeks)
}

func (m Map) findNextSteps(from utils.Vec, into []utils.Vec) []utils.Vec {
	fromVal := m.Get(from)

	for _, dir := range directions {
		newPos := from.Add(dir)
		newVal := m.Get(newPos)
		if newVal != fromVal+1 {
			continue
		}

		into = append(into, newPos)
	}

	return into
}

func task1(args []string) error {
	data, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	m := Map{data: data}

	sum := 0
	for head := range m.IterateTrailheads() {
		sum += m.FindScore(head)
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (m Map) FindRating(start utils.Vec) int {
	positions := []utils.Vec{start}
	trails := 0
	into := make([]utils.Vec, 4)

	for len(positions) > 0 {
		pos := positions[0]
		copy(positions, positions[1:])
		positions = positions[:len(positions)-1]

		if m.Get(pos) == 9 {
			trails++
			continue
		}

		into = m.findNextSteps(pos, into[:0])
		positions = append(positions, into...)
	}

	return trails
}

func task2(args []string) error {
	data, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	m := Map{data}

	start := time.Now()

	sum := 0
	for head := range m.IterateTrailheads() {
		sum += m.FindRating(head)
	}

	elapsed := time.Since(start)

	fmt.Printf("Sum: %d (%s)\n", sum, elapsed)

	return nil
}
