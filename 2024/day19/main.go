package main

import (
	"bufio"
	"fmt"
	"os"
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

func parseTowels(name string) (utils.Set[string], []string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	if !s.Scan() {
		return nil, nil, s.Err()
	}
	available := utils.SetFromSlice(strings.Split(s.Text(), ", "))

	var patterns []string
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}

		patterns = append(patterns, line)
	}

	return available, patterns, s.Err()
}

func findCombination(towels utils.Set[string], pattern string) bool {
	if pattern == "" {
		return true
	}

	for i := 1; i <= len(pattern); i++ {
		part := pattern[:i]

		// fmt.Printf("Trying to find pattern %s...\n", part)

		if !towels.Has(part) {
			continue
		}

		if findCombination(towels, pattern[i:]) {
			return true
		}
	}

	return false
}

func task1(args []string) error {
	towels, patterns, err := parseTowels(args[0])
	if err != nil {
		return err
	}

	counter := 0
	for _, pattern := range patterns {
		if findCombination(towels, pattern) {
			counter++
		}
	}

	fmt.Printf("Found %d valid patterns\n", counter)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func findAllCombinations(towels utils.Set[string], pattern string, maxTowel int, cache map[string]int) int {
	if pattern == "" {
		return 1
	}

	if sum, ok := cache[pattern]; ok {
		return sum
	}

	sum := 0

	for i := 1; i <= len(pattern); i++ {
		if i > maxTowel {
			break
		}

		part := pattern[:i]

		// fmt.Printf("Trying to find pattern %s...\n", part)

		if !towels.Has(part) {
			continue
		}

		sum += findAllCombinations(towels, pattern[i:], maxTowel, cache)
	}

	cache[pattern] = sum

	return sum
}

func task2(args []string) error {
	towels, patterns, err := parseTowels(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	maxTowel := 0
	for towel := range towels {
		maxTowel = max(maxTowel, len(towel))
	}

	cache := make(map[string]int)

	counter := 0
	for _, pattern := range patterns {
		// fmt.Printf("Checking pattern %d: %s\n", i, pattern)

		counter += findAllCombinations(towels, pattern, maxTowel, cache)
	}

	elapsed := time.Since(start)

	fmt.Printf("Found %d valid patterns in %s\n", counter, elapsed)

	return nil
}
