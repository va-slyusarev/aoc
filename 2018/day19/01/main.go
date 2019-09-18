package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type registers [6]int

func (r registers) copy() registers {
	var clone registers
	for i := range r {
		clone[i] = r[i]
	}
	return clone
}

type operation struct {
	name string
	a, b int
	c    int
}

type operations []operation

type do = func(reg registers, op operation) registers

var nameDo = map[string]do{
	"addr": addr,
	"addi": addi,
	"mulr": mulr,
	"muli": muli,
	"banr": banr,
	"bani": bani,
	"borr": borr,
	"bori": bori,
	"setr": setr,
	"seti": seti,
	"gtir": gtir,
	"gtri": gtri,
	"gtrr": gtrr,
	"eqir": eqir,
	"eqri": eqri,
	"eqrr": eqrr,
}

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

	var regs registers
	ip, ops, err := parse(lines)
	if err != nil {
		fmt.Printf("error scan iput data file: %v\n", err)
		return
	}

	for i := 0; i < len(ops); i++ {
		regs[ip] = i
		regs = nameDo[ops[i].name](regs, ops[i])
		i = regs[ip]
	}
	//fmt.Printf("%v\n", regs)
	fmt.Printf("Success! Target number is: %d\n", regs[0])
}

func parse(lines []string) (int, operations, error) {
	var ip int
	if _, err := fmt.Sscanf(lines[0], "#ip %d", &ip); err != nil {
		return 0, nil, fmt.Errorf("parsing ip value: %v", err)
	}

	ops := make(operations, 0)
	for _, line := range lines[1:] {
		op := operation{}
		if _, err := fmt.Sscanf(line, "%s %d %d %d", &op.name, &op.a, &op.b, &op.c); err != nil {
			return 0, nil, fmt.Errorf("parsing operation value: %v", err)
		}
		ops = append(ops, op)
	}

	return ip, ops, nil
}
