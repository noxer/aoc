package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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

func isLock(pattern []string) bool {
	return pattern[0] == "#####"
}

func parseKey(pattern []string) Key {
	slices.Reverse(pattern)
	return Key(parseLock(pattern))
}

func parseLock(pattern []string) Lock {
	var lock [5]int

	for i, row := range pattern {
		for j, b := range []byte(row) {
			if b == '#' {
				lock[j] = i
			}
		}
	}

	return lock
}

type Lock [5]int
type Key [5]int

func parseLocksAndKeys(name string) ([]Lock, []Key, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	pattern := make([]string, 0, 7)
	var locks []Lock
	var keys []Key

	for {
		for s.Scan() {
			line := s.Text()
			if line == "" {
				break
			}

			pattern = append(pattern, line)
		}

		if len(pattern) == 0 {
			break
		}

		if isLock(pattern) {
			locks = append(locks, parseLock(pattern))
		} else {
			keys = append(keys, parseKey(pattern))
		}

		pattern = pattern[:0]
	}

	return locks, keys, s.Err()
}

func keyFits(lock Lock, key Key) bool {
	for i, n := range lock {
		if key[i]+n > 5 {
			return false
		}
	}

	return true
}

func task1(args []string) error {
	locks, keys, err := parseLocksAndKeys(args[0])
	if err != nil {
		return err
	}

	fmt.Println(locks)
	fmt.Println(keys)

	count := 0
	for _, lock := range locks {
		for _, key := range keys {
			if keyFits(lock, key) {
				count++
			}
		}
	}

	fmt.Printf("%d combinations\n", count)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	return nil
}
