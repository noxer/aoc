package main

import (
	"fmt"
	"os"
	"sort"
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

func task1(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	_, modules := loadModules(lines)
	initializeConjunctionModules(modules)

	gq := &GlobalQueue{}

	for range 1000 {
		gq.Process(modules)
	}

	low, high := gq.Counters()
	fmt.Printf("Low: %d; High: %d; Product: %d\n", low, high, low*high)

	return nil
}

func initializeConjunctionModules(modules map[string]Module) {
	for name, module := range modules {
		outputs := module.Outputs()

		for _, output := range outputs {
			mod := modules[output]
			if conj, ok := mod.(*ConjunctionModule); ok {
				conj.memory[name] = Low
			}
		}
	}
}

func loadModules(lines []string) ([]string, map[string]Module) {
	modules := make(map[string]Module)
	var names []string

	for _, line := range lines {
		name, rawOutputs, _ := strings.Cut(line, " -> ")
		outputs := strings.Split(rawOutputs, ", ")

		if name == "broadcaster" {
			modules[name] = &BroadcasterModule{outputs: outputs}
			names = append(names, name)
			continue
		}

		switch name[0] {
		case '%':
			modules[name[1:]] = &FlipFlopModule{
				Name:    name[1:],
				outputs: outputs,
			}
			names = append(names, name[1:])
		case '&':
			modules[name[1:]] = &ConjunctionModule{
				name:    name[1:],
				memory:  make(map[string]Pulse),
				outputs: outputs,
			}
			names = append(names, name[1:])
		}
	}

	return names, modules
}

type Pulse bool

const (
	High = Pulse(true)
	Low  = Pulse(false)
)

type PendingPulse struct {
	Source      string
	Destination string
	Pulse       Pulse
}

type GlobalQueue struct {
	lowPulses  int
	highPulses int
	queue      []PendingPulse
}

func (gq *GlobalQueue) Enqueue(pulse Pulse, source string, destinations ...string) {
	for _, destination := range destinations {
		gq.queue = append(gq.queue, PendingPulse{
			Source:      source,
			Destination: destination,
			Pulse:       pulse,
		})
	}
}

func (gq *GlobalQueue) Counters() (int, int) {
	return gq.lowPulses, gq.highPulses
}

func (gq *GlobalQueue) Process(modules map[string]Module) {
	gq.queue = append(gq.queue, PendingPulse{
		Source:      "button",
		Destination: "broadcaster",
		Pulse:       Low,
	})

	for len(gq.queue) > 0 {
		pending := gq.queue[0]
		gq.queue = gq.queue[1:]
		if pending.Pulse {
			gq.highPulses++
		} else {
			gq.lowPulses++
		}

		fmt.Printf("%s -%t-> %s\n", pending.Source, pending.Pulse, pending.Destination)

		if module, ok := modules[pending.Destination]; ok {
			module.Trigger(gq, pending.Source, pending.Pulse)
		}
	}
}

type Module interface {
	Trigger(eq Enqueuer, source string, pulse Pulse)
	Outputs() []string
}

type FlipFlopModule struct {
	Name    string
	State   Pulse
	outputs []string
}

func (m *FlipFlopModule) Trigger(eq Enqueuer, source string, pulse Pulse) {
	if pulse == High {
		return
	}

	m.State = !m.State
	eq.Enqueue(m.State, m.Name, m.outputs...)
}

func (m *FlipFlopModule) Outputs() []string {
	return m.outputs
}

type ConjunctionModule struct {
	name    string
	memory  map[string]Pulse
	outputs []string
}

func (m *ConjunctionModule) Trigger(eq Enqueuer, source string, pulse Pulse) {
	m.memory[source] = pulse

	collected := High
	for _, remembered := range m.memory {
		collected = collected && remembered
	}

	eq.Enqueue(!collected, m.name, m.outputs...)
}

func (m *ConjunctionModule) Outputs() []string {
	return m.outputs
}

type BroadcasterModule struct {
	outputs []string
}

func (m *BroadcasterModule) Trigger(eq Enqueuer, source string, pulse Pulse) {
	eq.Enqueue(pulse, "broadcaster", m.outputs...)
}

func (m *BroadcasterModule) Outputs() []string {
	return m.outputs
}

///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	lines, err := utils.ReadLines(args[0])
	if err != nil {
		return err
	}

	_, modules := loadModules(lines)
	initializeConjunctionModules(modules)

	gq := &GlobalQueue2{}

	// state, bits := moduleState(names, modules)
	// fmt.Printf("%d (%048b)\n", len(state), bits)

	counter := 0
	for !gq.Check() {
		counter++

		gq.Process(counter, modules)

		// bits, count := gq.Number()
		// fmt.Printf("Bits: % 8d (%064b)\n", bits, bits)
		// fmt.Printf("                %s%s\n", strings.Repeat(" ", 64-count), strings.Repeat("-", count))

		// _, bits := moduleState(names, modules)
		// fmt.Printf("Iterated %d times (%048b)\r", counter, bits)
	}

	fmt.Printf("Counter: %d\n", counter)

	return nil
}

type Enqueuer interface {
	Enqueue(pulse Pulse, source string, destinations ...string)
}

type GlobalQueue2 struct {
	rxCount int
	rxBits  uint64
	queue   []PendingPulse
}

func moduleState(names []string, modules map[string]Module) ([]Pulse, uint64) {
	states := make([]Pulse, 0, len(modules))
	bits := uint64(0)
	for _, name := range names {
		module := modules[name]

		switch t := module.(type) {
		case *FlipFlopModule:
			states = append(states, t.State)
			bits <<= 1
			if t.State {
				bits |= 1
			}
		case *ConjunctionModule:
			continue
			if len(t.memory) <= 1 {
				continue
			}
			inputNames := sortedKeys(t.memory)
			for _, name := range inputNames {
				states = append(states, t.memory[name])
			}
		}
	}
	return states, bits
}

func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (gq *GlobalQueue2) Enqueue(pulse Pulse, source string, destinations ...string) {
	for _, destination := range destinations {
		gq.queue = append(gq.queue, PendingPulse{
			Source:      source,
			Destination: destination,
			Pulse:       pulse,
		})
	}
}

func (gq *GlobalQueue2) Check() bool {
	return gq.rxCount == 1 && gq.rxBits == 0
}

func (gq *GlobalQueue2) Number() (uint64, int) {
	return gq.rxBits, gq.rxCount
}

func (gq *GlobalQueue2) Process(counter int, modules map[string]Module) {
	gq.queue = append(gq.queue, PendingPulse{
		Source:      "button",
		Destination: "broadcaster",
		Pulse:       Low,
	})
	gq.rxCount = 0
	gq.rxBits = 0

	for len(gq.queue) > 0 {
		pending := gq.queue[0]
		copy(gq.queue, gq.queue[1:])
		gq.queue = gq.queue[:len(gq.queue)-1]

		// fmt.Printf("%s -%t-> %s\n", pending.Source, pending.Pulse, pending.Destination)

		if pending.Pulse == Low {
			switch pending.Destination {
			case "hz":
				fmt.Println(counter, "Section HP sends High!")
				panic("stop")
			case "pv":
				fmt.Println(counter, "Section QN sends High!")
			case "qh":
				fmt.Println(counter, "Section XV sends High!")
			case "xm":
				fmt.Println(counter, "Section ZB sends High!")
			}
		}

		if pending.Destination == "rx" {
			gq.rxCount++
			gq.rxBits <<= 1
			if pending.Pulse == High {
				gq.rxBits |= 1
			}
		}

		if module, ok := modules[pending.Destination]; ok {
			module.Trigger(gq, pending.Source, pending.Pulse)
		}
	}
}
