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

type xyz struct {
	x, y, z int
}

func (p xyz) manhattanDist(t xyz) int {
	dx := math.Abs(float64(t.x - p.x))
	dy := math.Abs(float64(t.y - p.y))
	dz := math.Abs(float64(t.z - p.z))
	return int(dx + dy + dz)
}

func (p xyz) String() string {
	return fmt.Sprintf("[%d,%d,%d]", p.x, p.y, p.z)
}

type nanobot struct {
	point xyz
	r     int
}

func (b nanobot) overlap(t nanobot) bool {
	return b.point.manhattanDist(t.point) <= b.r
}

func (b nanobot) String() string {
	return fmt.Sprintf("%v(r=%d)", b.point, b.r)
}

func main() {
	bots, err := parse()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	overlaps := make(map[nanobot][]nanobot)
	for _, bot := range bots {
		overlaps[bot] = make([]nanobot, 0)
		for _, target := range bots {
			if bot.overlap(target) {
				overlaps[bot] = append(overlaps[bot], target)
			}
		}
	}

	// Find the nanobot with the largest signal radius. How many nanobots are in range of its signals?
	var maxBot nanobot
	maxRadius := 0
	for bot := range overlaps {
		if bot.r > maxRadius {
			maxBot = bot
			maxRadius = bot.r
		}
	}

	fmt.Printf("Success! Target number is: %d\n", len(overlaps[maxBot]))
}

func parse() ([]nanobot, error) {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {

		return nil, fmt.Errorf("error read input data file: %v", err)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines []string

	for scanner.Scan() {
		value := scanner.Text()
		lines = append(lines, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scan input data file: %v", err)
	}

	result := make([]nanobot, 0)
	for _, line := range lines {
		newBot := nanobot{point: xyz{}}

		if _, err := fmt.Sscanf(line, "pos=<%d,%d,%d>, r=%d", &newBot.point.x, &newBot.point.y, &newBot.point.z, &newBot.r); err != nil {
			return nil, fmt.Errorf("error parse data file: %v", err)
		}
		result = append(result, newBot)
	}

	return result, nil
}
