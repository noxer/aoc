package main

import (
	"crypto/md5"
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

const key = "iwrupvqb"

func task1(_ []string) error {
	secret := []byte(key + strings.Repeat(" ", 9))[:len(key)]

	for n := int64(1); ; n++ {
		candidate := strconv.AppendInt(secret, n, 10)
		hash := md5.Sum(candidate)
		if hash[0] == 0 && hash[1] == 0 && hash[2] <= 0xf {
			fmt.Printf("Found hash at %d: %02x\n", n, hash[:])
			break
		}
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	secret := []byte(key + strings.Repeat(" ", 9))[:len(key)]

	for n := int64(1); ; n++ {
		candidate := strconv.AppendInt(secret, n, 10)
		hash := md5.Sum(candidate)
		if hash[0] == 0 && hash[1] == 0 && hash[2] == 0 {
			fmt.Printf("Found hash at %d: %02x\n", n, hash[:])
			break
		}
	}

	return nil
}
