package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines []string

	for scanner.Scan() {
		value := scanner.Text()
		lines = append(lines, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	var total int64
	for _, line := range lines {
		mass, _ := strconv.Atoi(line)
		total += fullCalc(mass)
	}

	fmt.Printf("Success! Target number is: %d\n", total)
}

func calcFuel(mass int) int64 {
	return int64(math.Floor(float64(mass)/3) - 2)
}

func fullCalc(mass int) int64 {
	var result int64
	for curr := calcFuel(mass); curr > 0; curr = calcFuel(int(curr)) {
		result += curr
	}
	return result
}
