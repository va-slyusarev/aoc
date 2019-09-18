package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type registers [4]int

func (r registers) copy() registers {
	var clone registers
	for i := range r {
		clone[i] = r[i]
	}
	return clone
}

func (r registers) equals(value registers) bool {
	for i := range r {
		if value[i] != r[i] {
			return false
		}
	}
	return true
}

type operation struct {
	code int
	a, b int
	c    int
}

type do = func(reg registers, op operation) registers

var operations = map[int]do{
	100: addr,
	101: addi,
	102: mulr,
	103: muli,
	104: banr,
	105: bani,
	106: borr,
	107: bori,
	108: setr,
	109: seti,
	110: gtir,
	111: gtri,
	112: gtrr,
	113: eqir,
	114: eqri,
	115: eqrr,
}

type test struct {
	before, after registers
	op            operation
}

type tTable []test

func addr(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] + result[op.b]
	return result
}

func addi(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] + op.b
	return result
}

func mulr(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] * result[op.b]
	return result
}

func muli(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] * op.b
	return result
}
func banr(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] & result[op.b]
	return result
}

func bani(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] & op.b
	return result
}
func borr(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] | result[op.b]
	return result
}

func bori(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a] | op.b
	return result
}

func setr(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = result[op.a]
	return result
}

func seti(reg registers, op operation) registers {
	result := reg.copy()
	result[op.c] = op.a
	return result
}

func gtir(reg registers, op operation) registers {
	result := reg.copy()
	if op.a > result[op.b] {
		result[op.c] = 1
		return result
	}
	result[op.c] = 0
	return result
}
func gtri(reg registers, op operation) registers {
	result := reg.copy()
	if result[op.a] > op.b {
		result[op.c] = 1
		return result
	}
	result[op.c] = 0
	return result
}
func gtrr(reg registers, op operation) registers {
	result := reg.copy()
	if result[op.a] > result[op.b] {
		result[op.c] = 1
		return result
	}
	result[op.c] = 0
	return result
}
func eqir(reg registers, op operation) registers {
	result := reg.copy()
	if op.a == result[op.b] {
		result[op.c] = 1
		return result
	}
	result[op.c] = 0
	return result
}
func eqri(reg registers, op operation) registers {
	result := reg.copy()
	if result[op.a] == op.b {
		result[op.c] = 1
		return result
	}
	result[op.c] = 0
	return result
}
func eqrr(reg registers, op operation) registers {
	result := reg.copy()
	if result[op.a] == result[op.b] {
		result[op.c] = 1
		return result
	}
	result[op.c] = 0
	return result
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	var lines []string

	for scanner.Scan() {
		value := scanner.Text()
		lines = append(lines, value)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	tests := parse(lines)

	goal := 0
	for _, tt := range tests {
		cnt := 0
		for _, op := range operations {
			after := op(tt.before, tt.op)
			if tt.after.equals(after) {
				cnt++
			}
		}

		if cnt > 2 {
			goal++
		}
	}

	fmt.Printf("Success! Target number is: %d\n", goal)
}

func parse(lines []string) tTable {
	var t test
	before, after := false, false
	tests := make(tTable, 0)
	for _, line := range lines {
		switch {
		case !before && !after:
			if _, err := fmt.Sscanf(line, "Before: [%d, %d, %d, %d]", &t.before[0], &t.before[1], &t.before[2], &t.before[3]); err == nil {
				before = true
			}
		case before && !after:
			if _, err := fmt.Sscanf(line, "%d %d %d %d", &t.op.code, &t.op.a, &t.op.b, &t.op.c); err == nil {
				after = true
			}
		case before && after:
			if _, err := fmt.Sscanf(line, "After:  [%d, %d, %d, %d]", &t.after[0], &t.after[1], &t.after[2], &t.after[3]); err == nil {
				before = false
				after = false
				tests = append(tests, t)
			}
		}
	}
	return tests
}
