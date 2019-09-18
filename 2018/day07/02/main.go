package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
)

var input = flag.String("input", "../input.txt", "Input data file path.")
var number = flag.Int("w", 5, "Number of workers.")
var step = flag.Int("t", 60, "Time per step.")

type item map[string]dependency

type dependency struct {
	hold bool
	dep  []string
}

// get next item and hold
func (it item) next() (string, error) {
	var independent []string
	for k, v := range it {
		if len(v.dep) == 0 && !v.hold {
			independent = append(independent, k)
		}
	}
	// alphabetical asc
	sort.Strings(independent)

	if len(independent) > 0 {
		value := independent[0]
		dep := it[value]
		dep.hold = true
		it[value] = dep
		return value, nil
	}
	return "", fmt.Errorf("waiting")
}

func (it item) delDependency(dep string) {
	for i := range it {
		c := it[i]
		for k, v := range c.dep {
			if v == dep {
				c.dep[k] = c.dep[len(c.dep)-1]
				c.dep = c.dep[:len(c.dep)-1]
				it[i] = c
			}
		}
	}
	delete(it, dep)
}

type worker map[string]task

type task struct {
	item string
	time int
}

func (w worker) nobodyBusy() bool {
	for _, v := range w {
		if v.time > 0 {
			return false
		}
	}
	return true
}

// reduce time and return completed item to remove hold
func (w worker) tick() []string {
	var unhold []string
	for i, v := range w {
		if v.time > 0 {
			v.time -= 1
			if v.time == 0 {
				unhold = append(unhold, v.item)
				v.item = ""
			}
			w[i] = v
		}
	}
	return unhold
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	items := make(item, 0)

	for scanner.Scan() {
		value := scanner.Text()
		var dep, it string
		if _, err := fmt.Sscanf(value, "Step %s must be finished before step %s can begin.", &dep, &it); err != nil {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}

		// add new item
		if _, ok := items[dep]; !ok {
			items[dep] = dependency{}
		}

		v := items[it]
		v.dep = append(v.dep, dep)

		items[it] = v
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	time, err := buildTime(items)
	if err != nil {
		fmt.Printf("Failure...No target time found, because %v.\n", err)
		return
	}
	fmt.Printf("Success! Target time is: %d\n", time)
}

func buildTime(items item) (int, error) {

	seconds := 0
	workers := make(worker, *number)

	// init workers
	for i := 0; i < *number; i++ {
		workers[string(i)] = task{}
	}

	for {

		// delete unhold dependency
		unhold := workers.tick()
		for _, v := range unhold {
			items.delDependency(v)
		}

		for k, v := range workers {
			if v.time == 0 && len(items) > 0 {

				next, err := items.next()
				if err != nil {
					break
				}

				v.item = next
				v.time = 1 + *step + int(next[0]) - 'A'
				workers[k] = v
			}
		}

		if len(items) == 0 && workers.nobodyBusy() {
			break
		}

		seconds++
	}

	return seconds, nil
}
