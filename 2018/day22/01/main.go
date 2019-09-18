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
	mouth  = 'M'
	target = 'T'
	rocky  = '.'
	wet    = '='
	narrow = '|'
)

type point struct {
	x, y int
}

func (p *point) dx(dx int) point { return point{p.x + dx, p.y} }

func (p *point) dy(dy int) point { return point{p.x, p.y + dy} }

func (p *point) equals(value *point) bool { return value != nil && p.x == value.x && p.y == value.y }

func (p *point) String() string { return fmt.Sprintf("%d,%d", p.x, p.y) }

type geology struct {
	geologic, erosion int
	symbol            rune
}

type cave struct {
	area   map[point]geology
	lt, rb point
}

func (c *cave) dangerousValue(depth int, target point) int {
	c.area = make(map[point]geology)
	c.lt = point{0, 0}
	c.rb = target

	value := 0
	for y := c.lt.y; y <= c.rb.y; y++ {
		for x := c.lt.x; x <= c.rb.x; x++ {
			cur := &point{x, y}
			geo := geology{}

			switch {
			case cur.equals(&c.lt) || cur.equals(&c.rb):
				geo.geologic = 0
			case cur.y == 0:
				geo.geologic = cur.x * 16807
			case cur.x == 0:
				geo.geologic = cur.y * 48271
			default:
				geo.geologic = c.area[cur.dx(-1)].erosion * c.area[cur.dy(-1)].erosion
			}

			geo.erosion = (geo.geologic + depth) % 20183

			switch geo.erosion % 3 {
			case 0:
				geo.symbol = rocky
				value += 0
			case 1:
				geo.symbol = wet
				value += 1
			case 2:
				geo.symbol = narrow
				value += 2
			}

			c.area[*cur] = geo
		}
	}

	return value
}

func (c *cave) String() string {
	var sb strings.Builder

	for y := c.lt.y; y <= c.rb.y; y++ {
		for x := c.lt.x; x <= c.rb.x; x++ {
			if c.lt.equals(&point{x, y}) {
				_, _ = fmt.Fprint(&sb, string(mouth))
				continue
			}
			if c.rb.equals(&point{x, y}) {
				_, _ = fmt.Fprint(&sb, string(target))
				continue
			}
			_, _ = fmt.Fprint(&sb, string(c.area[point{x, y}].symbol))
		}
		_, _ = fmt.Fprintln(&sb)
	}
	return sb.String()
}

func main() {
	now := time.Now()

	depth, target, err := parse()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	c := new(cave)
	val := c.dangerousValue(depth, *target)
	//fmt.Printf("%v\n", c)
	fmt.Printf("Success! Target number is: %d (elapsed time: %s)\n", val, time.Since(now))
}

func parse() (int, *point, error) {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return 0, nil, fmt.Errorf("error read input data file: %v", err)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))
	var lines []string
	for scanner.Scan() {
		value := scanner.Text()
		lines = append(lines, value)
	}

	if err := scanner.Err(); err != nil {
		return 0, nil, fmt.Errorf("error scan input data file: %v", err)
	}

	if len(lines) != 2 {
		return 0, nil, fmt.Errorf("error scan input data file: incorrect len 2 != %d", len(lines))
	}

	var depth int
	var target point

	if _, err := fmt.Sscanf(lines[0], "depth: %d", &depth); err != nil {
		return 0, nil, fmt.Errorf("error scan input data file: incorrect depth: %v", err)
	}
	if _, err := fmt.Sscanf(lines[1], "target: %d,%d", &target.x, &target.y); err != nil {
		return 0, nil, fmt.Errorf("error scan input data file: incorrect target: %v", err)
	}
	return depth, &target, nil
}
