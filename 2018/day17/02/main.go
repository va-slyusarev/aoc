package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

const (
	clay   = '#'
	sand   = '.'
	spring = '+'
	water  = '~'
	leak   = '|'
)

type point struct {
	x, y int
}

func (p *point) dx(dx int) point { return point{p.x + dx, p.y} }

func (p *point) dy(dy int) point { return point{p.x, p.y + dy} }

func (p *point) equals(value *point) bool { return value != nil && p.x == value.x && p.y == value.y }

func (p *point) takeLess(value point) point {
	if value.x < p.x {
		p.x = value.x
	}
	if value.y < p.y {
		p.y = value.y
	}

	return *p
}

func (p *point) takeMore(value point) point {
	if value.x > p.x {
		p.x = value.x
	}
	if value.y > p.y {
		p.y = value.y
	}

	return *p
}

func (p *point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

type ground struct {
	rocks   map[point]rune
	streams map[point]rune

	lt, rb point
}

func (g *ground) String() string {
	var sb strings.Builder

	for y := g.lt.y; y <= g.rb.y; y++ {
		for x := g.lt.x; x <= g.rb.x; x++ {
			cur := point{x, y}
			// clay
			if r, ok := g.rocks[cur]; ok {
				_, _ = fmt.Fprint(&sb, string(r))
				continue
			}
			// leak
			if r, ok := g.streams[cur]; ok {
				_, _ = fmt.Fprint(&sb, string(r))
				continue
			}
			// sand
			_, _ = fmt.Fprint(&sb, string(sand))
		}
		_, _ = fmt.Fprintln(&sb)
	}
	return sb.String()
}

func (g *ground) openValve() {
	leaks := []point{{500, 0}}

	addUnique := func(storage []point, values []point) []point {
		for _, value := range values {
			find := false
			for _, store := range storage {
				if value.equals(&store) {
					find = true
					break
				}
			}
			if !find {
				storage = append(storage, value)
			}
		}
		return storage
	}

	for len(leaks) > 0 {

		cur := leaks[len(leaks)-1]
		leaks = leaks[:len(leaks)-1]

		for cur.y >= g.lt.y && cur.y <= g.rb.y {
			if g.streams[cur] == spring {
				cur = cur.dy(1)
				continue
			}
			if _, ok := g.rocks[cur]; ok {
				leaks = addUnique(leaks, g.horizontalLeak(cur.dy(-1)))
				break
			}
			if g.streams[cur] == leak && g.streams[cur.dy(1)] == water {
				leaks = addUnique(leaks, g.horizontalLeak(cur))
				break
			}

			if g.streams[cur] == water {
				leaks = addUnique(leaks, []point{cur.dy(-1)})
				break
			}

			g.streams[cur] = leak
			cur = cur.dy(1)
		}
	}
}

func (g *ground) horizontalLeak(start point) []point {

	if g.streams[start] == water {
		return []point{start.dy(-1)}
	}

	tmp := make([]point, 0)
	var lClay, rClay *point

	// left
	for cur := start; cur.x >= g.lt.x; cur = cur.dx(-1) {
		if _, ok := g.rocks[cur]; ok {
			lClay = &cur
			break
		}
		if g.rocks[cur.dy(1)] != clay && g.streams[cur.dy(1)] != water {
			g.streams[cur] = leak
			tmp = append(tmp, cur)
			break
		}
		g.streams[cur] = leak
	}

	// right
	for cur := start; cur.x <= g.rb.x; cur = cur.dx(1) {
		if _, ok := g.rocks[cur]; ok {
			rClay = &cur
			break
		}
		if g.rocks[cur.dy(1)] != clay && g.streams[cur.dy(1)] != water {
			g.streams[cur] = leak
			tmp = append(tmp, cur)
			break
		}
		g.streams[cur] = leak
	}

	if lClay != nil && rClay != nil {
		for cur := lClay.dx(1); cur.x < rClay.x; cur = cur.dx(1) {
			g.streams[cur] = water
		}
		return []point{start.dy(-1)}
	}

	return tmp
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

	gr, err := parse(lines)
	if err != nil {
		fmt.Printf("error scan input data: %v", err)
		return
	}

	now := time.Now()

	gr.openValve()
	//fmt.Printf("%v\n", gr)

	countWater := func(g *ground) int {
		result := 0
		for _, v := range g.streams {
			if v == water {
				result++
			}
		}
		return result
	}

	fmt.Printf("Success! Target number is: %d (elapsed time: %s)\n", countWater(gr), time.Since(now))
}

func parse(lines []string) (*ground, error) {
	result := &ground{
		rocks:   make(map[point]rune),
		streams: map[point]rune{point{500, 0}: spring},
		lt:      point{500, 0}, rb: point{500, 0},
	}
	re := regexp.MustCompile(`([x,y])=(\d+),\s[x,y]=(\d+)[.]{2}(\d+)`)
	for _, line := range lines {
		values := re.FindStringSubmatch(line)
		if len(values) != 5 {
			return nil, fmt.Errorf("broken data: %q", line)
		}
		var a, b1, b2 int

		if _, err := fmt.Sscanf(values[2], "%d", &a); err != nil {
			return nil, fmt.Errorf("broken data: unknown a: %q", line)
		}
		if _, err := fmt.Sscanf(values[3], "%d", &b1); err != nil {
			return nil, fmt.Errorf("broken data: unknown b1: %q", line)
		}
		if _, err := fmt.Sscanf(values[4], "%d", &b2); err != nil {
			return nil, fmt.Errorf("broken data: unknown b2: %q", line)
		}

		switch values[1] {
		case "x":
			for i := b1; i <= b2; i++ {
				cur := point{x: a, y: i}
				result.rocks[cur] = clay
				_ = result.lt.takeLess(cur)
				_ = result.rb.takeMore(cur)
			}
		case "y":
			for i := b1; i <= b2; i++ {
				cur := point{x: i, y: a}
				result.rocks[cur] = clay
				_ = result.lt.takeLess(cur)
				_ = result.rb.takeMore(cur)
			}
		default:
			return nil, fmt.Errorf("broken data: unknown coord: %q", line)
		}
	}

	// to expand on x
	result.lt = result.lt.dx(-1)
	result.rb = result.rb.dx(1)

	return result, nil
}
