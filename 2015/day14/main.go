package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

var matchReindeer = regexp.MustCompile(`^([A-Za-z]+) can fly ([0-9]+) km/s for ([0-9]+) seconds, but then must rest for ([0-9]+) seconds.$`)

type Reindeer struct {
	Speed        int
	FlyDuration  int
	RestDuration int
}

func (r Reindeer) DistanceAfterSeconds(seconds int) int {
	cycleTime := r.FlyDuration + r.RestDuration
	cycleDist := r.Speed * r.FlyDuration

	distance := (seconds / cycleTime) * cycleDist
	seconds = seconds % cycleTime

	if seconds >= r.FlyDuration {
		distance += cycleDist
		return distance
	}

	distance += seconds * r.Speed
	return distance
}

func parseReindeer(line string) (string, Reindeer) {
	match := matchReindeer.FindStringSubmatch(line)
	if len(match) != 5 {
		fmt.Printf("Couldn't match line %s\n", line)
		return "", Reindeer{}
	}

	r := Reindeer{}
	r.Speed, _ = strconv.Atoi(match[2])
	r.FlyDuration, _ = strconv.Atoi(match[3])
	r.RestDuration, _ = strconv.Atoi(match[4])

	return match[1], r
}

func parseFile(name string) (map[string]Reindeer, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	rs := make(map[string]Reindeer)
	for s.Scan() {
		name, reindeer := parseReindeer(s.Text())
		rs[name] = reindeer
	}

	return rs, s.Err()
}

func task1(args []string) error {
	reindeers, err := parseFile(args[0])
	if err != nil {
		return err
	}

	best := 0
	seconds := 2503
	for _, reindeer := range reindeers {
		best = max(best, reindeer.DistanceAfterSeconds(seconds))
	}

	fmt.Printf("Best distance: %d\n", best)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	reindeers, err := parseFile(args[0])
	if err != nil {
		return err
	}

	points := make(map[string]int)
	for seconds := range 2503 {
		best := 0

		for _, reindeer := range reindeers {
			best = max(best, reindeer.DistanceAfterSeconds(seconds+1))
		}

		for name, reindeer := range reindeers {
			if reindeer.DistanceAfterSeconds(seconds+1) == best {
				points[name]++
			}
		}
	}

	for name, ps := range points {
		fmt.Printf("%s: %d\n", name, ps)
	}

	return nil
}
