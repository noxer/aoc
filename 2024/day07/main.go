package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

type Equation struct {
	Result int
	Values []int
}

func parseEquation(line string) Equation {
	parts := strings.Fields(line)
	result, _ := strconv.Atoi(strings.TrimSuffix(parts[0], ":"))
	values := make([]int, len(parts)-1)
	for i, str := range parts[1:] {
		values[i], _ = strconv.Atoi(str)
	}

	return Equation{
		Result: result,
		Values: values,
	}
}

func (e Equation) HasSolution() bool {
	return tryOps(e.Values[0], e.Values[1:], func(i int) bool {
		return i == e.Result
	})
}

func tryOps(current int, values []int, check func(int) bool) bool {
	if len(values) == 0 {
		return check(current)
	}

	sum := current + values[0]
	if tryOps(sum, values[1:], check) {
		return true
	}

	prod := current * values[0]
	return tryOps(prod, values[1:], check)
}

func task1(args []string) error {
	equations, err := utils.ReadLinesTransform(args[0], parseEquation)
	if err != nil {
		return err
	}

	start := time.Now()

	sum := 0
	for _, equation := range equations {
		if equation.HasSolution() {
			sum += equation.Result
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("Result: %d (%s)\n", sum, elapsed)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (e Equation) HasSolution2() bool {
	return tryOps2(e.Values[0], e.Values[1:], func(i int) bool {
		return i == e.Result
	})
}

func tryOps2(current int, values []int, check func(int) bool) bool {
	if len(values) == 0 {
		return check(current)
	}

	sum := current + values[0]
	if tryOps2(sum, values[1:], check) {
		return true
	}

	prod := current * values[0]
	if tryOps2(prod, values[1:], check) {
		return true
	}

	cat := concat(current, values[0])
	return tryOps2(cat, values[1:], check)
}

func concat(a, b int) int {
	res, _ := strconv.Atoi(strconv.Itoa(a) + strconv.Itoa(b))
	return res
}

func task2(args []string) error {
	equations, err := utils.ReadLinesTransform(args[0], parseEquation)
	if err != nil {
		return err
	}

	start := time.Now()

	sum := 0
	for _, equation := range equations {
		if equation.HasSolution2() {
			sum += equation.Result
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("Result: %d (%s)\n", sum, elapsed)

	return nil
}
