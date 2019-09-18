package main

import (
	"bufio"
	"bytes"
	"container/ring"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var input = flag.String("input", "../input.txt", "Input data file path.")
var divisor = flag.Int("d", 23, "Divisor for other rules.")
var numRemove = flag.Int("r", 7, "Number marble counter-clockwise remove for other rules.")
var multiplier = flag.Int("m", 100, "Multiplier for Part Two.")

type player struct {
	id    int
	score int
}

func main() {
	flag.Parse()

	f, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Printf("error read input data file: %v\n", err)
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(f))

	type input struct {
		players int
		marble  int
		score   int
	}
	var lines []*input

	for scanner.Scan() {
		value := scanner.Text()
		var players, marble, score int

		if strings.Contains(value, "high score is") {
			if _, err := fmt.Sscanf(value, "%d players; last marble is worth %d points: high score is %d", &players, &marble, &score); err != nil {
				fmt.Printf("broken data %q from input data file: %v\n", value, err)
				return
			}
		} else {
			if _, err := fmt.Sscanf(value, "%d players; last marble is worth %d points", &players, &marble); err != nil {
				fmt.Printf("broken data %q from input data file: %v\n", value, err)
				return
			}
		}

		lines = append(lines, &input{players: players, marble: marble, score: score})
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("error scan input data file: %v\n", err)
		return
	}

	for _, in := range lines {
		height := heightScore(in.players, in.marble)
		fmt.Printf("Success! For %2d players and %4d marlbe height score is: %d (%d)\n", in.players, in.marble, height, in.score)

	}
}

func heightScore(players, marble int) int {
	// init players
	p := ring.New(players)
	for i := 1; i <= players; i++ {
		p.Value = player{id: i}
		p = p.Next()
	}

	// init game
	game := ring.New(1)
	game.Value = 0

	for ball := 1; ball <= marble**multiplier; ball++ {

		switch ball%*divisor == 0 {
		case true:
			currPlayer := p.Value.(player)
			currPlayer.score += ball + game.Move(-*numRemove).Value.(int)
			p.Value = currPlayer
			game = game.Move(-*numRemove - 1)
			game.Unlink(1)
			game = game.Next()

		case false:
			newBall := ring.New(1)
			newBall.Value = ball
			game = game.Next()
			game.Link(newBall)
			game = game.Next()
		}

		p = p.Next()
	}

	// find height score
	height := 0
	for i := 0; i < players; i++ {
		score := p.Value.(player).score
		if score > height {
			height = score
		}
		p = p.Next()
	}

	return height
}
