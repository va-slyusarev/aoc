package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type point struct {
	x, y   int
	dx, dy int
}

func (p *point) timeMachine(sec int) {
	p.x += sec * p.dx
	p.y += sec * p.dy
}

func (p point) String() string {
	return fmt.Sprintf("(%d : %d) 🚚 (%d : %d)", p.x, p.y, p.dx, p.dy)
}

type field struct {
	points               []*point
	leftTop, rightBottom point
}

func (f *field) timeMachine(sec int) {
	for i := range f.points {
		f.points[i].timeMachine(sec)
	}
}

func (f *field) frame() {
	leftTop, rightBottom := point{x: math.MaxInt32, y: math.MaxInt32}, point{x: math.MinInt32, y: math.MinInt32}

	for _, p := range f.points {
		if p.x < leftTop.x {
			leftTop.x = p.x
		}
		if p.y < leftTop.y {
			leftTop.y = p.y
		}

		if p.x > rightBottom.x {
			rightBottom.x = p.x
		}
		if p.y > rightBottom.y {
			rightBottom.y = p.y
		}
	}

	f.leftTop = leftTop
	f.rightBottom = rightBottom
}

func (f *field) square(sec int) int {
	f.timeMachine(sec)
	f.frame()
	return (f.rightBottom.x - f.leftTop.x) * (f.rightBottom.y - f.leftTop.y)
}

func (f field) String() string {
	var sky strings.Builder
	_, _ = fmt.Fprintf(&sky, "leftTop: %s, rightBottom: %s\n", f.leftTop, f.rightBottom)

	m := make(map[point]struct{})
	for _, p := range f.points {
		n := point{x: p.x, y: p.y}
		m[n] = struct{}{}
	}

	for y := f.leftTop.y; y <= f.rightBottom.y; y++ {
		for x := f.leftTop.x; x <= f.rightBottom.x; x++ {
			p := point{x: x, y: y}
			if _, ok := m[p]; ok {
				_, _ = fmt.Fprint(&sky, " #")
				continue
			}
			_, _ = fmt.Fprint(&sky, " .")
		}
		_, _ = fmt.Fprintln(&sky)
	}

	return sky.String()
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var re = regexp.MustCompile("position=<\\s?(-?\\d+),\\s+(-?\\d+)>\\svelocity=<\\s?(-?\\d+),\\s+(-?\\d+)>") // position=< 50708, -40198> velocity=<-5,  4>
	sky := field{points: []*point{}}

	for scanner.Scan() {
		value := scanner.Text()
		values := re.FindStringSubmatch(value)

		if len(values) != 5 {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}

		var x, y, dx, dy int
		x, _ = strconv.Atoi(values[1])
		y, _ = strconv.Atoi(values[2])
		dx, _ = strconv.Atoi(values[3])
		dy, _ = strconv.Atoi(values[4])

		sky.points = append(sky.points, &point{x: x, y: y, dx: dx, dy: dy})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	fmt.Printf("Target seconds is: %d\n", lookingSkySec(sky))
}

func lookingSkySec(sky field) int {
	curSquare := sky.square(0)
	sec := 0
	// looking for a minimum area
	for {
		sec++
		s := sky.square(1)
		// fmt.Printf("%d ⏩ %d\n", curSquare, s)
		if s > curSquare {
			break
		}
		curSquare = s
	}

	// minimum - is target
	_ = sky.square(-1)

	return sec - 1
}
