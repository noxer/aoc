package main

import (
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

func countUnescaped(str string) int {
	count := 0

	for i := 0; i < len(str); i++ {
		count++

		if str[i] == '\\' {
			switch str[i+1] {
			case '\\', '"':
				i += 1
			case 'x':
				i += 3
			}
		}
	}

	return count - 2
}

func task1(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	sum := 0
	for _, line := range lines {
		fmt.Printf("%s -> %d\n", line, countUnescaped(line))
		sum += len(line) - countUnescaped(line)
	}

	fmt.Printf("Difference: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func countEscaped(str string) int {
	count := 0

	for _, b := range []byte(str) {
		count++

		switch b {
		case '\\', '"':
			count++
		}
	}

	return count + 2
}

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	sum := 0
	for _, line := range lines {
		// fmt.Printf("%s -> %d\n", line, countUnescaped(line))
		sum += countEscaped(line) - len(line)
	}

	fmt.Printf("Difference: %d\n", sum)

	return nil
}
