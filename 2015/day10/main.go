package main

import (
	"fmt"
	"os"
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

// 1122331

func lookSay(seq string) string {
	sb := strings.Builder{}

	for i := 0; i < len(seq); i++ {
		b := seq[i]
		count := 0
		for _, p := range []byte(seq[i:]) {
			if p == b {
				count++
			} else {
				break
			}
		}

		sb.WriteString(strconv.Itoa(count))
		sb.WriteByte(b)

		i += count - 1
	}

	return sb.String()
}

func task1(_ []string) error {
	seq := "1113222113"

	for range 40 {
		seq = lookSay(seq)
	}

	fmt.Printf("Length: %d\n", len(seq))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(_ []string) error {
	seq := "1113222113"

	for range 50 {
		seq = lookSay(seq)
	}

	fmt.Printf("Length: %d\n", len(seq))

	return nil
}
