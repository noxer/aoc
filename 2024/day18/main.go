package main

import (
	"fmt"
	"math"
	"os"
	"sort"
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

var Directions = []utils.Vec{
	{X: -1}, // left / west
	{Y: +1}, // down / south
	{X: +1}, // right / east
	{Y: -1}, // up / north
}

type MemorySpace struct {
	data map[utils.Vec]int
	size utils.Vec
}

func (ms MemorySpace) Contains(pos utils.Vec) bool {
	return pos.X >= 0 && pos.X < ms.size.X && pos.Y >= 0 && pos.Y < ms.size.Y
}

func (ms MemorySpace) FromSlice(bytes []utils.Vec) {
	for i, v := range bytes {
		ms.data[v] = i
	}
}

func (ms MemorySpace) FindPath(from, to utils.Vec) int {
	seen := make(map[utils.Vec]int)
	return ms.findPath(from, to, seen, 1) - 1
}

func (ms MemorySpace) findPath(pos, to utils.Vec, seen map[utils.Vec]int, steps int) int {
	// found the exit
	if pos == to {
		return 1
	}

	// outside of memory space
	if !ms.Contains(pos) {
		return -1
	}

	// memory corruption
	if _, ok := ms.data[pos]; ok {
		return -1
	}

	// we were here before
	if score, ok := seen[pos]; ok && score <= steps {
		return -1
	}
	seen[pos] = steps

	best := math.MaxInt
	for _, dir := range utils.Directions {
		newPos := pos.Add(dir)
		res := ms.findPath(newPos, to, seen, steps+1)
		if res >= 0 {
			best = min(best, res)
		}
	}

	if best != math.MaxInt {
		return best + 1
	} else {
		return -1
	}
}

func parseByte(line string) utils.Vec {
	vec := utils.Vec{}
	fmt.Sscanf(line, "%d,%d", &vec.X, &vec.Y)
	return vec
}

func task1(args []string) error {
	bytes, err := utils.ReadLinesTransform(args[0], parseByte)
	if err != nil {
		return err
	}

	fmt.Println(bytes)

	ms := MemorySpace{
		data: make(map[utils.Vec]int, len(bytes)),
		size: utils.Vec{
			X: 71,
			Y: 71,
		},
	}
	ms.FromSlice(bytes[:1024])

	start := utils.Vec{X: 0, Y: 0}
	end := utils.Vec{X: 70, Y: 70}

	begin := time.Now()

	length := ms.FindPath(start, end)

	elapsed := time.Since(begin)

	fmt.Printf("Path length: %d (%s)\n", length, elapsed)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

var seen = make(map[utils.Vec]int)

func (ms MemorySpace) FindPath2(from, to utils.Vec) []utils.Vec {
	clear(seen)
	return ms.findPath2(from, to, seen, 1)
}

func (ms MemorySpace) findPath2(pos, to utils.Vec, seen map[utils.Vec]int, steps int) []utils.Vec {
	// found the exit
	if pos == to {
		res := make([]utils.Vec, 1, steps)
		res[0] = pos
		return res
	}

	// outside of memory space
	if !ms.Contains(pos) {
		return nil
	}

	// memory corruption
	if _, ok := ms.data[pos]; ok {
		return nil
	}

	// we were here before
	if score, ok := seen[pos]; ok && score <= steps {
		return nil
	}
	seen[pos] = steps

	for _, dir := range utils.Directions {
		newPos := pos.Add(dir)
		res := ms.findPath2(newPos, to, seen, steps+1)
		if len(res) > 0 {
			return append(res, pos)
		}
	}

	return nil
}

func (ms MemorySpace) CheckPaths(paths [][]utils.Vec) bool {
outer:
	for _, path := range paths {
		for _, location := range path {
			if _, ok := ms.data[location]; !ok {
				continue outer
			}
		}

		return true
	}

	return false
}

func task2(args []string) error {
	bytes, err := utils.ReadLinesTransform(args[0], parseByte)
	if err != nil {
		return err
	}

	// fmt.Println(bytes)

	ms := MemorySpace{
		data: make(map[utils.Vec]int, len(bytes)),
		size: utils.Vec{
			X: 71,
			Y: 71,
		},
	}

	start := utils.Vec{X: 0, Y: 0}
	end := utils.Vec{X: 70, Y: 70}

	begin := time.Now()

	paths := make([][]utils.Vec, 0, 100)

	lastI := 0
	i, _ := sort.Find(len(bytes), func(i int) int {
		// fmt.Printf("Trying byte %v at %d\n", bytes[i], i)

		if i < lastI {
			clear(ms.data)
			ms.FromSlice(bytes[:i+1])
		} else {
			ms.FromSlice(bytes[lastI : i+1])
		}
		lastI = i

		if ms.CheckPaths(paths) {
			return 1
		}

		result := ms.FindPath2(start, end)
		if len(result) == 0 {
			return -1
		}
		paths = append(paths, result)

		return 1
	})
	if i < 0 || i >= len(bytes) {
		// fmt.Println("Couldn't find blocking byte...")
		return nil
	}

	elapsed := time.Since(begin)

	fmt.Printf("Found result: %d,%d (%s)\n", bytes[i].X, bytes[i].Y, elapsed)

	return nil
}
