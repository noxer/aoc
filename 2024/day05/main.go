package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
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

type Rule struct {
	L, R int
}

func (r Rule) Apply(batch Batch) bool {
	a := slices.Index(batch, r.L)
	b := slices.Index(batch, r.R)

	if a < 0 || b < 0 {
		return true
	}

	return a < b
}

func parseRule(line string) Rule {
	r := Rule{}
	left, right, _ := strings.Cut(line, "|")
	r.L, _ = strconv.Atoi(left)
	r.R, _ = strconv.Atoi(right)
	return r
}

type Batch []int

func (b Batch) MiddleElement() int {
	if len(b) == 0 {
		return -1
	}

	return b[len(b)/2]
}

func parseBatch(line string) Batch {
	parts := strings.Split(line, ",")
	b := make(Batch, len(parts))
	for i, part := range parts {
		b[i], _ = strconv.Atoi(part)
	}
	return b
}

func parseFile(name string) ([]Rule, []Batch, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var rules []Rule
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}

		rules = append(rules, parseRule(line))
	}

	var batches []Batch
	for s.Scan() {
		line := s.Text()
		batches = append(batches, parseBatch(line))
	}

	return rules, batches, s.Err()
}

func applyRules(rules []Rule, batch Batch) bool {
	for _, rule := range rules {
		if !rule.Apply(batch) {
			return false
		}
	}
	return true
}

func task1(args []string) error {
	rules, batches, err := parseFile(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	sum := 0
	for _, batch := range batches {
		if applyRules(rules, batch) {
			sum += batch.MiddleElement()
		}
	}

	elapsed := time.Since(start)

	fmt.Printf("Sum: %d, %s\n", sum, elapsed)
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func enforceRules(rules []Rule, batch Batch) {
	slices.SortFunc(batch, func(a, b int) int {
		for _, rule := range rules {
			if rule.L == a && rule.R == b {
				return -1
			}
			if rule.L == b && rule.R == a {
				return +1
			}
		}

		return -1
	})
}

func task2(args []string) error {
	rules, batches, err := parseFile(args[0])
	if err != nil {
		return err
	}

	sum := 0
	for _, batch := range batches {
		if !applyRules(rules, batch) {
			enforceRules(rules, batch)
			sum += batch.MiddleElement()
		}
	}

	fmt.Printf("Sum: %d\n", sum)
	return nil
}
