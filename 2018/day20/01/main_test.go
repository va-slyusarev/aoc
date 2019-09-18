package main

import (
	"testing"
)

func Test_layout_maxDist(t *testing.T) {
	tests := []struct {
		regexp string
		want   int
	}{
		{"^WNE$", 3},
		{"^ENWWW(NEEE|SSE(EE|N))$", 10},
		{"^ENNWSWW(NEWS|)SSSEEN(WNSE|)EE(SWEN|)NNN$", 18},
		{"^ESSWWN(E|NNENN(EESS(WNSE|)SSS|WWWSSSSE(SW|NNNE)))$", 23},
		{"^WSSEESWWWNW(S|NENNEEEENN(ESSSSW(NWSW|SSEN)|WSWWN(E|WWS(E|SS))))$", 31},
	}
	for _, tt := range tests {
		t.Run(tt.regexp, func(t *testing.T) {
			model := build(tt.regexp)
			if got := model.maxDist(); got != tt.want {
				t.Errorf("maxDist() = %v, want %v", got, tt.want)
			}
		})
	}
}
