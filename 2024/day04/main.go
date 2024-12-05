package main

import (
	"fmt"
	"os"
	"strings"

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

func task1(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	total := 0
	total += countHorizontal(lines, "XMAS", "SAMX")
	total += countVertical(lines, "XMAS", "SAMX")
	total += countDiagonalDown(lines, "XMAS", "SAMX")
	total += countDiagonalUp(lines, "XMAS", "SAMX")

	fmt.Printf("Total: %d\n", total)

	return nil
}

func countHorizontal(haystack []string, needles ...string) int {
	count := 0
	for _, line := range haystack {
		for _, needle := range needles {
			count += strings.Count(line, needle)
		}
	}

	return count
}

func countVertical(haystack []string, needles ...string) int {
	line0 := haystack[0]
	buf := make([]byte, len(haystack))
	count := 0

	for i := range line0 {
		for j, line := range haystack {
			buf[j] = line[i]
		}

		for _, needle := range needles {
			count += strings.Count(string(buf), needle)
		}
	}

	return count
}

type Pos struct {
	X, Y int
}

func (p Pos) Add(x, y int) Pos {
	return Pos{
		X: p.X + x,
		Y: p.Y + y,
	}
}

func countDiagonalDown(haystack []string, needles ...string) int {
	buf := make([]byte, 0, len(haystack))
	count := 0

	for x := -len(haystack) + 1; x < len(haystack[0]); x++ {
		p := Pos{
			X: x,
			Y: 0,
		}

		for ; p.Y < len(haystack); p = p.Add(1, 1) {
			if p.X < 0 || p.X >= len(haystack[0]) {
				continue
			}

			buf = append(buf, haystack[p.Y][p.X])
		}

		for _, needle := range needles {
			count += strings.Count(string(buf), needle)
		}

		buf = buf[:0]
	}

	return count
}

func countDiagonalUp(haystack []string, needles ...string) int {
	buf := make([]byte, 0, len(haystack))
	count := 0

	for x := len(haystack[0])*2 - 1; x >= 0; x-- {
		p := Pos{
			X: x,
			Y: 0,
		}

		for ; p.Y < len(haystack); p = p.Add(-1, 1) {
			if p.X < 0 || p.X >= len(haystack[0]) {
				continue
			}

			buf = append(buf, haystack[p.Y][p.X])
		}

		for _, needle := range needles {
			count += strings.Count(string(buf), needle)
		}

		buf = buf[:0]
	}

	return count
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func matchPattern(lines []string, p Pos) bool {
	if lines[p.Y+1][p.X+1] != 'A' {
		return false
	}

	a := lines[p.Y][p.X]
	b := lines[p.Y][p.X+2]
	c := lines[p.Y+2][p.X]
	d := lines[p.Y+2][p.X+2]

	first := (a == 'M' && d == 'S') || (a == 'S' && d == 'M')
	second := (b == 'M' && c == 'S') || (b == 'S' && c == 'M')

	return first && second
}

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	count := 0
	for y, line := range lines[:len(lines)-2] {
		for x := range line[:len(line)-2] {
			if matchPattern(lines, Pos{x, y}) {
				count++
			}
		}
	}

	fmt.Printf("Count: %d\n", count)
	return nil
}
