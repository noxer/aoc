package main

import (
	"fmt"
	"os"
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

func calcPresents(house int) int {
	presents := 0

	for i := range house {
		elf := i + 1
		if house%elf == 0 {
			presents += elf * 10
		}
	}

	return presents
}

func task1(args []string) error {
	const max = 33100000

	for i := range 1_000_000 {
		presents := calcPresents(i + 1)
		fmt.Printf("House %d gets %d presents\n", i+1, presents)

		if presents >= max {
			break
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func calcPresents2(house int) int {
	presents := 0

	for i := range house {
		elf := i + 1
		if house%elf == 0 && house/elf <= 50 {
			presents += elf * 11
		}
	}

	return presents
}

func task2(args []string) error {
	const max = 33100000

	for i := range 1_000_000 {
		presents := calcPresents2(i + 1)
		fmt.Printf("House %d gets %d presents\n", i+1, presents)

		if presents >= max {
			break
		}
	}

	return nil
}
