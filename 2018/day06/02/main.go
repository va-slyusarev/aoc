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
var max = flag.Int("max", 10000, "Max distance.")

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

	targetSize := findTargetSize(leftTop, rightBottom, coordinates)

	fmt.Printf("Target size is: %d\n", targetSize)
}

func findTargetSize(leftTop, rightBottom point, coordinates []point) int {

	area := make(map[point]int)

	for x := leftTop.x; x <= rightBottom.x; x++ {
		for y := leftTop.y; y <= rightBottom.y; y++ {

			p := point{x, y}
			for _, c := range coordinates {
				area[p] += p.manhattanDist(c)
			}
		}
	}

	size := 0
	for _, s := range area {
		if s < *max {
			size++
		}
	}

	return size
}
