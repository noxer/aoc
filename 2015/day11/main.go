package main

import (
	"fmt"
	"os"
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

const alphabet = "abcdefghjkmnpqrstuvwxyz"

func nextPassword(pass []byte, check func([]byte) bool) []byte {
	if recursivePassword(pass, pass, check) {
		return pass
	}

	return nil
}

func recursivePassword(full, pass []byte, check func([]byte) bool) bool {
	if len(pass) == 0 {
		return check(full)
	}

	i := strings.IndexByte(alphabet, pass[0])

	for _, b := range []byte(alphabet[i:]) {
		pass[0] = b
		if recursivePassword(full, pass[1:], check) {
			return true
		}
	}

	pass[0] = alphabet[0]

	return false
}

func checkStraight(pass []byte) bool {
	for i := range pass[:len(pass)-2] {
		if pass[i] == pass[i+1]-1 && pass[i] == pass[i+2]-2 {
			return true
		}
	}

	return false
}

func checkDoubles(pass []byte) bool {
	first := byte(0)

	for i := range pass[:len(pass)-1] {
		if pass[i] == pass[i+1] {
			if first != 0 && first != pass[i] {
				return true
			}

			first = pass[i]
		}
	}

	return false
}

func checkPass(pass []byte) bool {
	return checkStraight(pass) && checkDoubles(pass)
}

func task1(_ []string) error {
	pass := nextPassword([]byte("vzbxxyza"), checkPass)
	fmt.Println(string(pass))

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func task2(args []string) error {
	pass := nextPassword([]byte("vzbxxzaa"), checkPass)
	fmt.Println(string(pass))

	return nil
}
