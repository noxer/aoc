package main

import (
	"fmt"
	"os"
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

func checkVowels(str string) bool {
	count := 0
	for _, b := range []byte(str) {
		switch b {
		case 'a', 'e', 'i', 'o', 'u':
			count++
			if count >= 3 {
				return true
			}
		}
	}

	return false
}

func checkDouble(str string) bool {
	last := byte(0)
	for _, b := range []byte(str) {
		if b == last {
			return true
		}

		last = b
	}

	return false
}

func checkNotContains(str string) bool {
	return !strings.Contains(str, "ab") &&
		!strings.Contains(str, "cd") &&
		!strings.Contains(str, "pq") &&
		!strings.Contains(str, "xy")
}

func check(str string) bool {
	return checkVowels(str) && checkDouble(str) && checkNotContains(str)
}

func task1(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	count := 0
	for _, line := range lines {
		if check(line) {
			count++
		}
	}

	fmt.Printf("Nice: %d\n", count)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func checkTwice(str string) bool {
	for i := range []byte(str[:len(str)-2]) {
		needle := str[i : i+2]
		if strings.Contains(str[i+2:], needle) {
			return true
		}
	}

	return false
}

func checkRepeat(str string) bool {
	for i := range []byte(str[:len(str)-2]) {
		if str[i] == str[i+2] {
			return true
		}
	}

	return false
}

func check2(str string) bool {
	return checkTwice(str) && checkRepeat(str)
}

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	count := 0
	for _, line := range lines {
		nice := check2(line)
		fmt.Printf("%s is nice: %t\n", line, nice)

		if nice {
			count++
		}
	}

	fmt.Printf("Nice: %d\n", count)

	return nil
}
