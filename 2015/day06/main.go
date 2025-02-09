package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

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

type Cmd byte

const (
	TurnOn Cmd = iota
	TurnOff
	Toggle
)

type Inst struct {
	A, B    utils.Vec
	Command Cmd
}

var matchInst = regexp.MustCompile(`^(turn on|turn off|toggle) (\d+),(\d+) through (\d+),(\d+)$`)

func parseInst(line string) Inst {
	matches := matchInst.FindStringSubmatch(line)
	if len(matches) != 6 {
		panic(line)
	}

	inst := Inst{}

	switch matches[1] {
	case "turn on":
		inst.Command = TurnOn
	case "turn off":
		inst.Command = TurnOff
	case "toggle":
		inst.Command = Toggle
	}

	inst.A.X, _ = strconv.Atoi(matches[2])
	inst.A.Y, _ = strconv.Atoi(matches[3])

	inst.B.X, _ = strconv.Atoi(matches[4])
	inst.B.Y, _ = strconv.Atoi(matches[5])

	return inst
}

func (i Inst) Apply(grid *[1000][1000]bool) {
	for y := i.A.Y; y <= i.B.Y; y++ {
		for x := i.A.X; x <= i.B.X; x++ {
			switch i.Command {
			case TurnOn:
				(*grid)[y][x] = true
			case TurnOff:
				(*grid)[y][x] = false
			case Toggle:
				(*grid)[y][x] = !(*grid)[y][x]
			}
		}
	}
}

func task1(args []string) error {
	insts, err := utils.ReadLinesTransform(args[0], parseInst)
	if err != nil {
		return err
	}

	var grid [1000][1000]bool
	for _, inst := range insts {
		inst.Apply(&grid)
	}

	count := 0
	for y := range 1000 {
		for x := range 1000 {
			if grid[y][x] {
				count++
			}
		}
	}

	fmt.Printf("Lit lights: %d\n", count)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (i Inst) Apply2(grid *[1000][1000]int) {
	for y := i.A.Y; y <= i.B.Y; y++ {
		for x := i.A.X; x <= i.B.X; x++ {
			switch i.Command {
			case TurnOn:
				(*grid)[y][x]++
			case TurnOff:
				(*grid)[y][x]--
				if (*grid)[y][x] < 0 {
					(*grid)[y][x] = 0
				}
			case Toggle:
				(*grid)[y][x] += 2
			}
		}
	}
}

func task2(args []string) error {
	insts, err := utils.ReadLinesTransform(args[0], parseInst)
	if err != nil {
		return err
	}

	var grid [1000][1000]int
	for _, inst := range insts {
		inst.Apply2(&grid)
	}

	count := 0
	for y := range 1000 {
		for x := range 1000 {
			count += grid[y][x]
		}
	}

	fmt.Printf("Total brightness: %d\n", count)

	return nil

}
