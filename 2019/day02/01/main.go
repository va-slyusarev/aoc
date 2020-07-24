package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
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

	var nums []int

	for scanner.Scan() {
		strNums := strings.FieldsFunc(scanner.Text(), func(r rune) bool {
			return !unicode.IsDigit(r)
		})
		for _, str := range strNums {
			value, err := strconv.Atoi(str)
			if err != nil {
				fmt.Printf("error scan input data file: %v\n", err)
				return
			}
			nums = append(nums, value)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	// under the terms One Part
	if nums != nil && len(nums) > 2 {
		nums[1] = 12
		nums[2] = 2
	}

	fmt.Printf("Success! Target number is: %d\n", process(nums))
}

func process(nums []int) int {
	idx := 0

loop:
	for {
		op := nums[idx]

		switch op {
		case 1, 2:
			oneIdx := nums[idx+1]
			twoIdx := nums[idx+2]
			resultIdx := nums[idx+3]

			if op == 1 {
				nums[resultIdx] = nums[oneIdx] + nums[twoIdx]
			} else if op == 2 {
				nums[resultIdx] = nums[oneIdx] * nums[twoIdx]
			}
			idx += 4
		case 99:
			break loop
		default:
			fmt.Printf("error operation: %v\n", op)
			os.Exit(1)
		}
	}

	return nums[0]
}
