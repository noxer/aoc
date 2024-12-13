package main

import (
	"fmt"
	"os"
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

var directions = []utils.Vec{
	{X: -1}, // left
	{X: +1}, // right
	{Y: -1}, // up
	{Y: +1}, // down
}

type Node struct {
	Type  byte
	Links []*Node
}

type Map [][]Node

func (m Map) ConnectNodes() {
	for y, row := range m {
		for x, node := range row {
			pos := utils.Vec{X: x, Y: y}

			for _, dir := range directions {
				conn := m.Get(pos.Add(dir))
				if conn == nil || conn.Type != node.Type {
					continue
				}

				node.Links = append(node.Links, conn)
			}

			row[x] = node
		}
	}
}

func (m Map) CalculatePrice() int {
	seen := make(map[*Node]struct{})
	price := 0

	for _, row := range m {
		for x := range row {
			node := &row[x]

			if _, ok := seen[node]; ok {
				continue
			}

			area, perimeter := m.calculateRegion(node, seen)
			price += area * perimeter
		}
	}

	return price
}

func (m Map) calculateRegion(node *Node, seen map[*Node]struct{}) (int, int) {
	if _, ok := seen[node]; ok {
		return 0, 0
	}

	area := 1
	perimeter := 4 - len(node.Links)

	seen[node] = struct{}{}
	for _, link := range node.Links {
		a, p := m.calculateRegion(link, seen)
		area += a
		perimeter += p
	}

	return area, perimeter
}

func (m Map) Get(pos utils.Vec) *Node {
	if pos.X < 0 || pos.Y < 0 {
		return nil
	}

	if pos.Y >= len(m) {
		return nil
	}
	row := m[pos.Y]

	if pos.X >= len(row) {
		return nil
	}

	return &row[pos.X]
}

func parseNodes(line string) []Node {
	nodes := make([]Node, len(line))
	for i, b := range []byte(line) {
		nodes[i] = Node{Type: b}
	}
	return nodes
}

func task1(args []string) error {
	r, err := utils.ReadLinesTransform(args[0], parseNodes)
	if err != nil {
		return err
	}

	start := time.Now()

	m := Map(r)
	m.ConnectNodes()
	price := m.CalculatePrice()

	elapsed := time.Since(start)

	fmt.Printf("Price: %d (%s)\n", price, elapsed)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

var (
	turnLeft = []int{
		3, // left
		2, // right
		0, // up
		1, // down
	}
)

type Node2 struct {
	Type  byte
	Links [4]*Node2
}

type Map2 [][]Node2

func (m Map2) ConnectNodes() {
	for y, row := range m {
		for x, node := range row {
			pos := utils.Vec{X: x, Y: y}

			for i, dir := range directions {
				conn := m.Get(pos.Add(dir))
				if conn == nil || conn.Type != node.Type {
					continue
				}

				node.Links[i] = conn
			}

			row[x] = node
		}
	}
}

func (m Map2) CalculatePrice() int {
	seen := make(map[*Node2]struct{})
	price := 0

	for _, row := range m {
		for x := range row {
			node := &row[x]

			if _, ok := seen[node]; ok {
				continue
			}

			area, perimeter := m.calculateRegion(node, seen)
			price += area * perimeter

			// fmt.Printf("%c: %d x %d = %d\n", node.Type, area, perimeter, area*perimeter)
		}
	}

	return price
}

func (m Map2) calculateRegion(node *Node2, seen map[*Node2]struct{}) (int, int) {
	if _, ok := seen[node]; ok {
		return 0, 0
	}

	area := 1
	perimeter := m.checkPerimeter(node)

	seen[node] = struct{}{}
	for _, link := range node.Links {
		if link == nil {
			continue
		}

		a, p := m.calculateRegion(link, seen)
		area += a
		perimeter += p
	}

	return area, perimeter
}

func (m Map2) checkPerimeter(node *Node2) int {
	perimeter := 0
	for i := range 4 {
		next := node.Links[i]
		if next == nil && !m.hasPerimeter(node.Links[turnLeft[i]], i) {
			perimeter++
		}
	}
	return perimeter
}

func (m Map2) hasPerimeter(node *Node2, lookDir int) bool {
	return node != nil && node.Links[lookDir] == nil
}

func (m Map2) Get(pos utils.Vec) *Node2 {
	if pos.X < 0 || pos.Y < 0 {
		return nil
	}

	if pos.Y >= len(m) {
		return nil
	}
	row := m[pos.Y]

	if pos.X >= len(row) {
		return nil
	}

	return &row[pos.X]
}

func parseNodes2(line string) []Node2 {
	nodes := make([]Node2, len(line))
	for i, b := range []byte(line) {
		nodes[i] = Node2{Type: b}
	}
	return nodes
}

func task2(args []string) error {
	r, err := utils.ReadLinesTransform(args[0], parseNodes2)
	if err != nil {
		return err
	}

	start := time.Now()

	m := Map2(r)
	m.ConnectNodes()
	price := m.CalculatePrice()

	elapsed := time.Since(start)

	fmt.Printf("Price: %d (%s)\n", price, elapsed)

	return nil
}
