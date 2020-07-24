package main

import "testing"

func Test_calcFuel(t *testing.T) {
	tests := []struct {
		mass int
		want int64
	}{
		{12, 2},
		{14, 2},
		{1969, 654},
		{100756, 33583},
	}
	for _, tt := range tests {
		t.Run(string(tt.mass), func(t *testing.T) {
			if got := calcFuel(tt.mass); got != tt.want {
				t.Errorf("calcFuel() = %v, want %v", got, tt.want)
			}
		})
	}
}
