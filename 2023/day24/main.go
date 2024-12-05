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

type Vec3 struct {
	X, Y, Z float64
}

func parseVec3(str string) Vec3 {
	parts := strings.Split(str, ", ")
	v := Vec3{}
	v.X, _ = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	v.Y, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	v.Z, _ = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	return v
}

type Hailstone struct {
	Position Vec3
	Velocity Vec3
}

func (h Hailstone) CheckX(x float64) bool {
	if h.Velocity.X > 0 {
		return x >= h.Position.X
	}

	return x <= h.Position.X
}

func (h Hailstone) ToLine2() Trajectory {
	t := Trajectory{}
	t.Incline = h.Velocity.Y / h.Velocity.X
	t.Offset = t.Incline*-h.Position.X + h.Position.Y
	return t
}

type Trajectory struct {
	Offset, Incline float64
}

func (t Trajectory) Eval(x float64) float64 {
	return t.Incline*x + t.Offset
}

func (t Trajectory) Intersects(o Trajectory) (Vec3, bool) {
	// t.Incline * x + t.Offset = o.Incline * x + o.Offset | - o.Offset
	// t.Incline * x = o.Incline * x + o.Offset - t.Offset | - o.Incline * x
	// (t.Incline-o.Incline) * x = o.Offset - t.Offset     |
	inc := t.Incline - o.Incline
	if inc == 0 {
		return Vec3{}, false
	}

	x := (o.Offset - t.Offset) / inc

	return Vec3{X: x, Y: t.Eval(x)}, true
}

func parseHailstone(line string) Hailstone {
	pos, vel, _ := strings.Cut(line, " @ ")
	h := Hailstone{
		Position: parseVec3(pos),
		Velocity: parseVec3(vel),
	}
	return h
}

func hailstonesCollide(a, b Hailstone) (Vec3, bool) {
	pos, ok := a.ToLine2().Intersects(b.ToLine2())
	if !ok {
		return Vec3{}, false
	}

	valid := a.CheckX(pos.X) && b.CheckX(pos.X)
	return pos, valid
}

type BoundingBox struct {
	A, B Vec3
}

func (bb BoundingBox) HailstonesCollide(a, b Hailstone) bool {
	pos, ok := hailstonesCollide(a, b)
	if !ok {
		return false
	}

	return pos.X >= bb.A.X && pos.X <= bb.B.X &&
		pos.Y >= bb.A.Y && pos.Y <= bb.B.Y
}

func task1(args []string) error {
	hailstones, err := utils.ReadLinesTransform(args[0], parseHailstone)
	if err != nil {
		return err
	}

	// box := BoundingBox{
	// 	A: Vec3{X: 7, Y: 7},
	// 	B: Vec3{X: 27, Y: 27},
	// }

	box := BoundingBox{
		A: Vec3{X: 200000000000000, Y: 200000000000000},
		B: Vec3{X: 400000000000000, Y: 400000000000000},
	}

	counter := 0
	for i, a := range hailstones {
		for _, b := range hailstones[i+1:] {
			if box.HailstonesCollide(a, b) {
				counter++
			}
		}
	}

	fmt.Printf("Counter: %d\n", counter)

	return nil
}

func task2(args []string) error {
	return nil
}

// one vector to intersect them all

// |
// |     \
// -------\------x
// |\      x
// | \      \
// |  \      \
// |   \      \
// +----\------------------
// |     \
// |
// |
// |
// |
// |
