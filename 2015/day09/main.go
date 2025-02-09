package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"

	"github.com/noxer/aoc/2015/utils"
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

type Destination struct {
	To   string
	Dist int
}

func loadDistances(name string) (map[string][]Destination, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var (
		from, to string
		dist     int
		graph    = make(map[string][]Destination)
	)
	for s.Scan() {
		fmt.Sscanf(s.Text(), "%s to %s = %d", &from, &to, &dist)
		graph[from] = append(graph[from], Destination{to, dist})
		graph[to] = append(graph[to], Destination{from, dist})
	}

	return graph, s.Err()
}

type Path struct {
	Length   int
	Location string
	Seen     []string
}

func (p Path) Calc(graph map[string][]Destination) []Path {
	here := graph[p.Location]
	var results []Path

	seen := append(p.Seen[0:len(p.Seen):len(p.Seen)], p.Location)

	for _, dest := range here {
		if slices.Contains(p.Seen, dest.To) {
			continue
		}

		results = append(results, Path{
			Length:   p.Length + dest.Dist,
			Location: dest.To,
			Seen:     seen,
		})
	}

	return results
}

func (p Path) Score() int {
	return p.Length
}

func task1(args []string) error {
	graph, err := loadDistances(args[0])
	if err != nil {
		return err
	}

	paths := utils.NewHeap(Path.Score)
	for loc := range graph {
		paths.Push(Path{
			Location: loc,
		})
	}

	shortestPath := math.MaxInt
	for paths.Len() > 0 {
		path := paths.Pop()
		if path.Length >= shortestPath {
			continue
		}

		for _, next := range path.Calc(graph) {
			if len(next.Seen) == len(graph)-1 {
				if next.Length < shortestPath {
					fmt.Printf(
						"Found path: %s, %s with length %d\n",
						strings.Join(next.Seen, ", "),
						next.Location,
						next.Length,
					)
					shortestPath = next.Length
				}
				continue
			}

			paths.Push(next)
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	graph, err := loadDistances(args[0])
	if err != nil {
		return err
	}

	paths := utils.NewHeap(Path.Score)
	for loc := range graph {
		paths.Push(Path{
			Location: loc,
		})
	}

	longestPath := 0
	for paths.Len() > 0 {
		path := paths.Pop()

		for _, next := range path.Calc(graph) {
			if len(next.Seen) == len(graph)-1 {
				if next.Length > longestPath {
					fmt.Printf(
						"Found path: %s, %s with length %d\n",
						strings.Join(next.Seen, ", "),
						next.Location,
						next.Length,
					)
					longestPath = next.Length
				}
				continue
			}

			paths.Push(next)
		}
	}

	return nil
}
