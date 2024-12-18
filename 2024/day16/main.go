package main

import (
	"container/heap"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/noxer/aoc/2024/utils"
)

func main() {
	var err error

	// task1([]string{"input.txt"})
	// return

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

var (
	West      = utils.Vec{X: -1}
	East      = utils.Vec{X: +1}
	North     = utils.Vec{Y: -1}
	South     = utils.Vec{Y: +1}
	dirs      = [4]utils.Vec{West, East, North, South}
	turnScore = map[Turn]int{
		{West, North}: 1000,
		{North, East}: 1000,
		{East, South}: 1000,
		{South, West}: 1000,

		{North, West}: 1000,
		{East, North}: 1000,
		{South, East}: 1000,
		{West, South}: 1000,

		{North, South}: 2000,
		{South, North}: 2000,
		{East, West}:   2000,
		{West, East}:   2000,
	}
)

type Turn struct {
	From, To utils.Vec
}

type Maze struct {
	data  map[utils.Vec]byte
	size  utils.Vec
	start utils.Vec
	end   utils.Vec
}

func (m *Maze) FindStartAndEnd() {
	for pos, val := range m.data {
		switch val {
		case 'S':
			m.start = pos
			delete(m.data, pos)
		case 'E':
			m.end = pos
			delete(m.data, pos)
		}
	}
}

func (m Maze) Get(pos utils.Vec) byte {
	return m.data[pos]
}

func (m Maze) Wall(pos utils.Vec) bool {
	return m.Get(pos) == '#'
}

func (m Maze) MarkNodes() {
	for pos, val := range m.Iterate() {
		if val != 0 {
			continue
		}

		r := m.Wall(pos.Add(East))
		l := m.Wall(pos.Add(West))
		u := m.Wall(pos.Add(North))
		d := m.Wall(pos.Add(South))

		if l && r && !u && !d || !l && !r && u && d {
			continue
		}

		m.data[pos] = '+'
	}
}

type Edge struct {
	Cost int
	End  *Node
}

type Node struct {
	ID    int
	Edges [4]Edge
}

func (m Maze) Graph() map[utils.Vec]*Node {
	nodes := make(map[utils.Vec]*Node)

	id := 1
	for pos := range m.FilterIterator('+') {
		nodes[pos] = &Node{ID: id}
		id++
	}

	for pos, node := range nodes {
		node.Edges = m.FindEdges(nodes, pos)
	}

	return nodes
}

func (m Maze) FindEdges(nodes map[utils.Vec]*Node, pos utils.Vec) [4]Edge {
	var result [4]Edge

	for i, dir := range dirs {
		nextPos := pos.Add(dir)
		if m.Wall(nextPos) {
			continue
		}

		count := 1
		for m.Get(nextPos) != '+' {
			nextPos = nextPos.Add(dir)
			count++
		}

		result[i] = Edge{
			Cost: count,
			End:  nodes[nextPos],
		}
	}

	return result
}

func (m Maze) Iterate() func(func(utils.Vec, byte) bool) {
	return func(yield func(utils.Vec, byte) bool) {
		for y := range m.size.Y {
			for x := range m.size.X {
				vec := utils.Vec{X: x, Y: y}
				if !yield(vec, m.data[vec]) {
					return
				}
			}
		}
	}
}

func (m Maze) FilterIterator(filter byte) func(func(utils.Vec) bool) {
	return func(yield func(utils.Vec) bool) {
		for pos, val := range m.Iterate() {
			if val != filter {
				continue
			}

			if !yield(pos) {
				return
			}
		}

	}
}

func (m Maze) Print() {
	for y := range m.size.Y {
		for x := range m.size.X {
			if b, ok := m.data[utils.Vec{X: x, Y: y}]; ok {
				fmt.Print(string(b))
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

type Reindeer struct {
	Node      *Node
	Direction utils.Vec
	Seen      []int
	Score     int
}

func (r Reindeer) String() string {
	return fmt.Sprintf("Reindeer with Score %d at Node %d, looking %v, seen: %v", r.Score, r.Node.ID, r.Direction, r.Seen)
}

type ReindeerHeap []Reindeer

func (r *ReindeerHeap) Len() int {
	return len(*r)
}

func (r *ReindeerHeap) Less(i, j int) bool {
	return (*r)[i].Score < (*r)[j].Score
}

func (r *ReindeerHeap) Swap(i, j int) {
	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
}

func (r *ReindeerHeap) Push(x any) {
	rd := x.(Reindeer)
	*r = append(*r, rd)
}

func (r *ReindeerHeap) Pop() any {
	rd := (*r)[r.Len()-1]
	*r = (*r)[:r.Len()-1]
	return rd
}

func (r Reindeer) Move() []Reindeer {
	var results []Reindeer

	r.Seen = append(r.Seen, r.Node.ID)

	for i, edge := range r.Node.Edges {
		if edge.End == nil {
			continue
		}

		if slices.Contains(r.Seen, edge.End.ID) {
			continue
		}

		score := turnScore[Turn{r.Direction, dirs[i]}]
		reindeer := Reindeer{
			Node:      edge.End,
			Direction: dirs[i],
			Seen:      append([]int(nil), r.Seen...),
			Score:     r.Score + score + edge.Cost,
		}

		results = append(results, reindeer)
	}

	return results
}

type Tuple[S, T comparable] struct {
	A S
	B T
}

func task1(args []string) error {
	data, size, err := utils.ReadMapWithSize(args[0], '.')
	if err != nil {
		return err
	}

	m := Maze{
		data: data,
		size: size,
	}

	// m.Print()
	m.FindStartAndEnd()
	// m.Print()
	m.MarkNodes()
	// m.Print()

	graph := m.Graph()
	startNode := graph[m.start]
	endNode := graph[m.end]

	reindeers := &ReindeerHeap{{Node: startNode, Direction: East}}
	heap.Init(reindeers)

	seen := make(map[Tuple[int, utils.Vec]]int)

	score := 999999999999999999
	for reindeers.Len() > 0 {
		fmt.Printf("Length: %d\n", reindeers.Len())

		reindeer := heap.Pop(reindeers).(Reindeer)

		if reindeer.Score > score {
			fmt.Printf("Score %d >= %d, killing reindeer...\n", reindeer.Score, score)
			continue
		}

		if reindeer.Node.ID == endNode.ID {
			fmt.Printf("Reindeer found end with score: %d\n", reindeer.Score)
			score = min(score, reindeer.Score)
			continue
		}

		for _, rd := range reindeer.Move() {
			if best, ok := seen[Tuple[int, utils.Vec]{rd.Node.ID, rd.Direction}]; ok {
				if rd.Score < best {
					seen[Tuple[int, utils.Vec]{rd.Node.ID, rd.Direction}] = rd.Score
					heap.Push(reindeers, rd)
				}
			} else {
				seen[Tuple[int, utils.Vec]{rd.Node.ID, rd.Direction}] = rd.Score
				heap.Push(reindeers, rd)
			}
		}
	}

	fmt.Printf("Score: %d\n", score)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	data, size, err := utils.ReadMapWithSize(args[0], '.')
	if err != nil {
		return err
	}

	start := time.Now()

	m := Maze{
		data: data,
		size: size,
	}

	// m.Print()
	m.FindStartAndEnd()
	// m.Print()
	m.MarkNodes()
	// m.Print()

	graph := m.Graph()
	startNode := graph[m.start]
	endNode := graph[m.end]

	reindeers := &ReindeerHeap{{Node: startNode, Direction: East}}
	heap.Init(reindeers)

	seen := make(map[Tuple[int, utils.Vec]]int)

	score := 999999999999999999
	scoreSeens := make([][]int, 0, 20)
	for reindeers.Len() > 0 {
		// fmt.Printf("Length: %d\n", reindeers.Len())

		reindeer := heap.Pop(reindeers).(Reindeer)

		if reindeer.Score > score {
			// fmt.Printf("Score %d >= %d, killing reindeer...\n", reindeer.Score, score)
			continue
		}

		if reindeer.Node.ID == endNode.ID {
			// fmt.Printf("Reindeer found end with score: %d\n", reindeer.Score)

			if reindeer.Score < score {
				scoreSeens = scoreSeens[:1]
				scoreSeens[0] = reindeer.Seen
				score = reindeer.Score
			} else if reindeer.Score == score {
				scoreSeens = append(scoreSeens, reindeer.Seen)
			}

			continue
		}

		for _, rd := range reindeer.Move() {
			if best, ok := seen[Tuple[int, utils.Vec]{rd.Node.ID, rd.Direction}]; ok {
				if rd.Score <= best {
					seen[Tuple[int, utils.Vec]{rd.Node.ID, rd.Direction}] = rd.Score
					heap.Push(reindeers, rd)
				}
			} else {
				seen[Tuple[int, utils.Vec]{rd.Node.ID, rd.Direction}] = rd.Score
				heap.Push(reindeers, rd)
			}
		}
	}

	// fmt.Println(scoreSeens)

	elapsed := time.Since(start)

	fmt.Printf("Score: %d (%s)\n", countTiles(graph, scoreSeens, endNode.ID), elapsed)

	return nil
}

func searchNode(graph map[utils.Vec]*Node, nodeID int) *Node {
	for _, node := range graph {
		if node.ID == nodeID {
			return node
		}
	}

	return nil
}

func countTiles(graph map[utils.Vec]*Node, seens [][]int, endNodeID int) int {
	seenEdge := make(map[Tuple[int, int]]bool)
	seenNode := make(map[int]bool)
	score := 0

	for _, seen := range seens {
		for i, startID := range seen {
			endID := endNodeID
			if i+1 < len(seen) {
				endID = seen[i+1]
			}

			if seenEdge[Tuple[int, int]{startID, endID}] {
				continue
			}

			start := searchNode(graph, startID)
			for _, edge := range start.Edges {
				if edge.End != nil && edge.End.ID == endID {
					score += edge.Cost

					if seenNode[endID] {
						score--
					} else {
						seenNode[endID] = true
					}

					seenEdge[Tuple[int, int]{startID, endID}] = true
					seenEdge[Tuple[int, int]{endID, startID}] = true
					break
				}
			}
		}
	}

	return score + 1 // add start tile
}

// sOO
// o O
// o O
// o O
// oOO
// o
// o
// e
