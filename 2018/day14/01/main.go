package main

import (
	"flag"
	"fmt"
	"strings"
)

var number = flag.Int("n", 540561, "The number of recipes (your puzzle input).")

type lab struct {
	recipes    []int
	idx1, idx2 int
}

func (l lab) String() string {
	var sb strings.Builder
	for i, val := range l.recipes {
		switch i {
		case l.idx1:
			_, _ = fmt.Fprintf(&sb, "(%d)", val)
		case l.idx2:
			_, _ = fmt.Fprintf(&sb, "[%d]", val)
		default:
			_, _ = fmt.Fprintf(&sb, " %d ", val)
		}

	}
	return sb.String()
}

func (l lab) Len() int {
	return len(l.recipes)
}

func (l *lab) tick() {
	one, two := l.recipes[l.idx1], l.recipes[l.idx2]
	switch sum := one + two; {
	case sum < 10:
		l.recipes = append(l.recipes, sum)
	default:
		l.recipes = append(l.recipes, 1, sum%10)

	}

	l.idx1 = (l.idx1 + 1 + one) % l.Len()
	l.idx2 = (l.idx2 + 1 + two) % l.Len()
}

func (l *lab) last10(n int) string {
	for l.Len() < n+10 {
		l.tick()
	}
	var sb strings.Builder
	for _, val := range l.recipes[n : n+10] {
		_, _ = fmt.Fprintf(&sb, "%d", val)
	}
	return sb.String()
}

func main() {
	flag.Parse()

	lab := lab{[]int{3, 7}, 0, 1}
	fmt.Printf("Success! Target number is: %q\n", lab.last10(*number))
}
