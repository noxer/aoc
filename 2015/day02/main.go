package main

import (
	"bufio"
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

type Present struct {
	L, W, H int
}

func (p Present) WrappingPaper() int {
	a := p.L * p.W
	b := p.W * p.H
	c := p.H * p.L

	sum := 2*a + 2*b + 2*c
	sum += min(a, b, c)

	return sum
}

func parsePresents(name string) ([]Present, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var presents []Present

	for s.Scan() {
		line := s.Text()
		p := Present{}
		_, err = fmt.Sscanf(line, "%dx%dx%d", &p.L, &p.W, &p.H)
		if err != nil {
			return nil, err
		}

		presents = append(presents, p)
	}

	return presents, s.Err()
}

func task1(args []string) error {
	presents, err := parsePresents(args[0])
	if err != nil {
		return err
	}

	sum := 0
	for _, p := range presents {
		sum += p.WrappingPaper()
	}

	fmt.Printf("We need %d sqf of wrapping paper!\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (p Present) Ribbon() int {
	a := 2*p.L + 2*p.W
	b := 2*p.W + 2*p.H
	c := 2*p.H + 2*p.L

	return min(a, b, c) + p.L*p.W*p.H
}

func task2(args []string) error {
	presents, err := parsePresents(args[0])
	if err != nil {
		return err
	}

	sum := 0
	for _, p := range presents {
		sum += p.Ribbon()
	}

	fmt.Printf("We need %d f of ribbon!\n", sum)

	return nil
}
