package main

import (
	"fmt"
	"os"

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

type Gate func(map[string]uint16) bool

func parseGate(line string) Gate {
	c := uint16(0)
	out := ""
	x, y := "", ""

	if _, err := fmt.Sscanf(line, "%d -> %s", &c, &out); err == nil {
		fmt.Printf("%d -> %q\n", c, out)
		return func(m map[string]uint16) bool {
			if _, ok := m[out]; ok {
				return true
			}

			fmt.Printf("%d -> %s\n", c, out)
			m[out] = c
			return true
		}
	}

	if _, err := fmt.Sscanf(line, "%s -> %s", &x, &out); err == nil {
		fmt.Printf("%q -> %q\n", x, out)
		return func(m map[string]uint16) bool {
			if a, ok := m[x]; ok {
				fmt.Printf("%s (%d) -> %s\n", x, a, out)
				m[out] = a
				return true
			}

			// fmt.Printf("%s -> %s skipped\n", x, out)
			return false
		}
	}

	if _, err := fmt.Sscanf(line, "NOT %s -> %s", &x, &out); err == nil {
		fmt.Printf("NOT %q -> %q\n", x, out)
		return func(m map[string]uint16) bool {
			if n, ok := m[x]; ok {
				fmt.Printf("NOT %s (%d) -> %s\n", x, n, out)
				m[out] = ^n
				return true
			}

			// fmt.Printf("NOT %s -> %s skipped\n", x, out)
			return false
		}
	}

	if _, err := fmt.Sscanf(line, "%s LSHIFT %d -> %s", &x, &c, &out); err == nil {
		fmt.Printf("%q LSHIFT %d -> %q\n", x, c, out)
		return func(m map[string]uint16) bool {
			if n, ok := m[x]; ok {
				fmt.Printf("%s (%d) << %d -> %s\n", x, n, c, out)
				m[out] = n << c
				return true
			}

			// fmt.Printf("%s LSHIFT %d -> %s skipped\n", x, c, out)
			return false
		}
	}

	if _, err := fmt.Sscanf(line, "%s RSHIFT %d -> %s", &x, &c, &out); err == nil {
		fmt.Printf("%q RSHIFT %d -> %q\n", x, c, out)
		return func(m map[string]uint16) bool {
			if n, ok := m[x]; ok {
				fmt.Printf("%s (%d) >> %d -> %s\n", x, n, c, out)
				m[out] = n >> c
				return true
			}

			// fmt.Printf("%s RSHIFT %d -> %s skipped\n", x, c, out)
			return false
		}
	}

	if _, err := fmt.Sscanf(line, "%d AND %s -> %s", &c, &y, &out); err == nil {
		fmt.Printf("%d AND %q -> %q\n", c, y, out)
		return func(m map[string]uint16) bool {
			b, ok := m[y]
			if !ok {
				// fmt.Printf("%s AND %s -> %s skipped right\n", x, y, out)
				return false
			}

			fmt.Printf("%d AND %s (%d) -> %s\n", c, y, b, out)
			m[out] = c & b
			return true
		}
	}

	if _, err := fmt.Sscanf(line, "%s AND %s -> %s", &x, &y, &out); err == nil {
		fmt.Printf("%q AND %q -> %q\n", x, y, out)
		return func(m map[string]uint16) bool {
			a, ok := m[x]
			if !ok {
				// fmt.Printf("%s AND %s -> %s skipped left\n", x, y, out)
				return false
			}

			b, ok := m[y]
			if !ok {
				// fmt.Printf("%s AND %s -> %s skipped right\n", x, y, out)
				return false
			}

			fmt.Printf("%s (%d) AND %s (%d) -> %s\n", x, a, y, b, out)
			m[out] = a & b
			return true
		}
	}

	if _, err := fmt.Sscanf(line, "%s OR %s -> %s", &x, &y, &out); err == nil {
		fmt.Printf("%q OR %q -> %q\n", x, y, out)
		return func(m map[string]uint16) bool {
			a, ok := m[x]
			if !ok {
				// fmt.Printf("%s OR %s -> %s skipped left\n", x, y, out)
				return false
			}

			b, ok := m[y]
			if !ok {
				// fmt.Printf("%s OR %s -> %s skipped right\n", x, y, out)
				return false
			}

			fmt.Printf("%s (%d) OR %s (%d) -> %s\n", x, a, y, b, out)
			m[out] = a | b
			return true
		}
	}

	panic(line)
}

type CPU struct {
	Gates []Gate
	Wires map[string]uint16
}

func (c *CPU) Run() bool {
	next := c.Gates[:0]

	for _, gate := range c.Gates {
		if !gate(c.Wires) {
			next = append(next, gate)
		}
	}

	c.Gates = next
	return len(c.Gates) != 0
}

func task1(args []string) error {
	gates, err := utils.ReadLinesTransform(args[0], parseGate)
	if err != nil {
		return err
	}

	cpu := CPU{
		Gates: gates,
		Wires: make(map[string]uint16),
	}

	for cpu.Run() {
		// fmt.Println(len(cpu.Gates))
	}

	fmt.Printf("A: %d\n", cpu.Wires["a"])
	// fmt.Println(cpu.Wires)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	gates, err := utils.ReadLinesTransform(args[0], parseGate)
	if err != nil {
		return err
	}

	cpu := CPU{
		Gates: gates,
		Wires: make(map[string]uint16),
	}

	for cpu.Run() {
		// fmt.Println(len(cpu.Gates))
	}

	fmt.Printf("A: %d\n", cpu.Wires["a"])
	// fmt.Println(cpu.Wires)

	gates, err = utils.ReadLinesTransform(args[0], parseGate)
	if err != nil {
		return err
	}

	cpu2 := CPU{
		Gates: gates,
		Wires: map[string]uint16{
			"b": cpu.Wires["a"],
		},
	}

	for cpu2.Run() {

	}

	fmt.Printf("A: %d\n", cpu2.Wires["a"])

	return nil
}
