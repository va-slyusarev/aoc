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
		/*
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
				47, 590, 27730,
			},
		*/
		{"two",
			[]string{
				"#######",
				"#G..#E#",
				"#E#E.E#",
				"#G.##.#",
				"#...#E#",
				"#...E.#",
				"#######",
			},
			37, 982, 36334,
		},
		{"three",
			[]string{
				"#######",
				"#E..EG#",
				"#.#G.E#",
				"#E.##E#",
				"#G..#.#",
				"#..E#.#",
				"#######",
			},
			46, 859, 39514,
		},
		{"four",
			[]string{
				"#######",
				"#E.G#.#",
				"#.#G..#",
				"#G.#.G#",
				"#G..#.#",
				"#...E.#",
				"#######",
			},
			35, 793, 27755,
		},
		{"five",
			[]string{
				"#######",
				"#.E...#",
				"#.#..G#",
				"#.###.#",
				"#E#G#G#",
				"#...#G#",
				"#######",
			},
			54, 536, 28944,
		},
		{"six",
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
			20, 937, 18740,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newGame(tt.lines)
			for ; !g.end(); g.playRound() {

			}
			got1, got2, got3 := g.round, g.players.livingHP(), g.round*g.players.livingHP()

			if got1 != tt.want1 || got2 != tt.want2 || got3 != tt.want3 {
				t.Errorf("(got=want) round: %v=%v, hp: %v=%v, num: %v=%v", got1, tt.want1, got2, tt.want2, got3, tt.want3)
			}
		})
	}
}
