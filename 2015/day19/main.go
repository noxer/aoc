package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
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

func parseFile(name string) (map[string][]string, string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	m := make(map[string][]string)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}

		left, right, _ := strings.Cut(line, " => ")
		m[left] = append(m[left], right)
	}

	s.Scan()
	return m, s.Text(), s.Err()
}

func replacements(ts map[string][]string, str string) map[string]struct{} {
	set := make(map[string]struct{})
	for i := range str {
		for _, re := range ts[str[i:i+1]] {
			set[str[:i]+re+str[i+1:]] = struct{}{}
		}

		if len(str)-i >= 2 {
			for _, re := range ts[str[i:i+2]] {
				set[str[:i]+re+str[i+2:]] = struct{}{}
			}
		}
	}
	return set
}

func task1(args []string) error {
	ts, mo, err := parseFile(args[0])
	if err != nil {
		return err
	}

	re := replacements(ts, mo)
	fmt.Printf("Distict: %d\n", len(re))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type key struct {
	mo string
	st int
}

var cache = map[key]int{}

func expand(ts map[string][]string, expected, molecule string, steps int) int {
	if molecule == expected {
		fmt.Printf("Found steps: %d\n", steps)
		return steps
	}

	if st, ok := cache[key{molecule, steps}]; ok {
		return st
	}

	bestSteps := math.MaxInt
	for i := range molecule {
		for k, vs := range ts {
			if strings.HasPrefix(molecule[i:], k) {
				for _, v := range vs {
					re := molecule[:i] + v + molecule[i+len(k):]
					bestSteps = min(bestSteps, expand(ts, expected, re, steps+1))
				}
			}
		}
	}

	cache[key{molecule, steps}] = bestSteps

	return bestSteps
}

func invert(m map[string][]string) map[string][]string {
	o := make(map[string][]string)
	for k, vs := range m {
		for _, v := range vs {
			o[v] = append(o[v], k)
		}
	}
	return o
}

func task2(args []string) error {
	ts, mo, err := parseFile(args[0])
	if err != nil {
		return err
	}

	st := invert(ts)

	best := expand(st, "e", mo, 0)
	fmt.Printf("Best: %d (max int: %d)\n", best, math.MaxInt)

	return nil
}
