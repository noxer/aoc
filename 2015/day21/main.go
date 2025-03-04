package main

import (
	"fmt"
	"os"
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

type Item struct {
	Cost   int
	Damage int
	Armor  int
}

func (i Item) Apply(d *Duelist) {
	d.Damage += i.Damage
	d.Armor += i.Armor
}

var (
	weapons = []Item{
		{8, 4, 0},
		{10, 5, 0},
		{25, 6, 0},
		{40, 7, 0},
		{74, 8, 0},
	}

	armor = []Item{
		{0, 0, 0},
		{13, 0, 1},
		{31, 0, 2},
		{53, 0, 3},
		{75, 0, 4},
		{102, 0, 5},
	}

	rings = []Item{
		{0, 0, 0},
		{0, 0, 0},
		{25, 1, 0},
		{50, 2, 0},
		{100, 3, 0},
		{20, 0, 1},
		{40, 0, 2},
		{80, 0, 3},
	}
)

type Duelist struct {
	Hitpoints int
	Damage    int
	Armor     int
}

func (d Duelist) Hit(other *Duelist) {
	points := d.Damage - other.Armor
	if points < 1 {
		points = 1
	}

	other.Hitpoints -= points
}

func fight(player, boss *Duelist) bool {
	for {
		player.Hit(boss)
		if boss.Hitpoints <= 0 {
			return true
		}

		boss.Hit(player)
		if player.Hitpoints <= 0 {
			return false
		}
	}
}

func task1(args []string) error {
	cost := 99999

	for _, w := range weapons {
		for _, a := range armor {
			for i, r1 := range rings {
				for _, r2 := range rings[i+1:] {
					player := Duelist{100, 0, 0}
					w.Apply(&player)
					a.Apply(&player)
					r1.Apply(&player)
					r2.Apply(&player)

					boss := Duelist{100, 8, 2}

					if fight(&player, &boss) {
						cost = min(cost, w.Cost+a.Cost+r1.Cost+r2.Cost)
					}
				}
			}
		}
	}

	fmt.Printf("Best fight for money: %d\n", cost)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	cost := 0

	for _, w := range weapons {
		for _, a := range armor {
			for i, r1 := range rings {
				for _, r2 := range rings[i+1:] {
					player := Duelist{100, 0, 0}
					w.Apply(&player)
					a.Apply(&player)
					r1.Apply(&player)
					r2.Apply(&player)

					boss := Duelist{100, 8, 2}

					if !fight(&player, &boss) {
						cost = max(cost, w.Cost+a.Cost+r1.Cost+r2.Cost)
					}
				}
			}
		}
	}

	fmt.Printf("Worst fight for money: %d\n", cost)

	return nil
}
