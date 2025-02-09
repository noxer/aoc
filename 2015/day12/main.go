package main

import (
	"encoding/json"
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

func parseDocument(name string) (any, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	d.UseNumber()

	var doc any
	err = d.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func sumNumbers(d any) int64 {
	sum := int64(0)

	switch t := d.(type) {
	case map[string]any:
		for _, v := range t {
			sum += sumNumbers(v)
		}
	case []any:
		for _, v := range t {
			sum += sumNumbers(v)
		}
	case json.Number:
		num, err := t.Int64()
		if err != nil {
			panic("Heeeelp: " + err.Error())
		}
		sum += num
	}

	return sum
}

func task1(args []string) error {
	doc, err := parseDocument(args[0])
	if err != nil {
		return err
	}

	sum := sumNumbers(doc)
	fmt.Printf("Sum: %d\n", sum)

	return nil
}

///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////

func hasProperty(obj map[string]any, prop any) bool {
	for _, v := range obj {
		if v == prop {
			return true
		}
	}

	return false
}

func sumNonRedNumbers(d any) int64 {
	sum := int64(0)

	switch t := d.(type) {
	case map[string]any:
		if hasProperty(t, "red") {
			break
		}

		for _, v := range t {
			sum += sumNonRedNumbers(v)
		}
	case []any:
		for _, v := range t {
			sum += sumNonRedNumbers(v)
		}
	case json.Number:
		num, err := t.Int64()
		if err != nil {
			panic("Heeeelp: " + err.Error())
		}
		sum += num
	}

	return sum
}

func task2(args []string) error {
	doc, err := parseDocument(args[0])
	if err != nil {
		return err
	}

	sum := sumNonRedNumbers(doc)
	fmt.Printf("Sum: %d\n", sum)

	return nil
}
