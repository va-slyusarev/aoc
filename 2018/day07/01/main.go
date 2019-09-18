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

// key: item, value: dependency
type item map[string][]string

func (it item) next() (string, error) {
	var independent []string
	for k, v := range it {
		if len(v) == 0 {
			independent = append(independent, k)
		}
	}
	// alphabetical asc
	sort.Strings(independent)

	if len(independent) > 0 {
		return independent[0], nil
	}
	return "", fmt.Errorf("detect circular dependency")
}

func (it item) delDependency(dep string) {
	for i, v := range it {
		for j := range v {
			if v[j] == dep {
				v[j] = v[len(v)-1]
				it[i] = v[:len(v)-1]
			}
		}
	}
	delete(it, dep)
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
			items[dep] = []string{}
		}

		items[it] = append(items[it], dep)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	sequence, err := resolveDep(items)
	if err != nil {
		fmt.Printf("Failure...No target sequence found, because %v.\n", err)
		return
	}
	fmt.Printf("Success! Target sequence is: %v\n", sequence)
}

func resolveDep(items item) (string, error) {

	resolve := ""

	for len(items) > 0 {

		next, err := items.next()
		if err != nil {
			return "", err
		}

		resolve += next
		items.delDependency(next)
	}

	return resolve, nil
}
