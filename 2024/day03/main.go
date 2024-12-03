package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
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

var mul = regexp.MustCompile(`mul\((\d+),(\d+)\)`)

func task1(args []string) error {
	memory, err := loadString(args[0])
	if err != nil {
		return err
	}

	matches := mul.FindAllStringSubmatch(memory, -1)
	sum := 0
	for _, match := range matches {
		a, _ := strconv.Atoi(match[1])
		b, _ := strconv.Atoi(match[2])
		sum += a * b
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}

func loadString(name string) (string, error) {
	p, err := os.ReadFile(name)
	return string(p), err
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

var mulOrDo = regexp.MustCompile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)

func task2(args []string) error {
	memory, err := loadString(args[0])
	if err != nil {
		return err
	}

	matches := mulOrDo.FindAllStringSubmatch(memory, -1)
	sum := 0
	do := true
	for _, match := range matches {
		switch match[0] {
		case "do()":
			do = true
			continue
		case "don't()":
			do = false
			continue
		}

		if do {
			a, _ := strconv.Atoi(match[1])
			b, _ := strconv.Atoi(match[2])
			sum += a * b
		}
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}
