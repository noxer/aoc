package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"strings"
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

type Ingredient map[string]int

func parseIngredient(line string) (string, Ingredient) {
	name, props, _ := strings.Cut(line, ": ")
	properties := strings.Split(props, ", ")

	in := Ingredient{}

	for _, prop := range properties {
		key, value, _ := strings.Cut(prop, " ")
		val, _ := strconv.Atoi(value)

		in[key] = val
	}

	return name, in
}

func parseFile(name string) (map[string]Ingredient, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	m := make(map[string]Ingredient)
	for s.Scan() {
		name, in := parseIngredient(s.Text())
		m[name] = in
	}

	return m, s.Err()
}

func (in Ingredient) Score() int {
	prod := 1
	for _, name := range []string{"capacity", "durability", "flavor", "texture"} {
		if in[name] <= 0 {
			return 0
		}
		prod *= in[name]
	}
	return prod
}

func (in Ingredient) AddMul(other Ingredient, scalar int) Ingredient {
	in2 := maps.Clone(in)

	for k, v := range other {
		in2[k] += v * scalar
	}

	return in2
}

func score(ingredients []Ingredient, current Ingredient, remaining int) int {
	if remaining == 0 || len(ingredients) == 0 {
		return current.Score()
	}
	if len(ingredients) == 1 {
		return current.AddMul(ingredients[0], remaining).Score()
	}

	best := -1
	for spoons := range remaining + 1 {
		next := current.AddMul(ingredients[0], spoons)
		best = max(best, score(ingredients[1:], next, remaining-spoons))
	}
	return best
}

func task1(args []string) error {
	ingredients, err := parseFile(args[0])
	if err != nil {
		return err
	}

	values := slices.Collect(maps.Values(ingredients))
	best := score(values, make(Ingredient), 100)

	fmt.Printf("Best: %d\n", best)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (in Ingredient) ScoreFor500Cal() int {
	if in["calories"] == 500 {
		return in.Score()
	}

	return -1
}

func scoreCalories(ingredients []Ingredient, current Ingredient, remaining int) int {
	if remaining == 0 || len(ingredients) == 0 {
		return current.ScoreFor500Cal()
	}
	if len(ingredients) == 1 {
		return current.AddMul(ingredients[0], remaining).ScoreFor500Cal()
	}

	best := -1
	for spoons := range remaining + 1 {
		next := current.AddMul(ingredients[0], spoons)
		best = max(best, scoreCalories(ingredients[1:], next, remaining-spoons))
	}
	return best
}

func task2(args []string) error {
	ingredients, err := parseFile(args[0])
	if err != nil {
		return err
	}

	values := slices.Collect(maps.Values(ingredients))
	best := scoreCalories(values, make(Ingredient), 100)

	fmt.Printf("Best: %d\n", best)

	return nil
}
