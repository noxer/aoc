package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	var err error

	// strs := []string{"qff", "qnw", "z16", "pbv", "qqp", "z23", "fbq", "z36"}
	// sort.Strings(strs)
	// fmt.Printf("Result: %s\n", strings.Join(strs, ","))

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

type WireState byte

const (
	Undefined WireState = iota
	False
	True
)

func (ws WireState) String() string {
	switch ws {
	case False:
		return "false"
	case True:
		return "true"
	}

	return "undefined"
}

func parseWire(wire string) WireState {
	switch wire {
	case "0":
		return False
	case "1":
		return True
	}

	return Undefined
}

func (ws WireState) And(o WireState) WireState {
	if ws == True && o == True {
		return True
	}

	return False
}

func (ws WireState) Or(o WireState) WireState {
	if ws == True || o == True {
		return True
	}

	return False
}

func (ws WireState) Xor(o WireState) WireState {
	if ws != o {
		return True
	}

	return False
}

type Gate struct {
	Op         func(a, b WireState) WireState
	A, B, C, O string
}

func (g Gate) String() string {
	return fmt.Sprintf("%s %s %s -> %s", g.A, g.O, g.B, g.C)
}

func NewAnd(a, b, c string) Gate {
	return Gate{
		Op: WireState.And,
		A:  a,
		B:  b,
		C:  c,
		O:  "AND",
	}
}

func NewOr(a, b, c string) Gate {
	return Gate{
		Op: WireState.Or,
		A:  a,
		B:  b,
		C:  c,
		O:  "OR",
	}
}

func NewXor(a, b, c string) Gate {
	return Gate{
		Op: WireState.Xor,
		A:  a,
		B:  b,
		C:  c,
		O:  "XOR",
	}
}

func (g Gate) TryRun(wires Wires) bool {
	a := wires[g.A]
	b := wires[g.B]

	if a == Undefined || b == Undefined {
		return false
	}

	wires[g.C] = g.Op(a, b)

	return true
}

func parseGate(line string) Gate {
	var a, b, c, op string
	fmt.Sscanf(line, "%s %s %s -> %s", &a, &op, &b, &c)

	switch op {
	case "AND":
		return NewAnd(a, b, c)
	case "OR":
		return NewOr(a, b, c)
	case "XOR":
		return NewXor(a, b, c)
	}

	panic("Couldn't parse line: " + line)
}

type Wires map[string]WireState

func (w Wires) SetX(num uint64) {
	for i := range 45 {
		if (num>>i)&1 == 1 {
			w[fmt.Sprintf("x%02d", i)] = True
		} else {
			w[fmt.Sprintf("x%02d", i)] = False
		}
	}
}

func (w Wires) SetY(num uint64) {
	for i := range 45 {
		if (num>>i)&1 == 1 {
			w[fmt.Sprintf("y%02d", i)] = True
		} else {
			w[fmt.Sprintf("y%02d", i)] = False
		}
	}
}

func (w Wires) CalculateOutput() uint64 {
	output := uint64(0)
	for z := 0; z < 64; z++ {
		state := w[fmt.Sprintf("z%02d", z)]
		if state == True {
			output |= 1 << z
		}
	}
	return output
}

func parseWiresAndGates(name string) (Wires, []Gate, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	wires := make(Wires)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			break
		}

		wire, value, _ := strings.Cut(line, ": ")
		wires[wire] = parseWire(value)
	}

	var gates []Gate
	for s.Scan() {
		line := s.Text()
		gates = append(gates, parseGate(line))
	}

	return wires, gates, s.Err()
}

func runRound(wires Wires, gates []Gate) []Gate {
	var next []Gate

	for _, gate := range gates {
		if !gate.TryRun(wires) {
			next = append(next, gate)
		}
	}

	return next
}

func task1(args []string) error {
	wires, gates, err := parseWiresAndGates(args[0])
	if err != nil {
		return err
	}

	for len(gates) > 0 {
		gates = runRound(wires, gates)
	}

	output := wires.CalculateOutput()

	fmt.Printf("Output: %d\n", output)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func checkBit(gates []Gate, bit int) bool {
	wires := make(Wires)
	original := gates
	var last []Gate

	for a := range uint64(4) {
		for b := range uint64(4) {
			x := a << max(bit-1, 0)
			y := b << max(bit-1, 0)

			wires.SetX(x)
			wires.SetY(y)

			for len(gates) > 0 {
				// fmt.Printf("len(gates) == %d, len(last) == %d\n", len(gates), len(last))
				gates = runRound(wires, gates)

				if len(last) == len(gates) {
					return false
				}

				last = gates
			}

			output := wires.CalculateOutput()
			if output != x+y {
				// fmt.Printf("%d != %d + %d\n", output, x, y)
				return false
			}

			clear(wires)
			gates = original
		}
	}

	return true
}

func task2(args []string) error {
	_, gates, err := parseWiresAndGates(args[0])
	if err != nil {
		return err
	}

	var swappedPairs []string

	for n := range 45 {
		if !checkBit(gates, n) {
			fmt.Printf("Found error in bit %d\n", n)

			adder := extractAdder(gates, n)
			fmt.Printf("Extracted adder: %v\n", adder)

			swapped := mutateAndTest(gates, adder, n)
			if len(swapped) == 0 {
				fmt.Printf("Couldn't find a permutation for bit %d that works :(\n", n)
			}

			fmt.Printf("Found pair: %v\n", swapped)
			swappedPairs = append(swappedPairs, swapped...)
		}
	}

	sort.Strings(swappedPairs)
	fmt.Printf("Result: %s\n", strings.Join(swappedPairs, ","))

	return nil
}

func mutateAndTest(gates []Gate, adder []int, n int) []string {
	for i, a := range adder {
		for _, b := range adder[i+1:] {
			// fmt.Printf("Testing pair %d and %d\n", a, b)

			gates[a].C, gates[b].C = gates[b].C, gates[a].C

			if checkBit(gates, n) {
				return []string{gates[a].C, gates[b].C}
			}

			gates[a].C, gates[b].C = gates[b].C, gates[a].C
		}
	}

	return nil
}

func extractAdder(gates []Gate, bit int) []int {
	adder := findGatesWithInput(gates, fmt.Sprintf("x%02d", bit))
	adder = append(adder, findGatesWithInput(gates, gates[adder[0]].C)...)
	adder = append(adder, findGatesWithInput(gates, gates[adder[1]].C)...)

out:
	for i, a := range adder {
		ac := gates[a].C
		for _, b := range adder[i+1:] {
			bc := gates[b].C

			result := findGatesWithInputs(gates, ac, bc)
			if len(result) > 0 {
				adder = append(adder, result...)
				break out
			}
		}
	}

	return adder
}

func findGatesWithInput(gates []Gate, needle string) []int {
	var results []int

	for i, gate := range gates {
		if gate.A == needle || gate.B == needle {
			results = append(results, i)
		}
	}

	return results
}

func findGatesWithInputs(gates []Gate, a, b string) []int {
	var results []int

	for i, gate := range gates {
		if gate.A == a && gate.B == b || gate.A == b && gate.B == a {
			results = append(results, i)
		}
	}

	return results
}
