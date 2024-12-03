package main

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

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

type Report []int

func (r Report) Safe() bool {
	return checkAsc(r) || checkDesc(r)
}

func checkAsc(r Report) bool {
	elem := r[0]
	for _, n := range r[1:] {
		if n <= elem || n > elem+3 {
			return false
		}
		elem = n
	}
	return true
}

func checkDesc(r Report) bool {
	elem := r[0]
	for _, n := range r[1:] {
		if n >= elem || n < elem-3 {
			return false
		}
		elem = n
	}
	return true
}

func task1(args []string) error {
	reports, err := utils.ReadLinesTransform(args[0], func(line string) Report {
		parts := strings.Fields(line)
		report := make(Report, len(parts))
		for i, p := range parts {
			report[i], _ = strconv.Atoi(p)
		}
		return report
	})
	if err != nil {
		return err
	}

	count := 0
	for _, report := range reports {
		if report.Safe() {
			count++
		}
	}

	fmt.Printf("Safe reports: %d\n", count)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (r Report) SafeDampened() bool {
	if checkAscDampened(r) || checkDescDampened(r) {
		return true
	}

	slices.Reverse(r)

	return checkAscDampened(r) || checkDescDampened(r)
}

func checkAscDampened(r Report) bool {
	elem := r[0]
	for i, n := range r[1:] {
		if n <= elem || n > elem+3 {
			return checkAsc(deleteCopy(r, i+1))
		}
		elem = n
	}
	return true
}

func checkDescDampened(r Report) bool {
	elem := r[0]
	for i, n := range r[1:] {
		if n >= elem || n < elem-3 {
			return checkDesc(deleteCopy(r, i+1))
		}
		elem = n
	}
	return true
}

func deleteCopy(r Report, idx int) Report {
	out := make(Report, 0, len(r)-1)
	for i, n := range r {
		if i == idx {
			continue
		}
		out = append(out, n)
	}
	return out
}

func task2(args []string) error {
	reports, err := utils.ReadLinesTransform(args[0], func(line string) Report {
		parts := strings.Fields(line)
		report := make(Report, len(parts))
		for i, p := range parts {
			report[i], _ = strconv.Atoi(p)
		}
		return report
	})
	if err != nil {
		return err
	}

	count := 0
	for _, report := range reports {
		if report.SafeDampened() {
			count++
		}
	}

	fmt.Printf("Safe reports: %d\n", count)

	return nil
}
