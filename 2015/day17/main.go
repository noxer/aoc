package main

import (
	"fmt"
	"math"
	"os"
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

func fill(containers []int, remaining int) int {
	if remaining == 0 {
		return 1
	}
	if len(containers) == 0 {
		return 0
	}

	n := fill(containers[1:], remaining)

	if remaining >= containers[0] {
		n += fill(containers[1:], remaining-containers[0])
	}

	return n
}

func task1(args []string) error {
	containers, err := utils.ReadLinesTransform(args[0], func(line string) int {
		n, _ := strconv.Atoi(line)
		return n
	})
	if err != nil {
		return err
	}

	n := fill(containers, 150)
	fmt.Printf("Combinations: %d\n", n)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func minCont(containers []int, remaining, used int) int {
	if remaining == 0 {
		return used
	}
	if len(containers) == 0 {
		return math.MaxInt
	}

	n := minCont(containers[1:], remaining, used)

	if remaining >= containers[0] {
		n = min(n, minCont(containers[1:], remaining-containers[0], used+1))
	}

	return n
}

func matchCont(containers []int, remaining, used, match int) int {
	if remaining == 0 {
		if used == match {
			return 1
		}
		return 0
	}
	if len(containers) == 0 {
		return 0
	}

	n := matchCont(containers[1:], remaining, used, match)

	if remaining >= containers[0] {
		n += matchCont(containers[1:], remaining-containers[0], used+1, match)
	}

	return n
}

func task2(args []string) error {
	containers, err := utils.ReadLinesTransform(args[0], func(line string) int {
		n, _ := strconv.Atoi(line)
		return n
	})
	if err != nil {
		return err
	}

	min := minCont(containers, 150, 0)
	fmt.Printf("Min: %d\n", min)

	n := matchCont(containers, 150, 0, min)
	fmt.Printf("Match: %d\n", n)

	return nil
}
