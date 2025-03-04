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

func loadInstructions(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func task1(args []string) error {
	inst, err := loadInstructions(args[0])
	if err != nil {
		return err
	}

	floor := 0
	for _, b := range inst {
		switch b {
		case '(':
			floor++
		case ')':
			floor--
		}
	}

	fmt.Printf("Floor: %d\n", floor)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	inst, err := loadInstructions(args[0])
	if err != nil {
		return err
	}

	floor := 0
	for i, b := range inst {
		switch b {
		case '(':
			floor++
		case ')':
			floor--
		}

		if floor == -1 {
			fmt.Printf("First Instruction: %d\n", i+1)
			break
		}
	}

	return nil
}
