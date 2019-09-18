package main

import (
	"fmt"
)

func main() {

	fmt.Printf("Success! Target number is: %d\n", manually())
}

func manually() int {
	curr := 0
	last := 0
	seen := make(map[int]struct{})
	for {
		prev := curr | 65536
		curr = 1505483

		for {
			curr = (curr + prev&255) * 65899 & 16777215
			if prev < 256 {
				break
			}
			prev /= 256
		}

		if _, ok := seen[curr]; ok {
			return last
		}
		seen[curr] = struct{}{}
		last = curr
	}
}
