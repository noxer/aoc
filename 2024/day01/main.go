package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
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

	var left, right []int
	for _, line := range lines {
		a, b, _ := strings.Cut(line, "   ")
		ai, _ := strconv.Atoi(a)
		bi, _ := strconv.Atoi(b)

		left = append(left, ai)
		right = append(right, bi)
	}

	sort.Ints(left)
	sort.Ints(right)

	sum := 0
	for i, l := range left {
		r := right[i]
		diff := max(l, r) - min(l, r)
		sum += diff
	}

	fmt.Printf("Difference: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	var left []int
	right := make(map[int]int)
	for _, line := range lines {
		a, b, _ := strings.Cut(line, "   ")
		ai, _ := strconv.Atoi(a)
		bi, _ := strconv.Atoi(b)

		left = append(left, ai)
		right[bi]++
	}

	sum := 0
	for _, l := range left {
		times := right[l]
		score := l * times
		sum += score
	}

	fmt.Printf("Score: %d\n", sum)

	return nil
}
