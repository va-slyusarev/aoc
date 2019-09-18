package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

var input = flag.String("input", "../input.txt", "Input data file path.")
var generation = flag.Int("g", 20, "Number of generations.")

type rules map[string]string

type pot struct {
	order int
	value string
}

func (p *pot) String() string {
	return fmt.Sprintf("%d(%s)", p.order, p.value)
}

type pots []*pot

func (p pots) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "🍲\n")
	for _, v := range p {
		switch v.order % 1 {
		case 0:
			_, _ = fmt.Fprintf(&sb, "%-3d", v.order)
		default:
			_, _ = fmt.Fprintf(&sb, "%s", strings.Repeat(" ", 3))
		}
	}
	_, _ = fmt.Fprintln(&sb)

	for _, v := range p {
		_, _ = fmt.Fprintf(&sb, "%-3s", v.value)
	}

	return sb.String()
}

// increase by the number of empty pots left and right
func (p *pots) zoom(size int) {
	var left, right pots

	min, max := 0, 0
	if len(*p) > 0 {
		min = (*p)[0].order
		max = (*p)[len(*p)-1].order
	}

	for i := 1; i <= size; i++ {
		min -= 1
		max += 1
		left = append(pots{{order: min, value: "."}}, left...)
		right = append(right, &pot{order: max, value: "."})
	}

	*p = append(left, *p...)
	*p = append(*p, right...)
}

// remove the empty pots on the left and right
func (p *pots) trim() {

	// trim left
	for len(*p) > 0 {
		if (*p)[0].value == "#" {
			break
		}
		*p = (*p)[1:]
	}

	// trim right
	for len(*p) > 0 {
		if (*p)[len(*p)-1].value == "#" {
			break
		}
		*p = (*p)[:len(*p)-1]
	}
}

func (p pots) sumOrderSharp() int {
	sum := 0
	for _, it := range p {
		if it.value == "#" {
			sum += it.order
		}
	}
	return sum
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

	p := make(pots, 0)
	r := make(rules)

	for i, line := range lines {

		switch i {
		case 0:
			if err := parseState(line, &p); err != nil {
				fmt.Printf("error parse initial state %q: %v\n", line, err)
				return
			}
		case 1:
			continue
		default:
			if err := parseRule(line, r); err != nil {
				fmt.Printf("error parse rule %q: %v\n", line, err)
				return
			}
		}
	}

	for g := 1; g <= *generation; g++ {
		p.zoom(4)
		n := make(pots, len(p))
		for i := range p {
			n[i] = &pot{order: p[i].order, value: p[i].value}
		}

		for i := 2; i < len(p)-2; i++ {
			neighbourhood := p[i-2].value + p[i-1].value + p[i].value + p[i+1].value + p[i+2].value
			if v, ok := r[neighbourhood]; ok {
				n[i].value = v
				continue
			}
			n[i].value = "."
		}

		for i := range n {
			p[i].value = n[i].value
		}
		p.trim()
	}

	fmt.Printf("Success! Target number is: %d\n", p.sumOrderSharp())
}

func parseState(value string, p *pots) error {
	var re = regexp.MustCompile("initial state:\\s([.#]+)")
	values := re.FindStringSubmatch(value)
	if len(values) != 2 {
		return errors.New("broken state")
	}

	for i, v := range values[1] {
		*p = append(*p, &pot{order: i, value: string(v)})
	}
	return nil
}

func parseRule(value string, r rules) error {
	var re = regexp.MustCompile("([.#]{5})\\s=>\\s([.#])")
	values := re.FindStringSubmatch(value)
	if len(values) != 3 {
		return errors.New("broken rule")
	}
	r[values[1]] = values[2]
	return nil
}
