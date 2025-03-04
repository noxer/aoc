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

type Spell struct {
	Cost int

	InstantDamage int
	InstantHeal   int

	EffectTurns  int
	EffectDamage int
	EffectArmor  int
	EffectMana   int
}

var spells = []Spell{
	{53, 4, 0, 0},
	{},
	{},
	{},
	{},
}

type Duelist struct {
	Hitpoints int
	Mana      int
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

func task1(args []string) error {
	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	return nil
}
