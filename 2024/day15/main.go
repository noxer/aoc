package main

import (
	"bufio"
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

var (
	Left  = utils.Vec{X: -1}
	Right = utils.Vec{X: +1}
)

var dirs = map[byte]utils.Vec{
	'^': {Y: -1},
	'v': {Y: +1},
	'<': Left,
	'>': Right,
}

func vecToGPS(pos utils.Vec) int {
	return pos.Y*100 + pos.X
}

type Map struct {
	size utils.Vec
	data map[utils.Vec]byte
}

func (m Map) SumCoords() int {
	sum := 0
	for pos, val := range m.data {
		if val != 'O' {
			continue
		}

		sum += vecToGPS(pos)
	}
	return sum
}

func (m Map) Move(pos utils.Vec, command byte) utils.Vec {
	dir := dirs[command]

	freeSpace, ok := m.CanPush(pos, dir)
	if !ok {
		return pos
	}

	newPos := pos.Add(dir)
	if m.data[newPos] == 'O' {
		delete(m.data, newPos)
		m.data[freeSpace] = 'O'
	}

	return newPos
}

func (m Map) CanPush(pos, dir utils.Vec) (utils.Vec, bool) {
	for {
		pos = pos.Add(dir)
		val := m.data[pos]

		switch val {
		case 'O':
			continue
		case '#':
			return utils.Vec{}, false
		default:
			return pos, true
		}
	}
}

func (m Map) Print() {
	for y := 0; y < m.size.Y; y++ {
		for x := 0; x < m.size.X; x++ {
			pos := utils.Vec{X: x, Y: y}
			if b, ok := m.data[pos]; ok {
				fmt.Print(string(b))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}
}

func ReadMapAndCommands(name string) (Map, utils.Vec, []byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return Map{}, utils.Vec{}, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	y := 0
	m := make(map[utils.Vec]byte)
	size := utils.Vec{}
	start := utils.Vec{}
	for s.Scan() {
		if s.Text() == "" {
			break
		}

		for x, b := range s.Bytes() {
			if b == '.' {
				continue
			}

			if b == '@' {
				start = utils.Vec{X: x, Y: y}
				continue
			}

			m[utils.Vec{X: x, Y: y}] = b
		}

		y++
		size.X = max(size.X, len(s.Bytes()))
	}

	var cmds []byte
	for s.Scan() {
		cmds = append(cmds, s.Bytes()...)
	}

	size.Y = y
	return Map{size: size, data: m}, start, cmds, s.Err()
}

func task1(args []string) error {
	warehouse, start, commands, err := ReadMapAndCommands(args[0])
	if err != nil {
		return err
	}

	// warehouse.Print()

	pos := start
	for _, command := range commands {
		pos = warehouse.Move(pos, command)
		// warehouse.Print()
	}

	fmt.Printf("Sum: %d\n", warehouse.SumCoords())

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (m Map) Fat() Map {
	right := dirs['>']

	m2 := Map{
		size: utils.Vec{
			X: m.size.X * 2,
			Y: m.size.Y,
		},
		data: make(map[utils.Vec]byte, len(m.data)*2),
	}

	for pos, val := range m.data {
		newPosA := utils.Vec{X: pos.X * 2, Y: pos.Y}
		newPosB := newPosA.Add(right)

		switch val {
		case 'O':
			m2.data[newPosA] = '['
			m2.data[newPosB] = ']'
		case '#':
			m2.data[newPosA] = '#'
			m2.data[newPosB] = '#'
		}
	}

	return m2
}

func (m Map) MoveFat(pos utils.Vec, command byte) utils.Vec {
	dir := dirs[command]

	newPos := pos.Add(dir)
	if !m.CanPushFat(newPos, dir) {
		return pos
	}

	m.PushFat(newPos, dir)
	return newPos
}

func (m Map) PushFat(pos, dir utils.Vec) {
	val, ok := m.data[pos]
	if !ok {
		return
	}

	newPos := pos.Add(dir)

	// left and right movement
	if dir == Left || dir == Right {
		m.PushFat(newPos, dir)
		delete(m.data, pos)
		m.data[newPos] = val

		return
	}

	// up and down movement
	switch val {
	case '[':
		m.PushFat(newPos, dir)
		delete(m.data, pos)
		m.data[newPos] = '['

		pos = pos.Add(Right)
		newPos = newPos.Add(Right)

		m.PushFat(newPos, dir)
		delete(m.data, pos)
		m.data[newPos] = ']'

	case ']':
		m.PushFat(newPos, dir)
		delete(m.data, pos)
		m.data[newPos] = ']'

		pos = pos.Add(Left)
		newPos = newPos.Add(Left)

		m.PushFat(newPos, dir)
		delete(m.data, pos)
		m.data[newPos] = '['
	}
}

func (m Map) CanPushFat(pos, dir utils.Vec) bool {
	val, ok := m.data[pos]
	if !ok {
		return true
	}
	if val == '#' {
		return false
	}

	// left and right movement
	if dir == Left || dir == Right {
		newPos := pos.Add(dir)
		return m.CanPushFat(newPos, dir)
	}

	// up and down movement
	switch val {
	case '[':
		newPos := pos.Add(dir)
		if !m.CanPushFat(newPos, dir) {
			return false
		}
		newPos = newPos.Add(Right)
		return m.CanPushFat(newPos, dir)

	case ']':
		newPos := pos.Add(dir)
		if !m.CanPushFat(newPos, dir) {
			return false
		}
		newPos = newPos.Add(Left)
		return m.CanPushFat(newPos, dir)
	}

	return false
}

func (m Map) SumCoordsFat() int {
	sum := 0
	for pos, val := range m.data {
		if val != '[' {
			continue
		}

		sum += vecToGPS(pos)
	}
	return sum
}

func task2(args []string) error {
	warehouse, pos, commands, err := ReadMapAndCommands(args[0])
	if err != nil {
		return err
	}

	start := time.Now()

	warehouse = warehouse.Fat()
	pos = utils.Vec{X: pos.X * 2, Y: pos.Y}

	for _, command := range commands {
		pos = warehouse.MoveFat(pos, command)
	}

	elapsed := time.Since(start)

	// warehouse.Print()
	fmt.Printf("Sum: %d (%s)\n", warehouse.SumCoordsFat(), elapsed)

	return nil
}
