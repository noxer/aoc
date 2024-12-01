package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/noxer/aoc/2023/utils"
)

func main() {
	// task2([]string{"example.txt"})
	// return

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

type Part map[string]int

func (p Part) Sum() int {
	sum := 0
	for _, n := range p {
		sum += n
	}
	return sum
}

func parsePart(line string) Part {
	part := make(Part)
	content := line[1 : len(line)-1]
	attrs := strings.Split(content, ",")

	for _, attr := range attrs {
		key, rawValue, _ := strings.Cut(attr, "=")
		value, _ := strconv.Atoi(rawValue)
		part[key] = value
	}

	return part
}

func task1(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	insts := make(map[string][]func(Part) string)
	for i, line := range lines {
		if line == "" {
			lines = lines[i+1:]
			break
		}

		label, cmds := parseInstructions(line)
		insts[label] = cmds
	}

	var parts []Part
	for _, line := range lines {
		parts = append(parts, parsePart(line))
	}

	sum := 0
	for _, part := range parts {
		if runPart(insts, part) {
			sum += part.Sum()
		}
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}

func runPart(insts map[string][]func(Part) string, part Part) bool {
	label := "in"

	for {
		cmds := insts[label]

		for _, cmd := range cmds {
			res := cmd(part)
			if res != "" {
				label = res
				break
			}
		}

		switch label {
		case "A":
			return true
		case "R":
			return false
		}
	}
}

func parseInstructions(line string) (string, []func(Part) string) {
	label, rawCommands, _ := strings.Cut(line, "{")
	rawCommands = rawCommands[:len(rawCommands)-1]
	commands := strings.Split(rawCommands, ",")

	var cmds []func(Part) string
	for _, command := range commands {
		cond, lab, found := strings.Cut(command, ":")
		if !found {
			cmds = append(cmds, Always(command))
			continue
		}

		cmds = append(cmds, parseCommand(cond, lab))
	}

	return label, cmds
}

func parseCommand(cond, label string) func(p Part) string {
	attr, rawValue, found := strings.Cut(cond, ">")
	if found {
		value, _ := strconv.Atoi(rawValue)
		return GreaterThan(value, attr, label)
	}

	attr, rawValue, _ = strings.Cut(cond, "<")
	value, _ := strconv.Atoi(rawValue)
	return LessThan(value, attr, label)
}

func GreaterThan(n int, attr, label string) func(Part) string {
	// fmt.Println("GreaterThan", n, attr, label)

	return func(p Part) string {
		if p[attr] > n {
			return label
		}
		return ""
	}
}

func LessThan(n int, attr, label string) func(p Part) string {
	// fmt.Println("LessThan", n, attr, label)

	return func(p Part) string {
		if p[attr] < n {
			return label
		}
		return ""
	}
}

func Always(label string) func(p Part) string {
	// fmt.Println("Always", label)

	return func(_ Part) string {
		return label
	}
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	filters := make(map[string][]Filter)
	for _, line := range lines {
		if line == "" {
			break
		}

		label, cmds := parseFilters(line)
		filters[label] = cmds
	}

	count := 0
	accept := func(part RangePart) {
		count += part.Count()
	}

	walkTree(filters, "in", RangePart{
		X: Range{A: 1, B: 4000},
		M: Range{A: 1, B: 4000},
		A: Range{A: 1, B: 4000},
		S: Range{A: 1, B: 4000},
	}, accept)

	fmt.Printf("Count: %d\n", count)
	return nil
}

func walkTree(filters map[string][]Filter, current string, part RangePart, accept func(RangePart)) {
	if current == "R" {
		return
	}
	if current == "A" {
		accept(part)
	}

	currentFilters := filters[current]
	var recursePart RangePart
	for _, filter := range currentFilters {
		recursePart, part = filter.Apply(part)
		if recursePart.Valid() {
			walkTree(filters, filter.Label, recursePart, accept)
		}

		if !part.Valid() {
			break
		}
	}
}

type RangePart struct {
	X, M, A, S Range
}

func (p RangePart) Valid() bool {
	return p.X.Valid() && p.M.Valid() && p.A.Valid() && p.S.Valid()
}

func (p RangePart) Count() int {
	return p.X.Size() * p.M.Size() * p.A.Size() * p.S.Size()
}

type Range struct {
	A, B int
}

func (r Range) Size() int {
	return r.B - r.A + 1
}

func (r Range) Valid() bool {
	return r.Size() > 0
}

type Filter struct {
	Key   string
	Value int
	Op    string
	Label string
}

func (f Filter) Apply(rp RangePart) (a, b RangePart) {
	if f.Op == "" {
		return rp, rp
	}

	a = rp
	b = rp

	if f.Op == ">" {
		switch f.Key {
		case "x":
			a.X.A = max(a.X.A, f.Value+1)
			b.X.B = min(a.X.B, f.Value)
		case "m":
			a.M.A = max(a.M.A, f.Value+1)
			b.M.B = min(a.M.B, f.Value)
		case "a":
			a.A.A = max(a.A.A, f.Value+1)
			b.A.B = min(a.A.B, f.Value)
		case "s":
			a.S.A = max(a.S.A, f.Value+1)
			b.S.B = min(a.S.B, f.Value)
		}
	} else {
		switch f.Key {
		case "x":
			a.X.B = min(a.X.B, f.Value-1)
			b.X.A = max(b.X.A, f.Value)
		case "m":
			a.M.B = min(a.M.B, f.Value-1)
			b.M.A = max(b.M.A, f.Value)
		case "a":
			a.A.B = min(a.A.B, f.Value-1)
			b.A.A = max(b.A.A, f.Value)
		case "s":
			a.S.B = min(a.S.B, f.Value-1)
			b.S.A = max(b.S.A, f.Value)
		}
	}

	return a, b
}

func parseFilters(line string) (string, []Filter) {
	label, rawCommands, _ := strings.Cut(line, "{")
	rawCommands = rawCommands[:len(rawCommands)-1]
	commands := strings.Split(rawCommands, ",")

	var filters []Filter
	for _, command := range commands {
		cond, lab, found := strings.Cut(command, ":")
		if !found {
			filters = append(filters, Filter{Label: command})
			continue
		}

		filters = append(filters, parseFilter(cond, lab))
	}

	return label, filters
}

func parseFilter(cond, label string) Filter {
	attr, rawValue, found := strings.Cut(cond, ">")
	if found {
		value, _ := strconv.Atoi(rawValue)
		return Filter{Key: attr, Value: value, Op: ">", Label: label}
	}

	attr, rawValue, _ = strings.Cut(cond, "<")
	value, _ := strconv.Atoi(rawValue)
	return Filter{Key: attr, Value: value, Op: "<", Label: label}
}
