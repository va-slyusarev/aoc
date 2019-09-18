package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strings"
)

var input = flag.String("input", "../input.txt", "Input data file path.")

const (
	elf    = 'E'
	goblin = 'G'

	wall   = '#'
	cavern = '.'

	hp     = 200
	damage = 3
)

type point struct {
	x, y int
}

func (p point) String() string { return fmt.Sprintf("%d,%d", p.x, p.y) }

func (p *point) plus(value point) { p.x += value.x; p.y += value.y }

func (p *point) copy() point { return point{p.x, p.y} }

func (p point) equals(value *point) bool { return value != nil && p.x == value.x && p.y == value.y }

func (p *point) takeMore(value point) {
	if value.x > p.x {
		p.x = value.x
	}
	if value.y > p.y {
		p.y = value.y
	}
}

func (p point) less(value point) bool {
	if p.y != value.y {
		return p.y < value.y
	}
	return p.x < value.x
}

type unit struct {
	p      point
	kind   rune
	hp     int
	damage int
}

func (u *unit) String() string { return fmt.Sprintf("%v(%d)", string(u.kind), u.hp) }

func (u *unit) move(p *point) { u.p = *p }

func (u *unit) isDead() bool { return u.hp <= 0 }

func (u *unit) isEnemy(value *unit) bool { return u.kind != value.kind }

func (u *unit) attack(value *unit) { value.hp -= u.damage }

func (u *unit) enemy(players units) *unit {
	var foe *unit
	dirs := []point{{0, -1}, {-1, 0}, {1, 0}, {0, 1}}
	for _, dir := range dirs {
		p := u.p.copy()
		p.plus(dir)
		for _, player := range players {
			if p.equals(&player.p) && u.isEnemy(player) && !player.isDead() {
				// either a new enemy, or an enemy with less hp
				if foe == nil || player.hp < foe.hp {
					foe = player
				}
			}
		}
	}
	return foe
}

func (u *unit) enemies(players units) units {
	var enemies units
	for _, val := range players {
		if val.kind != u.kind && !val.isDead() {
			enemies = append(enemies, val)
		}
	}
	sort.Sort(enemies)
	return enemies
}

func (u *unit) path(g game, goal *unit) (int, *point) {

	frontier := []point{u.p}
	paths := make(map[point]point)

	find := false
	for len(frontier) > 0 {
		cur := frontier[0]
		frontier = frontier[1:]

		if cur.equals(&goal.p) {
			find = true
			break
		}

		for _, p := range g.neighbors(cur, *goal) {
			if _, ok := paths[p]; !ok {
				frontier = append(frontier, p)
				paths[p] = cur
			}
		}
	}

	if !find {
		return -1, nil
	}

	cur := goal.p
	path := []point{cur}

	for !cur.equals(&u.p) {
		cur = paths[cur]
		path = append(path, cur)
	}
	//fmt.Printf("%v - %v %v\n", u, goal, path)

	dp := path[len(path)-2]

	return len(path) - 2, &dp
}

type units []*unit

func (u units) Len() int           { return len(u) }
func (u units) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }
func (u units) Less(i, j int) bool { return u[i].p.less(u[j].p) }

func (u units) livingHP() int {
	totalHP := 0
	for _, unit := range u {
		if !unit.isDead() {
			totalHP += unit.hp
		}
	}
	return totalHP
}

func (u units) setDamage(kind rune, dmg int) {
	for _, unit := range u {
		if unit.kind == kind {
			unit.damage = dmg
		}
	}
}

type game struct {
	board map[point]rune
	max   point

	round int

	players units
}

func newGame(lines []string) game {
	g := game{
		board:   make(map[point]rune),
		max:     point{0, 0},
		round:   0,
		players: make(units, 0),
	}

	for y, line := range lines {
		for x, r := range line {
			p := point{x, y}

			switch r {
			case elf, goblin:
				g.players = append(g.players, &unit{p: p, kind: r, hp: hp, damage: damage})
				g.board[p] = cavern
			case wall, cavern:
				g.board[p] = r
			default:
				continue
			}
			g.max.takeMore(p)
		}
	}
	return g
}

func (g *game) neighbors(p point, goal unit) []point {
	result := make([]point, 0)
	dirs := []point{{0, -1}, {-1, 0}, {1, 0}, {0, 1}}

	for _, dir := range dirs {
		n := p.copy()
		n.plus(dir)

		// the goal we have to consider in the list of neighbors
		if n.equals(&goal.p) {
			result = append(result, n)
			continue
		}

		findPlayer := false
		for _, player := range g.players {
			if n.equals(&player.p) && !player.isDead() {
				findPlayer = true
				break
			}
		}
		if findPlayer {
			continue
		}

		if val := g.board[n]; val == cavern {
			result = append(result, n)
		}
	}
	return result
}

func (g game) playerByPoint(p point) *unit {
	for _, player := range g.players {
		if player.p.equals(&p) {
			return player
		}
	}
	return nil
}

func (g game) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "Round: #%d\n", g.round)

	for y := 0; y <= g.max.y; y++ {

		var addInfo strings.Builder
		for x := 0; x <= g.max.x; x++ {
			cur := point{x, y}

			if player := g.playerByPoint(cur); player != nil && !player.isDead() {
				_, _ = fmt.Fprint(&sb, string(player.kind))
				_, _ = fmt.Fprintf(&addInfo, " %s", player)
				continue
			}

			if v, ok := g.board[cur]; ok {
				_, _ = fmt.Fprint(&sb, string(v))
			}
		}
		_, _ = fmt.Fprintf(&sb, " %s\n", addInfo.String())
	}
	return sb.String()
}

func (g *game) playRound() {

	sort.Sort(g.players)
	for _, player := range g.players {

		if player.isDead() {
			continue
		}

		if enemy := player.enemy(g.players); enemy != nil {
			player.attack(enemy)
			continue
		}

		min := math.MaxInt64
		var dp *point
		for _, enemy := range player.enemies(g.players) {
			if l, p := player.path(*g, enemy); l < min && p != nil {
				min, dp = l, p
			}
		}

		if dp != nil {
			player.move(dp)

			// attack
			if enemy := player.enemy(g.players); enemy != nil {
				player.attack(enemy)
				continue
			}
		}
	}

	if !g.end() {
		g.round++
	}
}

func (g *game) end() bool {
	e, gob := false, false

	for _, player := range g.players {
		if !player.isDead() {

			switch player.kind {
			case elf:
				e = true
			case goblin:
				gob = true
			}
		}
		if e && gob {
			return false
		}
	}
	return true
}

func (g *game) elfsWin() bool {
	for _, player := range g.players {
		if (player.kind == elf && player.isDead()) || (player.kind == goblin && !player.isDead()) {
			return false
		}
	}
	return true
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

	var g game
	dmg := 4
	for {
		g = newGame(lines)
		g.players.setDamage(elf, dmg)

		for ; !g.end(); g.playRound() {

		}

		if g.elfsWin() {
			break
		}
		dmg++
	}

	fmt.Printf("Success! Target number is (%d * %d): %d\n", g.round, g.players.livingHP(), g.round*g.players.livingHP())
}
