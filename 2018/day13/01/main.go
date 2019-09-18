package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"time"
)

var input = flag.String("input", "../input.txt", "Input data file path.")
var visual = flag.Bool("v", false, "Visualization.")
var speed = flag.Int("ms", 1000, "Speed visualization (ms).")

const (
	up    = '^'
	down  = 'v'
	left  = '<'
	right = '>'

	vertical     = '|'
	horizontal   = '-'
	lcurve       = '/'
	rcurve       = '\\'
	intersection = '+'

	nobody = ' '
	crash  = 'X'
)

// current surface: current direction: move function
var movements = map[rune]map[rune]move{
	vertical:     {up: moveUp, down: moveDown},
	horizontal:   {left: moveLeft, right: moveRight},
	lcurve:       {up: moveRight, down: moveLeft, left: moveDown, right: moveUp},
	rcurve:       {up: moveLeft, down: moveRight, left: moveUp, right: moveDown},
	intersection: {up: moveIf, down: moveIf, left: moveIf, right: moveIf},
}

// current direction: last crossroad direction (left, up, right, just an indication of direction): move function
var interMovements = map[rune]map[rune]move{
	up:    {nobody: moveLeft, left: moveUp, up: moveRight, right: moveLeft},
	down:  {nobody: moveRight, left: moveDown, up: moveLeft, right: moveRight},
	left:  {nobody: moveDown, left: moveLeft, up: moveUp, right: moveDown},
	right: {nobody: moveUp, left: moveRight, up: moveDown, right: moveUp},
}

// last crossroad direction: next crossroad direction
var lastNextDirection = map[rune]rune{
	nobody: left,
	left:   up,
	up:     right,
	right:  left,
}

type move func(c *cart)

func moveUp(c *cart)    { c.p.y--; c.value = up }
func moveDown(c *cart)  { c.p.y++; c.value = down }
func moveLeft(c *cart)  { c.p.x--; c.value = left }
func moveRight(c *cart) { c.p.x++; c.value = right }
func moveIf(c *cart) {
	interMovements[c.value][c.lastCrossroad](c)
	c.lastCrossroad = lastNextDirection[c.lastCrossroad]
}

type point struct {
	x, y int
}

func (p point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p *point) takeMore(value point) {
	if value.x > p.x {
		p.x = value.x
	}
	if value.y > p.y {
		p.y = value.y
	}
}

func (p point) less(value point) bool {
	if p.y != value.y {
		return p.y < value.y
	}
	return p.x < value.x
}

func (p point) equals(value *point) bool {
	return value != nil && p.x == value.x && p.y == value.y
}

type cart struct {
	p point

	value         rune
	lastCrossroad rune
}

func (c cart) String() string {
	return fmt.Sprintf("[%s](%q - %q)", c.p, string(c.value), string(c.lastCrossroad))
}

type carts []*cart

func (c carts) Len() int           { return len(c) }
func (c carts) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c carts) Less(i, j int) bool { return c[i].p.less(c[j].p) }

type mine struct {
	pieces map[point]rune
	max    *point
	crash  *point

	workers carts
}

func (m mine) String() string {
	var sb strings.Builder

	for y := 0; y <= m.max.y; y++ {
		for x := 0; x <= m.max.x; x++ {
			cur := point{x, y}

			if cur.equals(m.crash) {
				_, _ = fmt.Fprint(&sb, string(crash))
				continue
			}

			findWorker := false
			for _, w := range m.workers {
				if w.p.equals(&cur) {
					_, _ = fmt.Fprint(&sb, string(w.value))
					findWorker = true
					break
				}
			}
			if findWorker {
				continue
			}

			if v, ok := m.pieces[cur]; ok {
				_, _ = fmt.Fprint(&sb, string(v))
			}
		}
		_, _ = fmt.Fprintln(&sb)
	}
	return sb.String()
}

func (m *mine) Sort() {
	sort.Sort(m.workers)
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

	m := initMine(lines)
	for {
		if *visual {
			fmt.Printf("\033[H\033[2J%s", m)
			time.Sleep(time.Duration(*speed) * time.Millisecond)
		}
		if p, crash := tick(&m); crash {
			if *visual {
				fmt.Printf("\033[H\033[2J%s", m)
			}
			fmt.Printf("Success! Target point is: %v\n", *p)
			return
		}
	}
}

func initMine(lines []string) mine {
	m := mine{pieces: make(map[point]rune), max: &point{0, 0}, workers: make(carts, 0)}

	for y, line := range lines {
		for x, r := range line {
			p := point{x, y}
			m.max.takeMore(p)

			if _, ok := interMovements[r]; ok {
				m.workers = append(m.workers, &cart{p: p, value: r, lastCrossroad: nobody})

				switch r {
				case up, down:
					m.pieces[p] = vertical
				case left, right:
					m.pieces[p] = horizontal
				}

				continue
			}

			if _, ok := movements[r]; ok {
				m.pieces[p] = r
				continue
			}

			m.pieces[p] = nobody
		}
	}

	return m
}

func tick(m *mine) (*point, bool) {

	busy := make(map[point]struct{})
	for _, w := range m.workers {
		busy[w.p] = struct{}{}
	}
	moves := make(map[point]struct{})

	m.Sort()
	for _, w := range m.workers {
		delete(busy, w.p)
		movements[m.pieces[w.p]][w.value](w)

		if _, cr := busy[w.p]; cr {
			m.crash = &w.p
			return &w.p, true
		}
		if _, cr := moves[w.p]; cr {
			m.crash = &w.p
			return &w.p, true
		}
		moves[w.p] = struct{}{}
	}

	return nil, false
}
