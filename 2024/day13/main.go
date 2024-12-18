package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
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

type Machine struct {
	ButtonA utils.Vec
	ButtonB utils.Vec
	Prize   utils.Vec
}

func (m Machine) Cost() int {
	cost := 999999999999999999

	for a := 0; lowerEq(m.ButtonA.Mul(a), m.Prize); a++ {
		posA := m.ButtonA.Mul(a)
		for b := 0; lowerEq(m.ButtonB.Mul(b).Add(posA), m.Prize); b++ {
			posB := m.ButtonB.Mul(b).Add(posA)
			if posB == m.Prize {
				cost = min(cost, a*3+b)
			}
		}
	}

	if cost == 999999999999999999 {
		cost = 0
	}

	return cost
}

func lowerEq(a, b utils.Vec) bool {
	return a.X <= b.X && a.Y <= b.Y
}

var matchOffsets = regexp.MustCompile(`^([^:]+):\sX[+=](\d+),\sY[+=](\d+)$`)

func parseMachines(name string) ([]Machine, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	var machines []Machine
	m := Machine{}
	for s.Scan() {
		line := s.Text()

		matches := matchOffsets.FindStringSubmatch(line)
		if len(matches) != 4 {
			machines = append(machines, m)
			m = Machine{}
			continue
		}

		v := utils.Vec{}
		v.X, _ = strconv.Atoi(matches[2])
		v.Y, _ = strconv.Atoi(matches[3])

		switch matches[1] {
		case "Button A":
			m.ButtonA = v
		case "Button B":
			m.ButtonB = v
		case "Prize":
			m.Prize = v
		}
	}

	zero := Machine{}
	if m != zero {
		machines = append(machines, m)
	}

	return machines, nil
}

func task1(args []string) error {
	machines, err := parseMachines(args[0])
	if err != nil {
		return err
	}

	sum := 0
	for _, machine := range machines {
		cost := machine.Cost()
		sum += cost
	}

	fmt.Printf("Cost: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type Line struct {
	Slope  float64
	Offset float64
}

type Vec2 struct {
	X, Y float64
}

func (v Vec2) Sub(o Vec2) Vec2 {
	return Vec2{
		X: v.X - o.X,
		Y: v.Y - o.Y,
	}
}

func (l Line) MoveThrough(pos utils.Vec) Line {
	y := l.Slope * float64(pos.X)
	l.Offset = float64(pos.Y) - y
	return l
}

func (l Line) Eval(x float64) float64 {
	return l.Slope*x + l.Offset
}

func (l Line) Intersection(o Line) (Vec2, bool) {
	slope := l.Slope - o.Slope
	if slope == 0 {
		return Vec2{}, false
	}

	x := (o.Offset - l.Offset) / slope

	return Vec2{X: x, Y: l.Eval(x)}, true
}

func lineFromVec(v utils.Vec) Line {
	return Line{
		Slope: float64(v.Y) / float64(v.X),
	}
}

func (m Machine) FastCost() int {
	lineA := lineFromVec(m.ButtonA)
	lineB := lineFromVec(m.ButtonB)

	lineB = lineB.MoveThrough(m.Prize)

	intersection, ok := lineA.Intersection(lineB)
	if !ok {
		return 0
	}

	a := checkMultiple(m.ButtonA, intersection)

	prizeVec2 := Vec2{X: float64(m.Prize.X), Y: float64(m.Prize.Y)}
	prizeOffset := prizeVec2.Sub(intersection)

	b := checkMultiple(m.ButtonB, prizeOffset)

	maybePrize := m.ButtonA.Mul(int(a)).Add(m.ButtonB.Mul(int(b)))
	if maybePrize != m.Prize {
		// fmt.Printf("\nFound discrepancy: %v vs %v\n", maybePrize, m.Prize)
		return 0
	}

	return a*3 + b
}

func checkMultiple(v utils.Vec, intersection Vec2) int {
	divX := intersection.X / float64(v.X)
	divY := intersection.Y / float64(v.Y)

	// fmt.Printf("checkMultiple %v %v: %f and %f\n", v, intersection, divX, divY)

	if int(math.Round(divX)) != int(math.Round(divY)) {
		return 0
	}

	return int(math.Round(divX))
}

func task2(args []string) error {
	machines, err := parseMachines(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	sum := 0
	for _, machine := range machines {
		// fmt.Printf("Solving machine %d...", i)

		machine.Prize.X += 10000000000000
		machine.Prize.Y += 10000000000000

		cost := machine.FastCost()
		// fmt.Printf(" (%d) ", cost)
		sum += cost

		// fmt.Println("ok")
	}

	elapsed := time.Since(start)

	fmt.Printf("Cost: %d (%s)\n", sum, elapsed)

	return nil
}
