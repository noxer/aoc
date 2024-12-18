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

type Robot struct {
	Pos utils.Vec
	Vel utils.Vec
}

func parseRobot(line string) Robot {
	rawP, rawV, _ := strings.Cut(line, " ")

	return Robot{
		Pos: parseVec(rawP),
		Vel: parseVec(rawV),
	}
}

func parseVec(part string) utils.Vec {
	var result utils.Vec

	part = part[2:]
	rawX, rawY, _ := strings.Cut(part, ",")

	result.X, _ = strconv.Atoi(rawX)
	result.Y, _ = strconv.Atoi(rawY)

	return result
}

func mapToRoom(width, height int, pos utils.Vec) utils.Vec {
	x := pos.X % width
	y := pos.Y % height

	if x < 0 {
		x += width
	}
	if y < 0 {
		y += height
	}

	return utils.Vec{
		X: x,
		Y: y,
	}
}

func countQuadrant(pos []utils.Vec, width, height int) (a, b, c, d int) {
	centerX := width / 2
	centerY := height / 2

	for _, p := range pos {
		switch {
		case p.X < centerX && p.Y < centerY:
			a++
		case p.X > centerX && p.Y < centerY:
			b++
		case p.X < centerX && p.Y > centerY:
			c++
		case p.X > centerX && p.Y > centerY:
			d++
		}
	}

	return a, b, c, d
}

func task1(args []string) error {
	robots, err := utils.ReadLinesTransform(args[0], parseRobot)
	if err != nil {
		return err
	}

	width := 101
	height := 103
	// width := 11
	// height := 7
	seconds := 100

	pos := make([]utils.Vec, len(robots))
	for i, robot := range robots {
		pos[i] = robot.Pos.Add(robot.Vel.Mul(seconds))
	}

	for i, p := range pos {
		pos[i] = mapToRoom(width, height, p)
	}

	a, b, c, d := countQuadrant(pos, width, height)

	fmt.Printf("%d * %d * %d * %d = %d\n", a, b, c, d, a*b*c*d)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

var directions = []utils.Vec{
	{X: -1},
	{X: +1},
	{Y: -1},
	{Y: +1},
}

func scoreMap(pos []utils.Vec) int {
	score := 0

	for _, p := range pos {
		for _, d := range directions {
			if slices.Contains(pos, p.Add(d)) {
				score++
			}
		}
	}

	return score
}

func printMap(pos []utils.Vec, width, height int) {
	for y := range height {
		for x := range width {
			if slices.Contains(pos, utils.Vec{X: x, Y: y}) {
				fmt.Print("X")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func task2(args []string) error {
	robots, err := utils.ReadLinesTransform(args[0], parseRobot)
	if err != nil {
		return err
	}

	width := 101
	height := 103
	// width := 11
	// height := 7

	pos := make([]utils.Vec, len(robots))

	maxScore := 0
	maxScoreSeconds := 0
	for seconds := range 10000 {
		for i, robot := range robots {
			pos[i] = robot.Pos.Add(robot.Vel.Mul(seconds))
		}

		for i, p := range pos {
			pos[i] = mapToRoom(width, height, p)
		}

		score := scoreMap(pos)
		if score > maxScore {
			maxScore = score
			maxScoreSeconds = seconds

			fmt.Printf("Seconds: %d; Score: %d\n", seconds, score)
			printMap(pos, width, height)
		}
	}

	fmt.Printf("Score: %d at %d\n", maxScore, maxScoreSeconds)

	return nil
}
