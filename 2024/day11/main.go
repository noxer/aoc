package main

import (
	"fmt"
	"os"
	"slices"
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

func applyRules(stones []int) []int {
	for i := 0; i < len(stones); i++ {
		stone := stones[i]

		if stone == 0 {
			stones[i] = 1
		} else if str := strconv.Itoa(stone); len(str)%2 == 0 {
			leftStone, _ := strconv.Atoi(str[:len(str)/2])
			rightStone, _ := strconv.Atoi(str[len(str)/2:])

			stones[i] = rightStone
			stones = slices.Insert(stones, i, leftStone)
			i++
		} else {
			stones[i] = stone * 2024
		}
	}

	return stones
}

func task1(args []string) error {
	p, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	stones := utils.ParseInts(string(p), " ")

	for range 25 {
		stones = applyRules(stones)
	}

	fmt.Printf("Stones count: %d\n", len(stones))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type Key struct {
	Stone int
	Steps int
}

func countStones(cache map[Key]int, stone, steps int) (count int) {
	if steps == 0 {
		return 1
	}

	if count, ok := cache[Key{stone, steps}]; ok {
		return count
	}

	defer func() {
		cache[Key{stone, steps}] = count
	}()

	if stone == 0 {
		return countStones(cache, 1, steps-1)
	} else if str := strconv.Itoa(stone); len(str)%2 == 0 {
		leftStone, _ := strconv.Atoi(str[:len(str)/2])
		rightStone, _ := strconv.Atoi(str[len(str)/2:])

		return countStones(cache, leftStone, steps-1) + countStones(cache, rightStone, steps-1)
	} else {
		return countStones(cache, stone*2024, steps-1)
	}
}

func task2(args []string) error {
	p, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}

	stones := utils.ParseInts(string(p), " ")

	start := time.Now()

	sum := 0
	cache := make(map[Key]int)
	for _, stone := range stones {
		sum += countStones(cache, stone, 75)
	}

	elapsed := time.Since(start)

	fmt.Printf("Stones count: %d (%s)\n", sum, elapsed)

	return nil
}
