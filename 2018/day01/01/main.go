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

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	sum := 0

	for scanner.Scan() {
		value := scanner.Text()
		intValue, err := strconv.Atoi(value)
		if err != nil {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}
		sum += intValue
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	fmt.Printf("Success! Target value is %d\n", sum)
}
