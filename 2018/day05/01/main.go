package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
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

	var polymers []string

	for scanner.Scan() {
		polymers = append(polymers, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	polymer := polymers[0]

	react := true

	for ; react; {
		polymer, react = reaction(polymer)
	}

	fmt.Printf("Target polymer units is: %d\n", len(polymer))
}

func reaction(polymer string) (string, bool) {
	for i := 0; i < len(polymer)-1; i++ {
		if polymer[i]-polymer[i+1] == 32 || polymer[i+1]-polymer[i] == 32 {
			return polymer[:i] + polymer[i+2:], true
		}
	}
	return polymer, false
}
