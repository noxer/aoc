package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
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

type Keys struct {
	From, To byte
}

var numpadPositions = map[byte]utils.Vec{
	'7': {X: 0, Y: 0},
	'8': {X: 1, Y: 0},
	'9': {X: 2, Y: 0},
	'4': {X: 0, Y: 1},
	'5': {X: 1, Y: 1},
	'6': {X: 2, Y: 1},
	'1': {X: 0, Y: 2},
	'2': {X: 1, Y: 2},
	'3': {X: 2, Y: 2},
	'0': {X: 1, Y: 3},
	'A': {X: 2, Y: 3},
}

func pressNumpadKeys(seq string) string {
	sb := strings.Builder{}
	last := byte('A')

	for _, current := range []byte(seq) {
		sb.WriteString(pressNumpadKey(last, current))
		last = current
	}

	return sb.String()
}

func pressNumpadKey(last, current byte) string {
	from := numpadPositions[last]
	to := numpadPositions[current]

	move := to.Sub(from)

	var seq []byte
	if from.X == 0 && to.Y == 3 {
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, 'A')

		return pressDirpadKeys(seq, 2)
	}

	if from.Y == 3 && to.X == 0 {
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, 'A')

		return pressDirpadKeys(seq, 2)
	}

	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, 'A')
	sequenceA := pressDirpadKeys(seq, 2)

	seq = seq[:0]

	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, 'A')
	sequenceB := pressDirpadKeys(seq, 2)

	if len(sequenceA) < len(sequenceB) {
		return sequenceA
	}

	return sequenceB
}

func verticalKeys(y int) []byte {
	switch {
	case y == 0:
		return nil
	case y < 0:
		return bytes.Repeat([]byte{'^'}, -y)
	default:
		return bytes.Repeat([]byte{'v'}, y)
	}
}

func horizontalKeys(x int) []byte {
	switch {
	case x == 0:
		return nil
	case x < 0:
		return bytes.Repeat([]byte{'<'}, -x)
	default:
		return bytes.Repeat([]byte{'>'}, x)
	}
}

var dirpadPositions = map[byte]utils.Vec{
	'^': {X: 1, Y: 0},
	'A': {X: 2, Y: 0},
	'<': {X: 0, Y: 1},
	'v': {X: 1, Y: 1},
	'>': {X: 2, Y: 1},
}

func pressDirpadKeys(seq []byte, level int) string {
	if level == 0 {
		return string(seq)
	}

	sb := strings.Builder{}
	last := byte('A')

	for _, current := range seq {
		sb.WriteString(pressDirpadKey(last, current, level))
		last = current
	}

	return sb.String()
}

func pressDirpadKey(last, current byte, level int) string {
	from := dirpadPositions[last]
	to := dirpadPositions[current]

	move := to.Sub(from)

	var seq []byte

	if from.Y == 0 && to.X == 0 {
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, 'A')

		return pressDirpadKeys(seq, level-1)
	}

	if from.X == 0 && to.Y == 0 {
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, 'A')

		return pressDirpadKeys(seq, level-1)
	}

	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, 'A')
	seqA := pressDirpadKeys(seq, level-1)

	seq = seq[:0]
	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, 'A')
	seqB := pressDirpadKeys(seq, level-1)

	if len(seqA) < len(seqB) {
		return seqA
	}

	return seqB
}

func task1(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	sum := 0

	for _, line := range lines {
		seq := pressNumpadKeys(line)
		// fmt.Printf("%s: %s\n", line, seq)

		num, _ := strconv.Atoi(strings.TrimSuffix(line, "A"))
		sum += num * len(seq)
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func pressNumpadKeys2(seq string) int {
	last := byte('A')
	cache := make(map[CacheKey]int)

	size := 0

	for _, current := range []byte(seq) {
		size += pressNumpadKey2(last, current, cache)
		last = current
	}

	return size
}

func pressNumpadKey2(last, current byte, cache map[CacheKey]int) int {
	from := numpadPositions[last]
	to := numpadPositions[current]

	move := to.Sub(from)

	var seq []byte
	if from.X == 0 && to.Y == 3 {
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, 'A')

		return pressDirpadKeys2(seq, 25, cache)
	}

	if from.Y == 3 && to.X == 0 {
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, 'A')

		return pressDirpadKeys2(seq, 25, cache)
	}

	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, 'A')
	sequenceA := pressDirpadKeys2(seq, 25, cache)

	seq = seq[:0]

	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, 'A')
	sequenceB := pressDirpadKeys2(seq, 25, cache)

	return min(sequenceA, sequenceB)
}

type CacheKey struct {
	Seq   string
	Level int
}

func pressDirpadKeys2(seq []byte, level int, cache map[CacheKey]int) int {
	if level == 0 {
		return len(seq)
	}

	// ask cache
	key := CacheKey{string(seq), level}
	if res, ok := cache[key]; ok {
		// fmt.Printf("Hit the cache for %v: %d\n", key, res)
		return res
	}

	size := 0
	last := byte('A')

	for _, current := range seq {
		size += pressDirpadKey2(last, current, level, cache)
		last = current
	}

	cache[key] = size
	return size
}

func pressDirpadKey2(last, current byte, level int, cache map[CacheKey]int) int {
	from := dirpadPositions[last]
	to := dirpadPositions[current]

	move := to.Sub(from)

	var seq []byte

	if from.Y == 0 && to.X == 0 {
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, 'A')

		return pressDirpadKeys2(seq, level-1, cache)
	}

	if from.X == 0 && to.Y == 0 {
		seq = append(seq, horizontalKeys(move.X)...)
		seq = append(seq, verticalKeys(move.Y)...)
		seq = append(seq, 'A')

		return pressDirpadKeys2(seq, level-1, cache)
	}

	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, 'A')
	seqA := pressDirpadKeys2(seq, level-1, cache)

	seq = seq[:0]
	seq = append(seq, horizontalKeys(move.X)...)
	seq = append(seq, verticalKeys(move.Y)...)
	seq = append(seq, 'A')
	seqB := pressDirpadKeys2(seq, level-1, cache)

	return min(seqA, seqB)
}

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	sum := 0

	for _, line := range lines {
		seq := pressNumpadKeys2(line)
		// fmt.Printf("%s: %s\n", line, seq)

		num, _ := strconv.Atoi(strings.TrimSuffix(line, "A"))
		sum += num * seq
	}

	elapsed := time.Since(start)

	fmt.Printf("Sum: %d (%s)\n", sum, elapsed)

	return nil
}
