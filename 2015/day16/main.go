package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

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

type Sue map[string]int

func (s Sue) Matches(other Sue) bool {
	for name, value := range s {
		if other[name] != value {
			return false
		}
	}

	return true
}

func parseAunt(line string) Sue {
	sue := make(Sue)

	_, propStr, _ := strings.Cut(line, ": ")
	for _, prop := range strings.Split(propStr, ", ") {
		name, valueStr, _ := strings.Cut(prop, ": ")
		value, _ := strconv.Atoi(valueStr)

		sue[name] = value
	}

	return sue
}

func task1(args []string) error {
	aunts, err := utils.ReadLinesTransform(args[0], parseAunt)
	if err != nil {
		return err
	}

	reference := Sue{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	for i, aunt := range aunts {
		if aunt.Matches(reference) {
			fmt.Printf("Sue #%d matches the reference!\n", i+1)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (s Sue) Matches2(other Sue) bool {
	for name, value := range s {
		switch name {
		case "cats", "trees":
			if other[name] >= value {
				return false
			}
		case "pomeranians", "goldfish":
			if other[name] <= value {
				return false
			}
		default:
			if other[name] != value {
				return false
			}
		}
	}

	return true
}

func task2(args []string) error {
	aunts, err := utils.ReadLinesTransform(args[0], parseAunt)
	if err != nil {
		return err
	}

	reference := Sue{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}

	for i, aunt := range aunts {
		if aunt.Matches2(reference) {
			fmt.Printf("Sue #%d matches the reference!\n", i+1)
		}
	}

	return nil
}
