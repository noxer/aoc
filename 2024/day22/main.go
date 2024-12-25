package main

import (
	"bytes"
	"fmt"
	"iter"
	"os"
	"slices"
	"strconv"
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

func nextRandom(n uint) uint {
	n ^= n * 64
	n %= 16777216

	n ^= n / 32
	n %= 16777216

	n ^= n * 2048
	n %= 16777216

	return n
}

func calculateRandom(n uint, rounds int) uint {
	for range rounds {
		n = nextRandom(n)
	}

	return n
}

func task1(args []string) error {
	buyers, err := utils.ReadLinesTransform(args[0], func(line string) uint {
		i, _ := strconv.ParseUint(line, 10, 0)
		return uint(i)
	})
	if err != nil {
		return err
	}

	sum := 0
	for _, buyer := range buyers {
		sum += int(calculateRandom(buyer, 2000))
	}

	fmt.Printf("Sum: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

type Random uint

func (r Random) IteratePrices(n int) func(func(uint) bool) {
	return func(yield func(uint) bool) {
		num := uint(r)

		for range n {
			num = nextRandom(num)
			if !yield(num % 10) {
				return
			}
		}
	}
}

func (r Random) IterateOffsets(n int) func(func(byte) bool) {
	return func(yield func(byte) bool) {
		num := uint(r)
		last := num

		for range n {
			num = nextRandom(num)
			offset := byte(int8(num%10) - int8(last%10))
			if !yield(offset) {
				return
			}

			last = num
		}
	}
}

func task2(args []string) error {
	buyers, err := utils.ReadLinesTransform(args[0], func(line string) uint {
		i, _ := strconv.ParseUint(line, 10, 0)
		return uint(i)
	})
	if err != nil {
		return err
	}

	start := time.Now()

	sequences := make([][]byte, len(buyers))
	prices := make([][]uint, len(buyers))
	for i, buyer := range buyers {
		sequences[i] = slices.Collect(iter.Seq[byte](Random(buyer).IterateOffsets(2000)))
		prices[i] = slices.Collect(iter.Seq[uint](Random(buyer).IteratePrices(2000)))
	}

	// maximumPrice := uint(0)
	// for a := int8(-9); a <= 9; a++ {
	// 	for b := int8(-9); b <= 9; b++ {
	// 		for c := int8(-9); c <= 9; c++ {
	// 			for d := int8(-9); d <= 9; d++ {
	// 				price := sumPrices(prices, sequences, []byte{byte(a), byte(b), byte(c), byte(d)})
	// 				if price > maximumPrice {
	// 					fmt.Printf("Price for %d, %d, %d, %d: %d\n", a, b, c, d, price)
	// 					maximumPrice = price
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	maxPrice := uint(0)

	ps := bestPrices(prices, sequences)
	for _, p := range ps {
		maxPrice = max(maxPrice, p)
	}

	// fmt.Println(ps)

	elapsed := time.Since(start)

	fmt.Printf("Maximum price: %d (%s)\n", maxPrice, elapsed)

	return nil
}

func sumPrices(prices [][]uint, sequences [][]byte, sub []byte) uint {
	price := uint(0)
	for i, sequence := range sequences {
		idx := bytes.Index(sequence, sub)
		if idx < 0 {
			continue
		}

		idx += len(sub) - 1
		price += prices[i][idx]
	}
	return price
}

func bestPrices(prices [][]uint, sequences [][]byte) map[uint32]uint {
	ps := make(map[uint32]uint)

	for i, sequence := range sequences {
		bestPrice(ps, prices[i], sequence)
	}

	return ps
}

func bestPrice(ps map[uint32]uint, prices []uint, sequence []byte) {
	seen := make(map[uint32]struct{})

	for i := range sequence[:len(sequence)-3] {
		key := encodeBytes(
			sequence[i],
			sequence[i+1],
			sequence[i+2],
			sequence[i+3],
		)

		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		ps[key] += prices[i+3]
	}
}

func encodeBytes(a, b, c, d byte) uint32 {
	return uint32(a) | uint32(b)<<8 | uint32(c)<<16 | uint32(d)<<24
}
