package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

type node struct {
	child byte
	meta  byte

	parent *node
	tchild tree
	sum    int
	leaf   bool
}

type licbytes []byte

func (l *licbytes) len() int { return len(*l) }

func (l *licbytes) pop() (byte, error) {
	if l.len() == 0 {
		return 0, fmt.Errorf("wrong pop")
	}
	b := (*l)[0]
	*l = (*l)[1:l.len()]

	return b, nil
}

// return: quantity node, quantity meta
func (l *licbytes) popNode() (byte, byte, error) {
	qnode, err := l.pop()
	if err != nil {
		return 0, 0, fmt.Errorf("wrong pop node")
	}
	qmeta, err := l.pop()
	if err != nil {
		return 0, 0, fmt.Errorf("wrong pop node")
	}
	return qnode, qmeta, nil
}

// return: sum n bytes
func (l *licbytes) popSum(n byte) (int, error) {
	sum := 0
	for i := 0; i < int(n); i++ {
		b, err := l.pop()
		if err != nil {
			return 0, fmt.Errorf("wrong pop sum")
		}
		sum += int(b)
	}
	return sum, nil
}

type tree []*node

func (t *tree) len() int { return len(*t) }

func (t *tree) push(n *node) { *t = append(*t, n) }

func (t *tree) touch() *node {
	if t.len() == 0 {
		return nil
	}
	return (*t)[t.len()-1]
}

func (t *tree) pop() (*node, error) {
	if t.len() == 0 {
		return nil, fmt.Errorf("empty")
	}
	n := t.touch()
	*t = (*t)[:t.len()-1]
	return n, nil
}

func (t *tree) amend(n *node) {
	_, _ = t.pop()
	t.push(n)
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))
	scanner.Split(bufio.ScanWords)

	lic := make(licbytes, 0)

	for scanner.Scan() {
		value := scanner.Text()
		var b byte
		if _, err := fmt.Sscanf(value, "%d", &b); err != nil {
			fmt.Printf("broken data %q from input data file: %v\n", value, err)
			return
		}

		lic = append(lic, b)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	sum, err := rootSum(lic)
	if err != nil {
		fmt.Printf("Failure...No target sum found, because %v.\n", err)
		return
	}
	fmt.Printf("Success! Target sum is: %d\n", sum)
}

func rootSum(lic licbytes) (int, error) {
	t := make(tree, 0)

	c, m, err := lic.popNode()
	if err != nil {
		return 0, err
	}
	n := &node{
		child:  c,
		meta:   m,
		parent: &node{},
		tchild: make(tree, 0),
		leaf:   c == 0,
	}
	t.push(n)

	for {

		root := t.touch()

		if root.child > 0 {

			c, m, err := lic.popNode()
			if err != nil {
				return 0, err
			}
			n := &node{
				child:  c,
				meta:   m,
				parent: root,
				tchild: make(tree, 0),
				leaf:   c == 0,
			}

			root.child -= 1
			root.tchild.push(n)
			t.amend(root)

			t.push(n)
			continue
		}

		if root.child == 0 {

			if root.leaf {
				s, err := lic.popSum(root.meta)
				if err != nil {
					return 0, err
				}
				root.sum = s
			}

			if !root.leaf {

				for i := 0; i < int(root.meta); i++ {
					idx, err := lic.pop()
					if err != nil {
						return 0, err
					}
					if int(idx) <= root.tchild.len() {
						root.sum += root.tchild[int(idx)-1].sum
					}
				}
			}

			root.parent.tchild.amend(root)

			// last element not pop -> root tree
			if lic.len() <= 0 {
				break
			}

			_, err := t.pop()
			if err != nil {
				return 0, err
			}
		}
	}

	root, err := t.pop()
	if err != nil {
		return 0, err
	}

	return root.sum, nil
}
