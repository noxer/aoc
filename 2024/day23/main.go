package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
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

type Computer struct {
	Name  string
	Peers []string
}

func (c *Computer) String() string {
	return fmt.Sprintf("%s: %v", c.Name, c.Peers)
}

func (c *Computer) AddPeer(name string) {
	c.Peers = append(c.Peers, name)
}

func (c *Computer) HasPeer(name string) bool {
	return slices.Contains(c.Peers, name)
}

type Network struct {
	Computers map[string]*Computer
}

func (n Network) GetComputer(name string) *Computer {
	if c, ok := n.Computers[name]; ok {
		return c
	}

	c := &Computer{
		Name: name,
	}
	n.Computers[name] = c

	return c
}

func (n Network) ParseConnection(line string) {
	a, b, _ := strings.Cut(line, "-")
	n.GetComputer(a).AddPeer(b)
	n.GetComputer(b).AddPeer(a)
}

func (n Network) Parse(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	s := bufio.NewScanner(f)

	for s.Scan() {
		n.ParseConnection(s.Text())
	}

	return s.Err()
}

func (n Network) FindSetsOfThree() [][]string {
	seen := make(map[string]struct{})
	var results [][]string

	for _, c := range n.Computers {
		for i, p := range c.Peers {
			peer := n.GetComputer(p)
			for _, o := range c.Peers[i+1:] {
				if peer.HasPeer(o) {
					triangle := normalize(c.Name, p, o)
					key := strings.Join(triangle, ",")
					if _, ok := seen[key]; ok {
						continue
					}
					seen[key] = struct{}{}

					results = append(results, triangle)
				}
			}
		}
	}

	return results
}

func normalize(c, p, o string) []string {
	r := []string{c, p, o}
	sort.Strings(r)
	return r
}

func filterStartsWithT(triangles [][]string) [][]string {
	result := triangles[:0]

	for _, triangle := range triangles {
		for _, name := range triangle {
			if strings.HasPrefix(name, "t") {
				result = append(result, triangle)
				break
			}
		}
	}

	return result
}

func task1(args []string) error {
	network := Network{
		Computers: make(map[string]*Computer),
	}

	err := network.Parse(args[0])
	if err != nil {
		return err
	}

	triangles := network.FindSetsOfThree()
	triangles = filterStartsWithT(triangles)

	fmt.Printf("Count of three connected computer triangles: %d\n", len(triangles))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func (n Network) CanMerge(as, bs []string) bool {
	if len(as) == 0 || len(bs) == 0 {
		return false
	}

	for _, a := range as {
		c := n.GetComputer(a)
		for _, b := range bs {
			if a != b && !c.HasPeer(b) {
				return false
			}
		}
	}

	return true
}

func (n Network) Merge(as, bs []string) []string {
	result := make([]string, 0, len(as)+len(bs))
	result = append(result, as...)

	for _, b := range bs {
		if slices.Contains(as, b) {
			continue
		}

		result = append(result, b)
	}

	return result
}

func (n Network) FindBiggestClusterForReal() []string {
	computers := n.FindSetsOfThree()

	return n.findBiggest(computers)
}

func (n Network) findBiggest(computers [][]string) []string {
	var best []string

	fmt.Println(countNonNil(computers))

	for i, as := range computers {
		if len(as) == 0 {
			continue
		}

		for j, bs := range computers[i+1:] {
			if !n.CanMerge(as, bs) {
				continue
			}

			computers[i] = n.Merge(as, bs)
			computers[i+1+j] = nil

			candidate := n.findBiggest(computers)
			if len(candidate) > len(best) {
				best = candidate
			}

			computers[i+1+j] = bs
		}

		computers[i] = as
	}

	if len(best) == 0 {
		for _, computer := range computers {
			if len(computer) > len(best) {
				best = computer
			}
		}

		sort.Strings(best)
		fmt.Printf("Best: %s\n", strings.Join(best, ","))
	}

	return best
}

func countNonNil(computers [][]string) int {
	count := 0
	for _, computer := range computers {
		if len(computer) != 0 {
			count++
		}
	}
	return count
}

func (n Network) FindBiggestCluster() []string {
	computers := make([][]string, 0, len(n.Computers))
	for name := range n.Computers {
		computers = append(computers, []string{name})
	}

	seen := make(map[string]struct{})
	wasMerged := make(map[string]struct{})
	var next [][]string
	for {
		next = next[:0]
		clear(wasMerged)
		clear(seen)

		for i, as := range computers {
			// key := strings.Join(as, ",")

			for _, bs := range computers[i+1:] {
				if !n.CanMerge(as, bs) {
					continue
				}

				merged := n.Merge(as, bs)
				sort.Strings(merged)
				key := strings.Join(merged, ",")
				if _, ok := seen[key]; ok {
					continue
				}
				seen[key] = struct{}{}

				// mark both clusters as merged
				wasMerged[key] = struct{}{}
				wasMerged[strings.Join(bs, ",")] = struct{}{}

				next = append(next, n.Merge(as, bs))
			}
		}

		if len(next) == 0 {
			slices.SortFunc(computers, func(as, bs []string) int {
				return len(as) - len(bs)
			})

			return computers[0]
		}

		fmt.Println("-", next)

		computers, next = next, computers
	}
}

func (n Network) FindBiggestClusterTheNextGeneration() []string {
	computers := n.FindSetsOfThree()

	for len(computers) > 1 {
		slices.SortFunc(computers, func(as, bs []string) int {
			return len(bs) - len(as)
		})

		last := computers[len(computers)-1]
		computers = computers[:len(computers)-1]

		for i, as := range computers {
			if !n.CanMerge(last, as) {
				continue
			}

			computers[i] = n.Merge(last, as)
		}
	}

	return computers[0]
}

func task2(args []string) error {
	network := Network{
		Computers: make(map[string]*Computer),
	}

	err := network.Parse(args[0])
	if err != nil {
		return err
	}

	biggest := network.FindBiggestClusterTheNextGeneration()
	fmt.Println(biggest)

	sort.Strings(biggest)
	password := strings.Join(biggest, ",")

	fmt.Printf("Password: %s\n", password)

	// fmt.Printf("Count of three connected computer triangles: %d\n", len(triangles))

	return nil
}
