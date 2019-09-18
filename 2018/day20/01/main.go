package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

var dirs = map[rune]point{'N': {0, -1}, 'E': {1, 0}, 'W': {-1, 0}, 'S': {0, 1}}

type point struct{ x, y int }

func (p *point) offset(xy point) point { return point{p.x + xy.x, p.y + xy.y} }

type layout map[point]int

func (l *layout) maxDist() int {
	max := 0
	for _, r := range *l {
		if r > max {
			max = r
		}
	}
	return max
}

func main() {
	flag.Parse()

	f, err := os.Open(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	defer func() { _ = f.Close() }()

	var regexp string
	if _, err := fmt.Fscanln(f, &regexp); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	now := time.Now()
	model := build(regexp)
	fmt.Printf("Success! Target number is: %d (elapsed time: %s)\n", model.maxDist(), time.Since(now))
}

func build(regexp string) layout {
	model := make(layout)
	walk(regexp, point{0, 0}, 0, model)
	return model
}

func walk(regexp string, p point, dist int, model layout) string {
	for {
		r := rune(regexp[0])
		switch r {
		case '^':
			regexp = regexp[1:]
		case 'N', 'E', 'W', 'S':
			p = p.offset(dirs[r])
			dist += 1
			if d, ok := model[p]; !ok || d > dist {
				model[p] = dist
			}
			regexp = regexp[1:]
		case '(':
			regexp = walkAnother(regexp, p, dist, model)
		case '|', ')':
			return regexp

		default:
			return regexp[1:]
		}
	}
}

func walkAnother(regexp string, p point, dist int, model layout) string {
	regexp = regexp[1:]
	for {
		regexp = walk(regexp, p, dist, model)
		r := rune(regexp[0])
		switch r {
		case '|':
			regexp = regexp[1:]
		case ')':
			return regexp[1:]
		}
	}
}
