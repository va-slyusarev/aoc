package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

const (
	ground     = '.'
	trees      = '|'
	lumberyard = '#'
)

type point struct {
	x, y int
}

func (p *point) dx(dx int) point { return point{p.x + dx, p.y} }

func (p *point) dy(dy int) point { return point{p.x, p.y + dy} }

func (p *point) String() string { return fmt.Sprintf("%d,%d", p.x, p.y) }

type collection struct {
	area map[point]rune

	lt, rb point
}

func (c *collection) String() string {
	var sb strings.Builder

	for y := c.lt.y; y <= c.rb.y; y++ {
		for x := c.lt.x; x <= c.rb.x; x++ {
			_, _ = fmt.Fprint(&sb, string(c.area[point{x, y}]))
		}
		_, _ = fmt.Fprintln(&sb)
	}
	return sb.String()
}

func (c *collection) tick(n int) {

	neighbors := func(p point) []rune {
		points := []point{{p.x - 1, p.y}, {p.x - 1, p.y - 1}, {p.x, p.y - 1}, {p.x + 1, p.y - 1}, {p.x + 1, p.y}, {p.x + 1, p.y + 1}, {p.x, p.y + 1}, {p.x - 1, p.y + 1}}
		result := make([]rune, 0)

		for _, cur := range points {
			if r, ok := c.area[cur]; ok {
				result = append(result, r)
			}
		}
		return result
	}

	mutator := func(p point) rune {

		switch c.area[p] {
		case ground:
			woods := 0
			for _, r := range neighbors(p) {
				if r == trees {
					woods++
				}
			}
			if woods > 2 {
				return trees
			}

		case trees:
			lumberyards := 0
			for _, r := range neighbors(p) {
				if r == lumberyard {
					lumberyards++
				}
			}
			if lumberyards > 2 {
				return lumberyard
			}
		case lumberyard:
			woods, lumberyards := 0, 0
			for _, r := range neighbors(p) {
				if r == trees {
					woods++
				}
				if r == lumberyard {
					lumberyards++
				}
			}
			if !(woods > 0 && lumberyards > 0) {
				return ground
			}

		}
		return c.area[p]
	}

	step := 0
	for step < n {

		newArea := make(map[point]rune)
		for p := range c.area {
			newArea[p] = mutator(p)
		}
		c.area = newArea

		step++
	}
}

func (c *collection) resource() int {
	woods, lumberyards := 0, 0
	for _, r := range c.area {
		switch r {
		case trees:
			woods++
		case lumberyard:
			lumberyards++
		}
	}
	return woods * lumberyards
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines []string

	for scanner.Scan() {
		value := scanner.Text()
		lines = append(lines, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	now := time.Now()

	collect := parse(lines)
	collect.tick(10)
	//fmt.Printf("%v", collect)
	fmt.Printf("Success! Target number is: %d (elapsed time: %s)\n", collect.resource(), time.Since(now))
}

func parse(lines []string) *collection {
	area := make(map[point]rune)
	cur := point{0, 0}

	for _, line := range lines {
		cur.x = 0
		for _, r := range line {
			area[cur] = r
			cur = cur.dx(1)
		}
		cur = cur.dy(1)
	}

	return &collection{
		area: area,
		lt:   point{0, 0},
		rb:   cur,
	}
}
