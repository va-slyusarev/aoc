package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
)

var input = flag.String("input", "../input.txt", "Input data file path.")
var maxIteration = flag.Int("it", 150, "Max iteration value.")

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}

	sum := 0
	combinations := map[int]int{sum: 1}
	iteration := 1

	for {

		fmt.Printf("Reading input data file interation count: %d\n", iteration)

		scanner := bufio.NewScanner(bytes.NewBuffer(f))
		for scanner.Scan() {
			value := scanner.Text()
			intValue, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("broken data %q from input data file: %v\n", value, err)
				return
			}
			sum += intValue
			if _, ok := combinations[sum]; ok {
				fmt.Printf("Success! Target value is %d\n", sum)
				return
			}
			combinations[sum] = 1
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("error scan input data file: %v\n", err)
			return
		}

		iteration++

		if iteration > *maxIteration {
			fmt.Printf("Failure... Reached iteration limit: %d\n", *maxIteration)
			return
		}
	}
}
