package main

import (
	"bufio"
	"fmt"
	"iter"
	"maps"
	"os"
	"regexp"
	"slices"
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

var matchHappiness = regexp.MustCompile(`^([A-Za-z]+) would (gain|lose) ([0-9]+) happiness units by sitting next to ([A-Za-z]+).$`)

type Friends map[string]map[string]int

func parseHappiness(f Friends, line string) {
	match := matchHappiness.FindStringSubmatch(line)
	if len(match) != 5 {
		fmt.Printf("Couldn't match line %s\n", line)
		return
	}

	happiness, _ := strconv.Atoi(match[3])
	if match[2] == "lose" {
		happiness = -happiness
	}

	if f[match[1]] == nil {
		f[match[1]] = make(map[string]int)
	}
	f[match[1]][match[4]] = happiness
}

func parseFile(name string) (Friends, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	friends := make(Friends)
	for s.Scan() {
		parseHappiness(friends, s.Text())
	}

	return friends, s.Err()
}

func calculateHappiness(friends Friends, table []string) int {
	happiness := 0

	for i, name := range table[:len(table)-1] {
		happiness += friends[name][table[i+1]]
		happiness += friends[table[i+1]][name]
	}

	happiness += friends[table[0]][table[len(table)-1]]
	happiness += friends[table[len(table)-1]][table[0]]

	return happiness
}

func permute[S ~[]T, T any](sl S) iter.Seq[S] {
	return func(yield func(S) bool) {
		permuteRecursive(sl, sl, yield)
	}
}

func permuteRecursive[S ~[]T, T any](full, part S, yield func(S) bool) bool {
	if len(part) <= 1 {
		return yield(full)
	}

	for i := range part {
		part[0], part[i] = part[i], part[0]
		if !permuteRecursive(full, part[1:], yield) {
			return false
		}
		part[0], part[i] = part[i], part[0]
	}

	return true
}

func task1(args []string) error {
	friends, err := parseFile(args[0])
	if err != nil {
		return err
	}

	table := slices.Collect(maps.Keys(friends))

	bestHappiness := 0
	for table := range permute(table) {
		fmt.Println(table)
		bestHappiness = max(bestHappiness, calculateHappiness(friends, table))
	}

	fmt.Printf("Happiness: %d\n", bestHappiness)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	friends, err := parseFile(args[0])
	if err != nil {
		return err
	}

	table := slices.Collect(maps.Keys(friends))
	table = append(table, "Tim")

	bestHappiness := 0
	for table := range permute(table) {
		fmt.Println(table)
		bestHappiness = max(bestHappiness, calculateHappiness(friends, table))
	}

	fmt.Printf("Happiness: %d\n", bestHappiness)

	return nil
}
