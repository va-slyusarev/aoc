package main

import (
	"flag"
	"fmt"
	"math"
	"sync"
	"time"
)

var serial = flag.Int("s", 3463, "Serial number (input data).")
var mXY = flag.Int("x", 300, "Max X & Y coordinate.")

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

type chFuel struct {
	size int
	p    point
	fuel int
}

func workers(c cells) <-chan chFuel {
	ch := make(chan chFuel, *mXY)
	go func() {
		var wg sync.WaitGroup
		for s := 1; s <= *mXY; s++ {
			wg.Add(1)
			go func(size int) {
				defer wg.Done()
				it, max := c.maxPowerSquare(size)
				ch <- chFuel{size: size, p: it, fuel: max}
			}(s)
		}
		go func() {
			wg.Wait()
			close(ch)
		}()
	}()

	return ch
}

func findMax(ch <-chan chFuel) chFuel {
	result := chFuel{}
	for c := range ch {
		fmt.Printf("it ⏩ %s,%d📦 (⛽  %d)\n", c.p, c.size, c.fuel)
		if c.fuel > result.fuel {
			result.size = c.size
			result.p = c.p
			result.fuel = c.fuel
		}
	}
	return result
}

func main() {
	flag.Parse()

	c := make(cells)
	c.init(*mXY, *serial)

	t := time.Now()
	r := findMax(workers(c))

	fmt.Printf("Target coordinate is: %s,%d📦 (⛽  %d), time: %s\n", r.p, r.size, r.fuel, time.Since(t).Round(time.Millisecond))
}
