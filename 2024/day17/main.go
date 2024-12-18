package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
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

type CPU struct {
	A, B, C int
	PC      int
	Memory  []int
	Output  []int
}

func (cpu *CPU) loadCombo(combo int) int {
	if combo <= 3 {
		return combo
	}

	switch combo {
	case 4:
		return cpu.A
	case 5:
		return cpu.B
	case 6:
		return cpu.C
	default:
		panic("unexpected combo operand")
	}
}

func (cpu *CPU) Step() bool {
	if cpu.PC < 0 || cpu.PC >= len(cpu.Memory) {
		return false
	}

	inst := cpu.Memory[cpu.PC]
	oper := cpu.Memory[cpu.PC+1]

	switch inst {
	case 0: // adv
		oper = cpu.loadCombo(oper)
		denom := 1 << oper
		cpu.A /= denom

	case 1: // bxl
		cpu.B ^= oper

	case 2: // bst
		oper = cpu.loadCombo(oper)
		cpu.B = oper & 7

	case 3: // jnz
		if cpu.A == 0 {
			break
		}

		cpu.PC = oper - 2

	case 4: // bxc
		cpu.B ^= cpu.C

	case 5: // out
		oper = cpu.loadCombo(oper)
		cpu.Output = append(cpu.Output, oper&7)

	case 6: // bdv
		oper = cpu.loadCombo(oper)
		denom := 1 << oper
		cpu.B = cpu.A / denom

	case 7: // cdv
		oper = cpu.loadCombo(oper)
		denom := 1 << oper
		cpu.C = cpu.A / denom
	}

	cpu.PC += 2
	return true
}

func (cpu *CPU) Run() {
	for cpu.Step() {
	}
}

func (cpu *CPU) Reset(a int) {
	cpu.A = a
	cpu.B = 0
	cpu.C = 0
	cpu.PC = 0
	cpu.Output = cpu.Output[:0]
}

func loadCPU(name string) (*CPU, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)
	cpu := &CPU{}

	if !s.Scan() {
		return nil, errors.New("unexpected end")
	}
	cpu.A, _ = strconv.Atoi(strings.TrimPrefix(s.Text(), "Register A: "))

	if !s.Scan() {
		return nil, errors.New("unexpected end")
	}
	cpu.B, _ = strconv.Atoi(strings.TrimPrefix(s.Text(), "Register B: "))

	if !s.Scan() {
		return nil, errors.New("unexpected end")
	}
	cpu.C, _ = strconv.Atoi(strings.TrimPrefix(s.Text(), "Register C: "))

	if !s.Scan() {
		return nil, errors.New("unexpected end")
	}
	if !s.Scan() {
		return nil, errors.New("unexpected end")
	}

	cpu.Memory = utils.ParseInts(strings.TrimPrefix(s.Text(), "Program: "), ",")

	return cpu, nil
}

func task1(args []string) error {
	cpu, err := loadCPU(args[0])
	if err != nil {
		return err
	}

	for cpu.Step() {
	}

	fmt.Println(cpu.Output)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	cpu, err := loadCPU(args[0])
	if err != nil {
		return err
	}

	// pow := math.Pow(8, 15)*6 +
	// 	math.Pow(8, 14)*5 +
	// 	math.Pow(8, 12)*18

	// math.Pow(8, 12) +
	// math.Pow(8, 11) +
	// math.Pow(8, 10)*2 +
	// math.Pow(8, 9)*7

	// match := []int{6, 6, 6, 4, 5, 1, 3, 0}

	start := time.Now()
	result := trySolve(cpu, 0)
	elapsed := time.Since(start)

	fmt.Printf("Result: %d (%s)\n", result, elapsed)

	return nil
}

func trySolve(cpu *CPU, a int) int {
	for b := range 8 {
		cpu.Reset(a + b)
		cpu.Run()

		// fmt.Printf("Trying %v\n", cpu.Output)

		if hasSuffix(cpu.Memory, cpu.Output) {
			if slices.Equal(cpu.Memory, cpu.Output) {
				return a + b
			}

			res := trySolve(cpu, (a+b)*8)
			if res > 0 {
				return res
			}
		}
	}

	return -1
}

func hasSuffix(sl, suffix []int) bool {
	if len(suffix) > len(sl) {
		return false
	}

	diff := len(sl) - len(suffix)
	return slices.Equal(sl[diff:], suffix)
}

// func reverseProgram(output []int) int {
// 	slices.Reverse(output)

// 	var a, b, c int

// 	a = 48

// 	for _, n := range output {
// 		b = n ^ 5

// 	}
// }
