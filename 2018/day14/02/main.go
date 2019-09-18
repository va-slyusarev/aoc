package main

import (
	"flag"
	"fmt"
	"strconv"
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

func (l *lab) tick() (twice bool) {
	one, two := l.recipes[l.idx1], l.recipes[l.idx2]
	switch sum := one + two; {
	case sum < 10:
		l.recipes = append(l.recipes, sum)
		twice = false
	default:
		l.recipes = append(l.recipes, 1, sum%10)
		twice = true

	}

	l.idx1 = (l.idx1 + 1 + one) % l.Len()
	l.idx2 = (l.idx2 + 1 + two) % l.Len()
	return twice
}

func (l *lab) countLeft(n int) int {
	seq := make([]int, 0)
	for _, v := range strconv.Itoa(n) {
		seq = append(seq, int(v-48))
	}

	equals := func(a, b []int) bool {
		if a == nil || b == nil || len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}

	for {
		twice := l.tick()
		if l.Len() < len(seq)+1 {
			continue
		}
		switch twice {
		case true:
			if equals(l.recipes[l.Len()-len(seq)-1:l.Len()-1], seq) {
				return l.Len() - len(seq) - 1
			}
			fallthrough
		default:
			if equals(l.recipes[l.Len()-len(seq):], seq) {
				return l.Len() - len(seq)
			}
		}
	}
}

func main() {
	flag.Parse()

	lab := lab{[]int{3, 7}, 0, 1}
	fmt.Printf("Success! Target value is: %d\n", lab.countLeft(*number))
}
