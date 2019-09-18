package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
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
const (
	torch = iota
	gear
	neither
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

	coord    point
	tool     int // current tool
	distance int // to target
}

func (g geology) String() string {
	return fmt.Sprintf("[%v], tool: %d dist: %d", &g.coord, g.tool, g.distance)
}

func (g geology) hash() string {
	return fmt.Sprintf("%v, tool: %d", &g.coord, g.tool)
}

type geos []geology

func (g geos) Len() int           { return len(g) }
func (g geos) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g geos) Less(i, j int) bool { return g[i].distance < g[j].distance } // asc

type cave struct {
	area   map[point]geology
	lt, rb point
	target point
}

func (c *cave) geodesy(depth int, target point) {
	c.area = make(map[point]geology)
	c.lt = point{0, 0}
	c.rb = point{target.x + 40, target.y + 40} // expand for another routes
	c.target = target

	for y := c.lt.y; y <= c.rb.y; y++ {
		for x := c.lt.x; x <= c.rb.x; x++ {
			geo := geology{coord: point{x, y}}

			switch {
			case geo.coord.equals(&c.lt) || geo.coord.equals(&c.rb):
				geo.geologic = 0
			case geo.coord.y == 0:
				geo.geologic = geo.coord.x * 16807
			case geo.coord.x == 0:
				geo.geologic = geo.coord.y * 48271
			default:
				geo.geologic = c.area[geo.coord.dx(-1)].erosion * c.area[geo.coord.dy(-1)].erosion
			}

			geo.erosion = (geo.geologic + depth) % 20183

			switch geo.erosion % 3 {
			case 0:
				geo.symbol = rocky
			case 1:
				geo.symbol = wet
			case 2:
				geo.symbol = narrow
			}

			c.area[geo.coord] = geo
		}
	}
}

func (c *cave) minWalk() int {

	allowedTool := func(geo geology, tool int) bool {
		return (geo.symbol == rocky && tool != neither) ||
			(geo.symbol == wet && tool != torch) ||
			(geo.symbol == narrow && tool != gear)
	}

	neighbors := func(g geology) geos {
		points := []point{{g.coord.x - 1, g.coord.y}, {g.coord.x, g.coord.y - 1}, {g.coord.x + 1, g.coord.y}, {g.coord.x, g.coord.y + 1}}
		result := make(geos, 0)

		for _, cur := range points {
			if r, ok := c.area[cur]; ok {

				if cur.equals(&c.lt) || cur.equals(&c.target) {
					r.tool = torch
					result = append(result, r)
					continue
				}
				for tool := torch; tool <= neither; tool++ {
					if allowedTool(r, tool) {
						r.tool = tool
						result = append(result, r)
					}
				}
			}
		}
		return result
	}

	distance := func(from, to geology) int {
		if from.tool == to.tool && allowedTool(to, from.tool) {
			return 1
		}

		// valid change tool
		for tool := torch; tool <= neither; tool++ {
			if allowedTool(from, tool) && allowedTool(to, tool) {
				return 7 + 1
			}
		}

		return math.MaxInt32
	}

	queue := make(geos, 0)
	for _, v := range c.area {

		v.distance = math.MaxInt32

		switch {
		case v.coord.equals(&c.lt):
			v.tool = torch
			v.distance = 0
			queue = append(queue, v)
		case v.coord.equals(&c.target):
			v.tool = torch
			queue = append(queue, v)
		default:
			for tool := torch; tool <= neither; tool++ {
				if allowedTool(v, tool) {
					v.tool = tool
					queue = append(queue, v)
				}
			}
		}
	}

	dists := make(map[string]int)
	for _, v := range queue {
		dists[v.hash()] = v.distance
	}

	for len(queue) > 0 {
		sort.Sort(queue)
		cur := queue[0]
		queue = queue[1:]

		for _, v := range neighbors(cur) {
			newDist := dists[cur.hash()] + distance(cur, v)
			if newDist < dists[v.hash()] {
				v.distance = newDist
				dists[v.hash()] = newDist
				queue = append(queue, v)
			}
		}
	}

	return dists[geology{coord: c.target, tool: torch}.hash()] + 2
}

func (c *cave) String() string {
	var sb strings.Builder

	for y := c.lt.y; y <= c.rb.y; y++ {
		for x := c.lt.x; x <= c.rb.x; x++ {
			if c.lt.equals(&point{x, y}) {
				_, _ = fmt.Fprint(&sb, string(mouth))
				continue
			}
			if c.target.equals(&point{x, y}) {
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
	c.geodesy(depth, *target)
	//fmt.Printf("%v\n", c)
	fmt.Printf("Success! Target minutes is: %d (elapsed time: %s)\n", c.minWalk(), time.Since(now))
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
