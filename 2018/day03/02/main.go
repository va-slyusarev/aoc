package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type claim struct {
	id            string
	top           point
	width, height int
}

type point struct {
	x, y int
}

func (c claim) isIntersect(value claim) bool {

	return !((c.top.x >= value.top.x+value.width) ||
		(value.top.x >= c.top.x+c.width) ||
		(c.top.y >= value.top.y+value.height) ||
		(value.top.y) >= c.top.y+c.height)
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var claims []claim

	for scanner.Scan() {
		value := scanner.Text()
		it := claim{}
		if _, err := fmt.Sscanf(value, "#%s @ %d,%d: %dx%d", &it.id, &it.top.x, &it.top.y, &it.width, &it.height); err != nil {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}
		claims = append(claims, it)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	for _, c1 := range claims {
		target := true
		for _, c2 := range claims {

			if c1.id != c2.id && c1.isIntersect(c2) {
				target = false
				break
			}
		}
		if target {
			fmt.Printf("Success! Target id: %s\n", c1.id)
			return
		}
	}

	fmt.Printf("Failure...No target id found.\n")
}
