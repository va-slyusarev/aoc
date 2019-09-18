package main

import (
	"strconv"
	"testing"
)

func Test_lab_last10(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{9, "5158916779"},
		{5, "0124515891"},
		{18, "9251071085"},
		{2018, "5941429882"},
	}
	for _, tt := range tests {
		t.Run(strconv.Itoa(tt.n), func(t *testing.T) {
			lab := lab{[]int{3, 7}, 0, 1}
			if got := lab.last10(tt.n); got != tt.want {
				t.Errorf("last10() = %v, want %v", got, tt.want)
			}
		})
	}
}
