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

	tests, opers := parse(lines)

	reorder(tests)

	var reg registers
	for _, op := range opers {
		reg = operations[op.code](reg, op)
	}

	fmt.Printf("Success! Target number is: %d\n", reg[0])
}

func parse(lines []string) (tTable, []operation) {
	var t test
	opers := make([]operation, 0)
	before, after := false, false
	tests := make(tTable, 0)
	for _, line := range lines {
		switch {
		case !before && !after:
			if _, err := fmt.Sscanf(line, "Before: [%d, %d, %d, %d]", &t.before[0], &t.before[1], &t.before[2], &t.before[3]); err == nil {
				before = true
			} else if _, err := fmt.Sscanf(line, "%d %d %d %d", &t.op.code, &t.op.a, &t.op.b, &t.op.c); err == nil {
				opers = append(opers, t.op)
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
	return tests, opers
}

func reorder(tbl tTable) {
	// input code: my codes
	matrix := make(map[int]map[int]struct{})
	for _, tt := range tbl {
		if _, ok := matrix[tt.op.code]; !ok {
			matrix[tt.op.code] = make(map[int]struct{})
		}
		for k, op := range operations {
			after := op(tt.before, tt.op)
			if tt.after.equals(after) {
				matrix[tt.op.code][k] = struct{}{}
			}
		}
	}

	filter := func(value int, m map[int]struct{}) map[int]struct{} {
		if len(m) == 1 {
			return m
		}
		delete(m, value)
		return m
	}

	for range matrix {
		for _, v1 := range matrix {
			if len(v1) == 1 {
				val := 0
				for k := range v1 {
					val = k
				}
				for k2, v2 := range matrix {
					matrix[k2] = filter(val, v2)
				}
			}
		}
	}

	for k, v := range matrix {

		for val := range v {
			operations[k] = operations[val]
		}
	}
}
