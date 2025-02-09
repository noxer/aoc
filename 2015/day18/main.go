package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/noxer/aoc/2015/utils"
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

type Grid map[utils.Vec]struct{}

func (g Grid) CountNeighbors(pos utils.Vec) int {
	off := utils.Vec{}
	count := 0

	for off.X = -1; off.X <= 1; off.X++ {
		for off.Y = -1; off.Y <= 1; off.Y++ {
			if off.X == 0 && off.Y == 0 {
				continue
			}

			if _, ok := g[pos.Add(off)]; ok {
				count++
			}
		}
	}

	return count
}

func (g Grid) NextCell(pos utils.Vec) bool {
	_, alive := g[pos]
	if !alive {
		return g.CountNeighbors(pos) == 3
	}

	switch g.CountNeighbors(pos) {
	case 2, 3:
		return true
	default:
		return false
	}
}

func (g Grid) Next() Grid {
	next := make(Grid)
	pos := utils.Vec{}

	for pos.X = 0; pos.X < 100; pos.X++ {
		for pos.Y = 0; pos.Y < 100; pos.Y++ {
			if g.NextCell(pos) {
				next[pos] = struct{}{}
			}
		}
	}

	return next
}

func loadGrid(name string) (Grid, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	grid := make(Grid)
	pos := utils.Vec{}
	for s.Scan() {
		for x, b := range s.Bytes() {
			if b != '#' {
				continue
			}

			pos.X = x
			grid[pos] = struct{}{}
		}

		pos.Y++
	}

	return grid, nil
}

func task1(args []string) error {
	grid, err := loadGrid(args[0])
	if err != nil {
		return err
	}

	for range 100 {
		grid = grid.Next()
	}

	fmt.Printf("%d lights are on\n", len(grid))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (g Grid) Next2() Grid {
	next := make(Grid)
	next[utils.Vec{0, 0}] = struct{}{}
	next[utils.Vec{0, 99}] = struct{}{}
	next[utils.Vec{99, 0}] = struct{}{}
	next[utils.Vec{99, 99}] = struct{}{}

	pos := utils.Vec{}

	for pos.X = 0; pos.X < 100; pos.X++ {
		for pos.Y = 0; pos.Y < 100; pos.Y++ {
			if g.NextCell(pos) {
				next[pos] = struct{}{}
			}
		}
	}

	return next
}

func task2(args []string) error {
	grid, err := loadGrid(args[0])
	if err != nil {
		return err
	}

	grid[utils.Vec{0, 0}] = struct{}{}
	grid[utils.Vec{0, 99}] = struct{}{}
	grid[utils.Vec{99, 0}] = struct{}{}
	grid[utils.Vec{99, 99}] = struct{}{}

	for range 100 {
		grid = grid.Next2()
	}

	fmt.Printf("%d lights are on\n", len(grid))

	return nil
}
