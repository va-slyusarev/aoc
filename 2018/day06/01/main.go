package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type point struct {
	x, y int
}

func (p *point) String() string {
	return fmt.Sprintf("[%d:%d]", p.x, p.y)
}

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

func (p *point) manhattanDist(value point) int {
	return int(math.Abs(float64(p.x-value.x)) + math.Abs(float64(p.y-value.y)))
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var coordinates []point
	leftTop, rightBottom := point{math.MaxUint32, math.MaxUint32}, point{0, 0}

	for scanner.Scan() {
		value := scanner.Text()
		it := point{}
		if _, err := fmt.Sscanf(value, "%d, %d", &it.x, &it.y); err != nil {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}

		_ = leftTop.takeLess(it)
		_ = rightBottom.takeMore(it)

		coordinates = append(coordinates, it)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	targetArea := findMaxArea(leftTop, rightBottom, coordinates)

	fmt.Printf("Target area is: %d\n", targetArea)
}

func infinityCoordinates(coordinates []point) []point {

	var infinity []point
	for _, p1 := range coordinates {

		leftTop, rightTop, leftBottom, rightBottom := false, false, false, false

		for _, p2 := range coordinates {
			if p1.x > p2.x && p1.y > p2.y {
				leftTop = true
			}
			if p1.x < p2.x && p1.y > p2.y {
				rightTop = true
			}
			if p1.x > p2.x && p1.y < p2.y {
				leftBottom = true
			}
			if p1.x < p2.x && p1.y < p2.y {
				rightBottom = true
			}
		}

		if !(leftTop && rightTop && leftBottom && rightBottom) {
			infinity = append(infinity, p1)
		}
	}

	return infinity
}

func findMaxArea(leftTop, rightBottom point, coordinates []point) int {

	type who struct {
		coordinate point
		distance   int
		nobody     bool
	}

	area := make(map[point]who)

	for x := leftTop.x; x <= rightBottom.x; x++ {
		for y := leftTop.y; y <= rightBottom.y; y++ {

			p := point{x, y}
			area[p] = who{distance: math.MaxUint32}

			for _, c := range coordinates {

				currentWho := area[p]
				manhattan := c.manhattanDist(p)
				if manhattan < currentWho.distance {
					currentWho.coordinate = c
					currentWho.distance = manhattan
					currentWho.nobody = false
					area[p] = currentWho
					continue
				}

				if manhattan == currentWho.distance {
					currentWho.nobody = true
					area[p] = currentWho
					continue
				}
			}
		}
	}

	infinity := infinityCoordinates(coordinates)

	areas := make(map[point]int)
	for _, v := range area {
		if v.nobody {
			continue
		}
		// skip infinity
		exist := false
		for _, inf := range infinity {
			if inf == v.coordinate {
				exist = true
			}
		}
		if exist {
			continue
		}
		areas[v.coordinate] += 1
	}

	max := 0
	for _, v := range areas {
		if v > max {
			max = v
		}
	}

	return max
}
