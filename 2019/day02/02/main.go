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
var number = flag.Int("n", 19690720, "The program to produce the output value")

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

	for noun := 0; noun < 100; noun += 1 {
		for verb := 0; verb < 100; verb += 1 {
			operNums := make([]int, len(nums))
			copy(operNums, nums)

			// under the terms Two Part
			if len(operNums) > 2 {
				operNums[1] = noun
				operNums[2] = verb
			}

			if process(operNums) == *number {
				fmt.Printf("Success! Target number is: %d\n", 100*noun+verb)
				return
			}
		}
	}
	fmt.Printf("Failure...Number not found.\n")
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
