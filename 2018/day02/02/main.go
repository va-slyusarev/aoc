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

	var letters1 []string
	var letters2 []string

	for scanner.Scan() {
		letters1 = append(letters1, scanner.Text())
		letters2 = append(letters2, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	for _, value1 := range letters1 {
		for _, value2 := range letters2 {
			if targetValue, ok := targetPair(value1, value2); ok {
				fmt.Printf("Success! Common letters are: %q\n", targetValue)
				return
			}
		}
	}

	fmt.Printf("Failure...No common letters found.\n")
}

func targetPair(letters1 string, letters2 string) (string, bool) {
	v1 := strings.Split(letters1, "")
	v2 := strings.Split(letters2, "")

	if len(v1) != len(v2) {
		return "", false
	}

	diff := 0
	lastIdx := 0

	for i, c := range v1 {
		if v2[i] != c {
			lastIdx = i
			diff += 1
		}
	}

	if diff == 1 {
		return strings.Join(append(v1[:lastIdx], v1[lastIdx+1:]...), ""), true
	}
	return "", false
}
