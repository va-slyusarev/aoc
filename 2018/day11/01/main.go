package main

import (
	"flag"
	"fmt"
	"math"
)

var serial = flag.Int("s", 3463, "Serial number (input data).")
var mXY = flag.Int("x", 300, "Max X & Y coordinate.")
var square = flag.Int("z", 3, "Size square.")

type point struct {
	x, y int
}

type cells map[point]int

func (p point) String() string {
	return fmt.Sprintf("%d,%d", p.x, p.y)
}

func (p point) power(serial int) int {
	id := p.x + 10
	pwr := id * p.y
	pwr += serial
	pwr *= id
	pwr %= 1000
	pwr /= 100
	pwr -= 5
	return pwr
}

func (c cells) init(mXY int, serial int) {
	for y := 1; y <= mXY; y++ {
		for x := 1; x <= mXY; x++ {
			p := point{x, y}
			c[p] = p.power(serial)
		}
	}
}

func (c cells) powerSquare(p point, size int) int {
	power := 0
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			power += c[point{p.x + x, p.y + y}]
		}
	}
	return power
}

func (c cells) maxPowerSquare(size int) (point, int) {
	var p point
	max := math.MinInt32

	for y := 1; y <= *mXY-size+1; y++ {
		for x := 1; x <= *mXY-size+1; x++ {
			it := point{x, y}
			if sq := c.powerSquare(it, size); sq > max {
				p = it
				max = sq
			}
		}
	}

	return p, max
}

func main() {
	flag.Parse()

	c := make(cells)
	c.init(*mXY, *serial)

	p, fuel := c.maxPowerSquare(*square)

	fmt.Printf("Target coordinate is: %s (⛽  %d)\n", p, fuel)
}
