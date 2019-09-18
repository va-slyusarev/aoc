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

	// point and IDs elfs
	wholePiece := make(map[point][]string)

	for _, v := range claims {

		for x := v.top.x; x < v.top.x+v.width; x++ {

			for y := v.top.y; y < v.top.y+v.height; y++ {

				p := point{x, y}

				if _, ok := wholePiece[p]; ok {
					wholePiece[p] = append(wholePiece[p], v.id)
					continue
				}

				wholePiece[p] = []string{v.id}
			}
		}
	}

	targetSquare := 0
	for _, elfs := range wholePiece {
		if len(elfs) > 1 {
			targetSquare++
		}
	}

	fmt.Printf("Success! Target square is: %d\n", targetSquare)
}
