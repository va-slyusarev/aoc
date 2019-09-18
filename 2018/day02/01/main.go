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

	twoLetters := 0
	threeLetters := 0

	for scanner.Scan() {
		value := scanner.Text()
		twoLetters += containsDuplicate(value, 2)
		threeLetters += containsDuplicate(value, 3)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	fmt.Printf("Success! Target checksum is %d\n", twoLetters*threeLetters)
}

// 0 - no duplicate
// 1 - contains duplicate (no matter how many times)
func containsDuplicate(value string, quantity int) int {

	letters := map[string]int{}
	for _, letter := range strings.Split(value, "") {
		if _, ok := letters[letter]; ok {
			letters[letter] += 1
			continue
		}
		letters[letter] = 1
	}

	result := 0

	for _, count := range letters {
		if count == quantity {
			result += 1
			break
		}
	}
	return result
}
