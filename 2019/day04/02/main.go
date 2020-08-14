package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
)

var min = flag.Int("min", 123257, "Minimal value in range")
var max = flag.Int("max", 647015, "Maximal value in range")

type digits [6]int

func main() {
	flag.Parse()

	count := 0
	for value := *min; value <= *max; value++ {

		d := parseDigits(value)
		if sameDigits(d) && neverDecrease(d) {
			count++
		}
	}

	fmt.Printf("Success! Target count is: %d\n", count)
}

func parseDigits(value int) digits {
	getN := func(value int, idx int) int {
		s := strconv.Itoa(value)
		if len(s) < idx {
			log.Fatalf("error parsing digits: number is %d, digit index is %d\n", value, idx)
		}
		i, err := strconv.Atoi(s[idx-1 : idx])
		if err != nil {
			log.Fatalf("error parsing digits: %v\n", err)
		}
		return i
	}

	return digits{getN(value, 1), getN(value, 2), getN(value, 3), getN(value, 4), getN(value, 5), getN(value, 6)}
}

func sameDigits(d digits) bool {
	return (d[0] == d[1] && d[1] != d[2]) ||
		(d[1] == d[2] && d[2] != d[3] && d[2] != d[0]) ||
		(d[2] == d[3] && d[3] != d[4] && d[3] != d[1]) ||
		(d[3] == d[4] && d[4] != d[5] && d[4] != d[2]) ||
		(d[4] == d[5] && d[5] != d[3])
}

func neverDecrease(d digits) bool {
	pr := d[0]
	for i := 1; i < 6; i++ {
		if d[i] < pr {
			return false
		}
		pr = d[i]
	}
	return true
}
