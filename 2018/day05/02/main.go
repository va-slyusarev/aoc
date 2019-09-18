package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
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
	minUnits := len(polymer)

	for i := 'A'; i <= 'Z'; i++ {

		correct := strings.ReplaceAll(polymer, string(i), "")
		correct = strings.ReplaceAll(correct, string(i+32), "")

		react := true

		for ; react; {
			correct, react = reaction(correct)
		}

		if len(correct) < minUnits {
			minUnits = len(correct)
		}
	}

	fmt.Printf("Target polymer units is: %d\n", minUnits)
}

func reaction(polymer string) (string, bool) {
	for i := 0; i < len(polymer)-1; i++ {
		if polymer[i]-polymer[i+1] == 32 || polymer[i+1]-polymer[i] == 32 {
			return polymer[:i] + polymer[i+2:], true
		}
	}
	return polymer, false
}
