package main

import (
	"bufio"
	"fmt"
	"iter"
	"os"

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

func IterateFile(name string) iter.Seq[byte] {
	return func(yield func(byte) bool) {
		f, err := os.Open(name)
		if err != nil {
			fmt.Printf("Error opening file: %s\n", err)
			return
		}
		defer f.Close()

		br := bufio.NewReader(f)

		for {
			b, err := br.ReadByte()
			if err != nil {
				return
			}

			if !yield(b) {
				return
			}
		}
	}
}

var dirs = map[byte]utils.Vec{
	'<': {X: -1}, // left / west
	'>': {X: +1}, // right / east
	'^': {Y: -1}, // up / north
	'v': {Y: +1}, // down / south
}

func task1(args []string) error {
	pos := utils.Vec{}
	houses := map[utils.Vec]bool{
		pos: true,
	}

	for b := range IterateFile(args[0]) {
		dir := dirs[b]
		pos = pos.Add(dir)
		houses[pos] = true
	}

	fmt.Printf("Visited houses: %d\n", len(houses))
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func IterateDouble[T any](it iter.Seq[T]) iter.Seq2[T, T] {
	return func(yield func(T, T) bool) {
		next, stop := iter.Pull(it)
		defer stop()

		for {
			a, ok := next()
			if !ok {
				return
			}

			b, _ := next()

			if !yield(a, b) {
				return
			}
		}
	}

}

func task2(args []string) error {
	santa := utils.Vec{}
	robo := utils.Vec{}
	houses := map[utils.Vec]bool{
		santa: true,
	}

	for a, b := range IterateDouble(IterateFile(args[0])) {
		dir := dirs[a]
		santa = santa.Add(dir)
		houses[santa] = true

		dir = dirs[b]
		robo = robo.Add(dir)
		houses[robo] = true
	}

	fmt.Printf("Visited houses: %d\n", len(houses))
	return nil

}
