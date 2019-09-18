package main

// TODO: not resolved

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
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

func (p xyz) offset(point xyz) xyz {
	return xyz{p.x + point.x, p.y + point.y, p.z + point.z}
}

func (p xyz) String() string {
	return fmt.Sprintf("[%d,%d,%d]", p.x, p.y, p.z)
}

type nanobot struct {
	point xyz
	r     int

	count  int
	dist0  int
	inside bool
}

func (b nanobot) overlap(t nanobot) bool {
	return b.point.manhattanDist(t.point) <= b.r+t.r
}

func (b nanobot) String() string {
	return fmt.Sprintf("%v(r=%d, c=%d)", b.point, b.r, b.count)
}

func (b nanobot) hash() string {
	return fmt.Sprintf("%vr%d", b.point, b.r)
}

func main() {
	bots, err := parse()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	countOverlap := func(bots []nanobot, b nanobot) int {
		result := 0
		for _, bot := range bots {
			if bot.overlap(b) && bot.hash() != b.hash() {
				result++
			}
		}
		return result
	}

	/*
		iAmInside := func(bots []nanobot, b nanobot) bool {
			for _, bot := range bots {

				dist := bot.point.manhattanDist(b.point)
				r := bot.r-b.r
				fmt.Printf("bot: %v, dist: %d, r: %d\n", bot, dist, r)

				if dist <= r && bot.String() != b.String() {
					return true
				}
			}
			return false
		}
	*/

	candidates := func(bots []nanobot) []nanobot {

		for i, bot := range bots {
			bot.count = countOverlap(bots, bot)
			bot.dist0 = bot.point.manhattanDist(xyz{})
			bots[i] = bot
		}

		sort.Slice(bots, func(i, j int) bool {
			return bots[i].count > bots[j].count
		})

		fmt.Printf("sorted by overlaps: %v\n", bots)

		max := bots[0].count
		for i, bot := range bots {
			if bot.count < max {
				bots = bots[:i]
				break
			}
		}

		fmt.Printf("get max overlaps: %v\n", bots)

		sort.Slice(bots, func(i, j int) bool {
			return bots[i].dist0 < bots[j].dist0
		})

		fmt.Printf("sorted by dist0: %v\n", bots)

		fmt.Printf("cand: %v\n", bots[0])

		return bots[:1]
	}

	minDist := math.MaxInt32
	maxOver := 0
	for _, bot := range candidates(bots) {
		fmt.Printf("bot: %v\n", bot)
		for x := 0; x < bot.r+1; x++ {
			for y := 0; y < bot.r+1; y++ {
				for z := 0; z < bot.r+1; z++ {
					point := bot.point.offset(xyz{x, y, z})
					if dist, over := point.manhattanDist(xyz{}), countOverlap(bots, nanobot{point: point}); over >= maxOver {

						fmt.Printf("point: %v dist: %d, overlaps: %d\n", point, dist, over)

						if over > maxOver {
							minDist = dist
						}
						if over == maxOver && dist < minDist {
							minDist = dist
						}
						maxOver = over
					}
				}
			}
		}
	}

	fmt.Printf("Success! Target number is: %d (overlaps: %d)\n", minDist, maxOver)
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
