package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/noxer/aoc/2023/utils"
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

var directions = map[string]Pos{
	"R": {X: +1},
	"L": {X: -1},
	"D": {Y: +1},
	"U": {Y: -1},

	"0": {X: +1},
	"1": {Y: +1},
	"2": {X: -1},
	"3": {Y: -1},
}

type Command struct {
	Direction Pos
	Length    int
	Color     string
}

type Pos struct {
	X, Y int
}

func (p Pos) Add(o Pos) Pos {
	return Pos{
		X: p.X + o.X,
		Y: p.Y + o.Y,
	}
}

func (p Pos) Mul(m int) Pos {
	return Pos{
		X: p.X * m,
		Y: p.Y * m,
	}
}

type Trench struct {
	Start, End Pos
}

func (t Trench) String() string {
	return fmt.Sprintf("(%d, %d) -- (%d, %d)", t.Start.X, t.Start.Y, t.End.X, t.End.Y)
}

func (t *Trench) Normalize() {
	if t.Start.X > t.End.X || t.Start.Y > t.End.Y {
		t.Start, t.End = t.End, t.Start
	}
}

func (t Trench) Intersects(y int) bool {
	return t.Start.Y <= y && t.End.Y >= y
}

func (t Trench) Length() int {
	return (t.End.X - t.Start.X) + (t.End.Y - t.Start.Y) + 1
}

func (t Trench) Vertical() bool {
	return t.Start.X == t.End.X
}

func (t Trench) Horizontal() bool {
	return t.Start.Y == t.End.Y
}

func (t Trench) IsStart(y int) bool {
	return t.Start.Y == y
}

func (t Trench) IsEnd(y int) bool {
	return t.End.Y == y
}

func task1(args []string) error {
	if len(args) == 0 {
		return errors.New("need file name")
	}

	cmds, err := utils.ReadLinesTransform(args[0], func(s string) Command {
		cmd := Command{}
		fields := strings.Fields(s)
		cmd.Direction = directions[fields[0]]
		cmd.Length, _ = strconv.Atoi(fields[1])
		cmd.Color = fields[2]

		return cmd
	})
	if err != nil {
		return err
	}

	m := make(map[Pos]string)
	digTrenches(m, cmds)

	printMap(m)

	trenchCount := len(m)
	count := countHole(m)

	printMap(m)

	fmt.Println("Count:", trenchCount+int(count))

	return nil
}

func digHole(m map[Pos]string) {
	l := 0
	r := 0
	u := 0
	d := 0

	for pos := range m {
		l = min(pos.X, l)
		r = max(pos.X, r)
		u = min(pos.Y, u)
		d = max(pos.Y, d)
	}

	for x := l; x <= r; x++ {
		for y := u; y <= d; y++ {
			if shootRay(m, Pos{x, y}, r) {
				m[Pos{x, y}] = "#"
			}
		}
	}
}

func shootRay(m map[Pos]string, pos Pos, maxX int) bool {
	if _, ok := m[pos]; ok {
		return false
	}

	y := pos.Y
	wasTrench := false
	wasTrenchFromAbove := false
	count := 0
	for x := pos.X; x <= maxX+1; x++ {
		if _, ok := m[Pos{X: x, Y: y}]; ok {
			if !wasTrench {
				count++
				wasTrench = true
				_, wasTrenchFromAbove = m[Pos{X: x, Y: y - 1}]
			}
		} else {
			if wasTrench {
				wasTrench = false
				_, wasTrenchToBelow := m[Pos{X: x - 1, Y: y + 1}]
				if wasTrenchFromAbove != wasTrenchToBelow {
					count--
				}
			}
		}
	}

	return count%2 == 1
}

func digTrenches(m map[Pos]string, cmds []Command) {
	pos := Pos{0, 0}
	for _, cmd := range cmds {
		fmt.Printf("Digging trench from %v with %d length\n", pos, cmd.Length)
		pos = digTrench(m, cmd, pos)
	}
}

func digTrench(m map[Pos]string, cmd Command, pos Pos) Pos {
	for range cmd.Length {
		pos = pos.Add(cmd.Direction)
		m[pos] = cmd.Color
	}

	return pos
}

func printMap(m map[Pos]string) {
	l := 0
	r := 0
	u := 0
	d := 0

	for pos := range m {
		l = min(pos.X, l)
		r = max(pos.X, r)
		u = min(pos.Y, u)
		d = max(pos.Y, d)
	}

	for y := u; y <= d; y++ {
		for x := l; x <= r; x++ {
			if _, ok := m[Pos{x, y}]; ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func calcBounds(trenches []Trench) (l, r, u, d int) {
	for _, trench := range trenches {
		l = min(trench.Start.X, l)
		r = max(trench.Start.X, r)
		u = min(trench.Start.Y, u)
		d = max(trench.Start.Y, d)

		l = min(trench.End.X, l)
		r = max(trench.End.X, r)
		u = min(trench.End.Y, u)
		d = max(trench.End.Y, d)
	}

	return l, r, u, d
}

func task2(args []string) error {
	if len(args) == 0 {
		return errors.New("need file name")
	}

	cmds, err := utils.ReadLinesTransform(args[0], func(s string) Command {
		// cmd := Command{}
		// fields := strings.Fields(s)
		// cmd.Direction = directions[fields[0]]
		// cmd.Length, _ = strconv.Atoi(fields[1])
		// cmd.Color = fields[2]

		// return cmd

		cmd := Command{}
		inst := strings.Fields(s)[2]
		length, _ := strconv.ParseInt(inst[2:7], 16, 64)
		cmd.Length = int(length)
		cmd.Direction = directions[inst[7:8]]

		return cmd
	})
	if err != nil {
		return err
	}

	fmt.Printf("Loaded %d commands\n", len(cmds))

	trenches := convertTrenches(cmds)

	// for _, trench := range trenches {
	// 	fmt.Printf("(%d, %d) -- (%d, %d)\n", trench.Start.X, trench.Start.Y, trench.End.X, trench.End.Y)
	// }

	_, _, u, d := calcBounds(trenches)

	fmt.Printf("Bounds: %d - %d\n", u, d)

	count := 0
	for y := u; y <= d; y++ {
		row := trenchesForY(trenches, y)
		// fmt.Printf("Row %d: %v\n", y, row)

		inside := false
		start := -1
		skip := -1
		for _, trench := range row {
			if trench.Vertical() && !trench.IsEnd(y) {
				inside = !inside
				if inside {
					start = trench.Start.X
					skip = 1
				} else {
					end := trench.End.X
					diff := end - start
					count += diff - skip
				}
			}

			if trench.Vertical() && trench.IsEnd(y) {
				skip++
			}

			if trench.Horizontal() {
				skip += trench.Length() - 2
			}
		}
	}

	trenchCount := 0
	for _, trench := range trenches {
		// fmt.Printf("Trench %s: %d\n", trench, trench.Length())
		trenchCount += trench.Length() - 1
	}

	fmt.Printf("Trench Count: %d\n", trenchCount)
	fmt.Printf("Count: %d\n", count+trenchCount)

	return nil
}

func convertTrenches(cmds []Command) []Trench {
	var trenches []Trench
	pos := Pos{0, 0}

	for _, cmd := range cmds {
		scaledDir := cmd.Direction.Mul(cmd.Length)
		end := pos.Add(scaledDir)

		trench := Trench{Start: pos, End: end}
		trench.Normalize()
		trenches = append(trenches, trench)
		pos = end
	}

	return trenches
}

func trenchesForY(trenches []Trench, y int) []Trench {
	var row []Trench

	// collect trenches
	for _, trench := range trenches {
		if !trench.Intersects(y) {
			continue
		}

		row = append(row, trench)
	}

	slices.SortFunc(row, func(a, b Trench) int {
		c := a.Start.X - b.Start.X
		if c != 0 {
			return c
		}

		if a.Vertical() {
			return -1
		}
		return +1
	})

	return row
}

func countHole(m map[Pos]string) uint64 {
	start := time.Now()
	rows := splitMap(m)
	fmt.Printf("Splitting took %d\n", time.Since(start))

	count := uint64(0)
	for y, row := range rows {
		fmt.Printf("%d / %d\r", y, len(rows))

		sort.Ints(row)

		areWeInYet := false
		start := 0
		skip := 0

		for _, x := range row {
			if _, ok := m[Pos{X: x, Y: y + 1}]; ok {
				areWeInYet = !areWeInYet
				if areWeInYet {
					start = x + 1
					skip = 0
				} else {
					width := x - start
					width -= skip
					count += uint64(width)
				}
			} else {
				skip++
			}
		}
	}
	fmt.Println()

	return count
}

func splitMap(m map[Pos]string) [][]int {
	var rows [][]int
	for pos := range m {
		if pos.Y >= len(rows) {
			rows = append(rows, slices.Repeat([][]int{nil}, pos.Y-len(rows)+1)...)
		}

		rows[pos.Y] = append(rows[pos.Y], pos.X)
	}
	return rows
}
