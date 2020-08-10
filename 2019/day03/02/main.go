package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"unicode"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type state struct {
	num      []int
	intersec bool
}

type wires struct {
	pt   map[image.Point]state
	rect image.Rectangle
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines [][]string

	for scanner.Scan() {
		lines = append(lines, strings.FieldsFunc(scanner.Text(), func(r rune) bool {
			return !unicode.IsDigit(r) && !unicode.IsLetter(r)
		}))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	w := paveWires(rect(lines), lines)

	fmt.Printf("Success! Target number is: %d\n", manhattanDist(w))
}

func rect(lines [][]string) *wires {
	min := image.Pt(0, 0)
	max := image.Pt(0, 0)

	for _, line := range lines {

		current := image.Pt(0, 0)

		for _, direction := range line {
			var dir string
			var value int
			if _, err := fmt.Sscanf(direction, "%1s%d", &dir, &value); err != nil {
				log.Fatalf("Error parsing direction: %v\n", err)
			}
			switch dir {
			case "U":
				current = current.Add(image.Pt(0, value))
				if max.Y < current.Y {
					max.Y = current.Y
				}
			case "R":
				current = current.Add(image.Pt(value, 0))
				if max.X < current.X {
					max.X = current.X
				}
			case "D":
				current = current.Add(image.Pt(0, -value))
				if min.Y > current.Y {
					min.Y = current.Y
				}
			case "L":
				current = current.Add(image.Pt(-value, 0))
				if min.X > current.X {
					min.X = current.X
				}
			default:
				log.Fatalf("Unknown direction: %v\n", dir)
			}
		}
	}

	// extended area
	min = min.Sub(image.Pt(1, 1))
	max = max.Add(image.Pt(2, 2))

	return &wires{
		pt:   make(map[image.Point]state),
		rect: image.Rectangle{Min: min, Max: max},
	}
}

func paveWires(w *wires, lines [][]string) *wires {

	// update state current point
	updState := func(w *wires, cur image.Point, step int, idx int) {
		if _, ok := w.pt[cur]; !ok {
			st := state{num: []int{0, 0}, intersec: false}
			st.num[idx] = step
			w.pt[cur] = st
		} else {
			st := w.pt[cur]
			st.num[idx] = step

			if st.num[0] > 0 && st.num[1] > 0 {
				st.intersec = true
			}
			w.pt[cur] = st
		}
	}

	for idx, line := range lines {

		current := image.Pt(0, 0)
		step := 0

		for _, direction := range line {
			var dir string
			var value int
			if _, err := fmt.Sscanf(direction, "%1s%d", &dir, &value); err != nil {
				log.Fatalf("Error parsing direction: %v\n", err)
			}
			switch dir {
			case "U":
				for y := 1; y <= value; y++ {
					current = current.Add(image.Pt(0, 1))
					step++
					updState(w, current, step, idx)
				}

			case "R":
				for x := 1; x <= value; x++ {
					current = current.Add(image.Pt(1, 0))
					step++
					updState(w, current, step, idx)
				}
			case "D":
				for y := 1; y <= value; y++ {
					current = current.Add(image.Pt(0, -1))
					step++
					updState(w, current, step, idx)
				}
			case "L":
				for x := 1; x <= value; x++ {
					current = current.Add(image.Pt(-1, 0))
					step++
					updState(w, current, step, idx)
				}
			default:
				log.Fatalf("Unknown direction: %v\n", dir)
			}
		}
	}
	return w
}

func manhattanDist(w *wires) int {

	result := math.MaxInt32
	for _, v := range w.pt {

		// intersection wire's
		if v.intersec {
			dist := v.num[0] + v.num[1]
			if dist < result {
				result = dist
			}
		}
	}
	return result
}
