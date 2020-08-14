package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type instruction [5]int
type params [3]int

// A B C D E
// 0 1 2 3 4
func newInstruction(value int) instruction {

	getN := func(value int, idx int) int {
		s := fmt.Sprintf("%05s", strconv.Itoa(value))
		if len(s) < idx {
			log.Fatalf("error parsing instruction: number is %d, digit index is %d\n", value, idx)
		}
		i, err := strconv.Atoi(s[idx-1 : idx])
		if err != nil {
			log.Fatalf("error parsing digit: %v\n", err)
		}
		return i
	}

	return instruction{getN(value, 1), getN(value, 2), getN(value, 3), getN(value, 4), getN(value, 5)}
}

// A B C DE
// 0 1 2 34
//       operation
func (i instruction) operation() int {
	return 10*i[3] + i[4]
}

// A B C DE
// 3 2 1
func (i instruction) isPosMode(pos int) bool {
	if pos < 1 && pos > 3 {
		log.Fatalf("wrong position (available is [1..3]): %d\n", pos)
	}
	return i[3-pos] == 0
}

func (i instruction) params(nums []int, currentIdx int) params {
	oneParam, twoParam, threeParam := 0, 0, 0

	if currentIdx+1 < len(nums) {
		oneParam = nums[currentIdx+1]
	}
	if currentIdx+2 < len(nums) {
		twoParam = nums[currentIdx+2]
	}
	if currentIdx+3 < len(nums) {
		threeParam = nums[currentIdx+3]
	}

	// one & two params, because three param only output
	if i.operation() != 3 && i.operation() != 4 && i.operation() != 99 {
		if i.isPosMode(1) {
			oneParam = nums[oneParam]
		}
		if i.isPosMode(2) {
			twoParam = nums[twoParam]
		}
	}
	return params{oneParam, twoParam, threeParam}
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		log.Fatalf("error read input data file: %v\n", err)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var nums []int

	for scanner.Scan() {
		strNums := strings.FieldsFunc(scanner.Text(), func(r rune) bool {
			return !unicode.IsNumber(r) && r != '-'
		})
		for _, str := range strNums {
			value, err := strconv.Atoi(str)
			if err != nil {
				log.Fatalf("error scan input data file: %v\n", err)
			}
			nums = append(nums, value)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error scan input data file: %v\n", err)
	}

	log.Printf("Success! Target number is: %d\n", process(nums))
}

func process(nums []int) int {
	result := 0
	idx := 0
loop:
	for {
		inst := newInstruction(nums[idx])
		op := inst.operation()
		pms := inst.params(nums, idx)

		switch op {
		case 1, 2: // add, multiply
			if op == 1 {
				nums[pms[2]] = pms[0] + pms[1]
			} else if op == 2 {
				nums[pms[2]] = pms[0] * pms[1]
			}
			idx += 4

		case 3: // input
			nums[pms[0]] = 5 // ship's thermal radiator controller
			idx += 2

		case 4: // output
			result = nums[pms[0]]
			idx += 2

		case 5: // jump-if-true
			if pms[0] != 0 {
				idx = pms[1]
			} else {
				idx += 3
			}

		case 6: // jump-if-false
			if pms[0] == 0 {
				idx = pms[1]
			} else {
				idx += 3
			}

		case 7: // less
			if pms[0] < pms[1] {
				nums[pms[2]] = 1
			} else {
				nums[pms[2]] = 0
			}
			idx += 4

		case 8: // equals
			if pms[0] == pms[1] {
				nums[pms[2]] = 1
			} else {
				nums[pms[2]] = 0
			}
			idx += 4

		case 99: // exit
			break loop
		default:
			log.Fatalf("wrong operation: %v, (index = %d)\n", op, idx)
		}
	}
	return result
}
