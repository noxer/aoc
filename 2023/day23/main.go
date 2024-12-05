package main

import (
	"fmt"
	"math/bits"
	"os"
	"slices"
	"strings"

	"github.com/noxer/aoc/2024/utils"
)

const (
	StartID = 0
	EndID   = -1
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

func task1(args []string) error {
	m, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	edges := generateGraph(m)
	fmt.Println(edges)

	paths := []Path{
		{
			Current: StartID,
			Length:  0,
			Visited: make(Set[int]),
		},
	}

	maxHike := 0
	for len(paths) > 0 {
		path := paths[0]
		copy(paths, paths[1:])
		paths = paths[:len(paths)-1]

		if path.Current == EndID {
			maxHike = max(maxHike, path.Length)
			continue
		}

		paths = append(paths, path.WalkAll(edges)...)
	}
	fmt.Printf("Max hike: %d\n", maxHike)

	return nil
}

type Path struct {
	Current int
	Length  int
	Visited Set[int]
}

func (p Path) Walk(e Edge) Path {
	p.Current = e.End
	p.Length += e.Length
	p.Visited = p.Visited.Copy()
	p.Visited[e.End] = struct{}{}

	return p
}

func (p Path) WalkAll(edges []Edge) []Path {
	next := make([]Path, 0, 3)

	for _, edge := range edges {
		// Filter other edges
		if p.Current != edge.Start {
			continue
		}

		// Filter visited end nodes
		if _, ok := p.Visited[edge.End]; ok {
			continue
		}

		next = append(next, p.Walk(edge))
	}

	return next
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Copy() Set[T] {
	r := make(Set[T], len(s))
	for k := range s {
		r[k] = struct{}{}
	}
	return r
}

func generateGraph(m []string) []Edge {
	nodes := make(map[Pos]int)
	edges := make([]Edge, 0, 100)

	startPos := Pos{findGap(m[0]), 0}
	nodes[startPos] = StartID
	nodes[Pos{findGap(m[len(m)-1]), len(m) - 1}] = EndID

	toVisit := []Pos{startPos}

	for len(toVisit) > 0 {
		startPos = toVisit[0]
		copy(toVisit, toVisit[1:])
		toVisit = toVisit[:len(toVisit)-1]

		dirs := findOptions(m, startPos, startPos)
		for dir := range dirs.Iterate() {
			pathLength, endPos := followPath(m, startPos.Add(dir.RelativePos()), startPos)
			if endID, ok := nodes[endPos]; ok {
				edges = append(edges, Edge{
					Length: pathLength,
					Start:  nodes[startPos],
					End:    endID,
				})
			} else {
				endID := len(nodes)
				nodes[endPos] = endID

				toVisit = append(toVisit, endPos)

				edges = append(edges, Edge{
					Length: pathLength,
					Start:  nodes[startPos],
					End:    endID,
				})
			}
		}
	}

	return edges
}

func followPath(m []string, current, last Pos) (int, Pos) {
	length := 1

	for {
		dirs := findOptions(m, current, last)
		switch dirs.Count() {
		case 0:
			return length, current
		case 1:
			last = current
			current = current.Add(dirs.RelativePos())
			length++
		default:
			return length, current
		}
	}
}

func findOptions(m []string, current, last Pos) Direction {
	possibleDirs := Direction(0)
	for a, dir := range dirs {
		newPos := current.Add(a)
		if newPos == last ||
			newPos.X < 0 || newPos.X >= len(m[0]) ||
			newPos.Y < 0 || newPos.Y >= len(m) {

			continue
		}

		tile := m[newPos.Y][newPos.X]
		switch {
		case tile == '.',
			tile == '>' && dir == Right,
			tile == '<' && dir == Left,
			tile == 'v' && dir == Down,
			tile == '^' && dir == Up:

			possibleDirs |= dir
		}
	}

	return possibleDirs
}

func findGap(line string) int {
	return strings.Index(line, ".")
}

type Direction byte

func (d Direction) Count() int {
	return bits.OnesCount8(byte(d))
}

func (d Direction) RelativePos() Pos {
	switch d {
	case Up:
		return Pos{Y: -1}
	case Down:
		return Pos{Y: +1}
	case Left:
		return Pos{X: -1}
	case Right:
		return Pos{X: +1}
	}

	return Pos{}
}

func (d Direction) Iterate() func(yield func(Direction) bool) {
	return func(yield func(Direction) bool) {
		for i := range 4 {
			dir := Direction(1 << i)
			if d&dir != 0 {
				if !yield(dir) {
					return
				}
			}
		}
	}
}

const (
	Up Direction = 1 << iota
	Down
	Left
	Right
)

var dirs = map[Pos]Direction{
	{Y: -1}: Up,
	{Y: +1}: Down,
	{X: -1}: Left,
	{X: +1}: Right,
}

type Pos struct {
	X, Y int
}

func (p Pos) Add(o Pos) Pos {
	return Pos{p.X + o.X, p.Y + o.Y}
}

type Node struct {
	ID  int
	Pos Pos
}

type Edge struct {
	Length int
	Start  int
	End    int
}

//////////////////////////////////////////////////////

func generateGraph2(m []string) []Edge {
	nodes := make(map[Pos]int)
	edges := make([]Edge, 0, 100)

	startPos := Pos{findGap(m[0]), 0}
	nodes[startPos] = StartID
	nodes[Pos{findGap(m[len(m)-1]), len(m) - 1}] = EndID

	toVisit := []Pos{startPos}

	for len(toVisit) > 0 {
		startPos = toVisit[0]
		copy(toVisit, toVisit[1:])
		toVisit = toVisit[:len(toVisit)-1]

		dirs := findOptions2(m, startPos, startPos)
		for dir := range dirs.Iterate() {
			pathLength, endPos := followPath2(m, startPos.Add(dir.RelativePos()), startPos)
			if endID, ok := nodes[endPos]; ok {
				edges = append(edges, Edge{
					Length: pathLength,
					Start:  nodes[startPos],
					End:    endID,
				})
			} else {
				endID := len(nodes)
				nodes[endPos] = endID

				toVisit = append(toVisit, endPos)

				edges = append(edges, Edge{
					Length: pathLength,
					Start:  nodes[startPos],
					End:    endID,
				})
			}
		}
	}

	return edges
}

func findOptions2(m []string, current, last Pos) Direction {
	possibleDirs := Direction(0)
	for a, dir := range dirs {
		newPos := current.Add(a)
		if newPos == last ||
			newPos.X < 0 || newPos.X >= len(m[0]) ||
			newPos.Y < 0 || newPos.Y >= len(m) {

			continue
		}

		tile := m[newPos.Y][newPos.X]
		if tile != '#' {
			possibleDirs |= dir
		}
	}

	return possibleDirs
}

func followPath2(m []string, current, last Pos) (int, Pos) {
	length := 1

	for {
		dirs := findOptions2(m, current, last)
		switch dirs.Count() {
		case 0:
			return length, current
		case 1:
			last = current
			current = current.Add(dirs.RelativePos())
			length++
		default:
			return length, current
		}
	}
}

func task2(args []string) error {

	m, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	edges := generateGraph2(m)
	fmt.Println(edges)

	paths := []Path{
		{
			Current: StartID,
			Length:  0,
			Visited: make(Set[int]),
		},
	}

	maxLen := make(map[int]int)
	for len(paths) > 0 {
		path := paths[0]
		copy(paths, paths[1:])
		paths = paths[:len(paths)-1]

		newPaths := path.WalkAll(edges)
		for _, newPath := range newPaths {
			maxLen[newPath.Current] = max(maxLen[newPath.Current], newPath.Length)
			if newPath.Current == EndID {
				fmt.Printf("Found end: %d\n", maxLen[EndID])
			}
		}

		paths = append(paths, newPaths...)
		slices.SortFunc(paths, func(a, b Path) int {
			return b.Length - a.Length
		})
	}
	fmt.Printf("Max hike: %d\n", maxLen[EndID])

	return nil
}
