package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		log.Fatalf("error read input data file: %v\n", err)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines []string

	for scanner.Scan() {
		value := scanner.Text()
		lines = append(lines, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error scan input data file: %v\n", err)
	}

	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
}
