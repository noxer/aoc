package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseInts(str, sep string) []int {
	parts := strings.Split(str, sep)
	ints := make([]int, len(parts))

	for i, part := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			panic(fmt.Sprintf("Error trying to parse %q (index %d) as int: %s", part, i, err))
		}

		ints[i] = n
	}

	return ints
}
