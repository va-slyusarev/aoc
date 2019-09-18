package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
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

	ip, ops, err := parse(lines)
	if err != nil {
		fmt.Printf("error scan iput data file: %v\n", err)
		return
	}

	now := time.Now()

	_ = run(0, ops, ip, 0)

	fmt.Printf("Success! Target number is: %d (elapsed time: %s)\n", manually(), time.Since(now))
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

// for analyze
func run(reg0 int, ops operations, ip int, n int) int {
	var regs registers
	regs[0] = reg0

	step := 0
	for i := 0; i < len(ops) && step < n; i++ {
		var sb strings.Builder
		_, _ = fmt.Fprintf(&sb, "step: %6d, i = %2d, ", step, i)
		regs[ip] = i
		_, _ = fmt.Fprintf(&sb, "%v, ", regs)
		regs = nameDo[ops[i].name](regs, ops[i])
		_, _ = fmt.Fprintf(&sb, "%v, %v, ", ops[i], regs)
		i = regs[ip]
		_, _ = fmt.Fprintf(&sb, "next i = %d", i+1)
		step++
		fmt.Println(sb.String())
	}

	if step >= n {
		return -1
	}
	return regs[0]
}

// I had to analyze the input data. As a result, you should get into 28 instructions (eqrr 3 0 4) and the desired value
// reg0 should be equal to the value reg3. To get to the 28th instruction it is necessary through 16 (seti 27 6 1).
// And in 16 through 14 (addr 4 1 1) provided reg4 == 1. This condition will be achievable with a "positive" outcome of
// 13 instructions (gtir 256 2 4).
// Actually, everything will revolve around operations with registers reg2 and reg3 (commands 6-12).
func manually() int {
	prev := 65536
	curr := 1505483

	for {
		curr = (curr + prev&255) * 65899 & 16777215
		if prev < 256 {
			return curr
		}
		prev /= 256
	}
}
