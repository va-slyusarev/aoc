package main

import "testing"

func Test_game_play(t *testing.T) {
	tests := []struct {
		name  string
		lines []string
		want1 int
		want2 int
		want3 int
	}{
		{"one",
			[]string{
				"#######",
				"#.G...#",
				"#...EG#",
				"#.#.#G#",
				"#..G#E#",
				"#.....#",
				"#######",
			},
			29, 172, 4988,
		},
		{"two",
			[]string{
				"#######",
				"#E..EG#",
				"#.#G.E#",
				"#E.##E#",
				"#G..#.#",
				"#..E#.#",
				"#######",
			},
			33, 948, 31284,
		},
		/*
			{"three",
				[]string{
					"#######",
					"#E.G#.#",
					"#.#G..#",
					"#G.#.G#",
					"#G..#.#",
					"#...E.#",
					"#######",
				},
				37, 94, 3478,
			},
			{"four",
				[]string{
					"#######",
					"#.E...#",
					"#.#..G#",
					"#.###.#",
					"#E#G#G#",
					"#...#G#",
					"#######",
				},
				39, 166, 6474,
			},
			{"five",
				[]string{
					"#########",
					"#G......#",
					"#.E.#...#",
					"#..##..G#",
					"#...##..#",
					"#...#...#",
					"#.G...G.#",
					"#.....G.#",
					"#########",
				},
				30, 38, 1140,
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var g game
			dmg := 4
			for {
				g = newGame(tt.lines)
				g.players.setDamage(elf, dmg)

				for ; !g.end(); g.playRound() {

				}

				if g.elfsWin() {
					break
				}
				dmg++
			}
			got1, got2, got3 := g.round, g.players.livingHP(), g.round*g.players.livingHP()

			if got1 != tt.want1 || got2 != tt.want2 || got3 != tt.want3 {
				t.Errorf("(got=want) round: %v=%v, hp: %v=%v, num: %v=%v", got1, tt.want1, got2, tt.want2, got3, tt.want3)
			}
		})
	}
}
