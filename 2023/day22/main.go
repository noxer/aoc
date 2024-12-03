package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/noxer/aoc/2023/utils"
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

type Cube struct {
	X, Y, Z int
}

func parseCube(str string) Cube {
	parts := strings.Split(str, ",")
	p := Cube{}

	p.X, _ = strconv.Atoi(parts[0])
	p.Y, _ = strconv.Atoi(parts[1])
	p.Z, _ = strconv.Atoi(parts[2])

	return p
}

type Brick struct {
	Start, End Cube
}

func (b Brick) Contains(cube Cube) bool {
	checkX := cube.X >= b.Start.X && cube.X <= b.End.X
	checkY := cube.Y >= b.Start.Y && cube.Y <= b.End.Y
	checkZ := cube.Z >= b.Start.Z && cube.Z <= b.End.Z

	return checkX && checkY && checkZ
}

func (b Brick) MoveDown() Brick {
	b.Start.Z--
	b.End.Z--

	return b
}

func (b Brick) Iterate() func(yield func(Cube) bool) {
	return func(yield func(Cube) bool) {
		for x := b.Start.X; x <= b.End.X; x++ {
			for y := b.Start.Y; y <= b.End.Y; y++ {
				for z := b.Start.Z; z <= b.End.Z; z++ {
					if !yield(Cube{x, y, z}) {
						return
					}
				}
			}
		}
	}
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Any() T {
	for t := range s {
		return t
	}

	var t T
	return t
}

func parseBrick(str string) Brick {
	start, end, _ := strings.Cut(str, "~")
	return Brick{parseCube(start), parseCube(end)}
}

func task1(args []string) error {
	bricks, err := utils.ReadLinesTransform(args[0], parseBrick)
	if err != nil {
		return err
	}

	hasMoved := true
	for hasMoved {
		hasMoved = false

		for i, brick := range bricks {
			if checkCanFall(bricks, brick) {
				bricks[i] = brick.MoveDown()
				hasMoved = true
			}
		}
	}

	essentials := make(Set[int])
	for _, brick := range bricks {
		touches := checkTouches(bricks, brick)
		if len(touches) == 1 {
			essentials[touches.Any()] = struct{}{}
		}
	}

	fmt.Printf("Disintegrate: %d\n", len(bricks)-len(essentials))
	return nil
}

func checkCanFall(bricks []Brick, brick Brick) bool {
	for cube := range brick.Iterate() {
		down := Cube{cube.X, cube.Y, cube.Z - 1}
		if down.Z < 1 {
			return false
		}

		for _, b := range bricks {
			if b == brick {
				continue
			}

			if b.Contains(down) {
				return false
			}
		}
	}

	return true
}

func checkTouches(bricks []Brick, brick Brick) Set[int] {
	touches := make(Set[int])
	for cube := range brick.Iterate() {
		down := Cube{cube.X, cube.Y, cube.Z - 1}

		for i, b := range bricks {
			if brick == b {
				continue
			}

			if b.Contains(down) {
				touches[i] = struct{}{}
			}
		}
	}

	return touches
}

func task2(args []string) error {
	bricks, err := utils.ReadLinesTransform(args[0], parseBrick)
	if err != nil {
		return err
	}

	hasMoved := true
	for hasMoved {
		hasMoved = false

		for i, brick := range bricks {
			if checkCanFall(bricks, brick) {
				bricks[i] = brick.MoveDown()
				hasMoved = true
			}
		}
	}

	nodes := make([]Node, len(bricks))
	for i := range nodes {
		nodes[i].ID = i
	}

	for i, brick := range bricks {
		touches := checkTouches(bricks, brick)
		for touch := range touches {
			nodes[i].parents = append(nodes[i].parents, &nodes[touch])
			nodes[touch].children = append(nodes[touch].children, &nodes[i])
		}
	}

	fmt.Println(nodes)

	sum := 0
	for i := range nodes {
		sum += nodes[i].Destroys() - 1
		for i := range nodes {
			nodes[i].Destroyed = false
		}
	}
	fmt.Printf("Destroyed: %d\n", sum)

	return nil

}

type Node struct {
	ID        int
	Destroyed bool
	parents   []*Node
	children  []*Node
}

func (n *Node) orphaned() bool {
	for _, parent := range n.parents {
		if !parent.Destroyed {
			return false
		}
	}
	return true
}

func (n *Node) Destroys() int {
	count := 1
	n.Destroyed = true

	for _, child := range n.children {
		if child.orphaned() {
			count += child.Destroys()
		}
	}

	return count
}
