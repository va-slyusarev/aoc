package main

import "testing"

func Test_fullCalc(t *testing.T) {
	tests := []struct {
		mass int
		want int64
	}{
		{14, 2},
		{1969, 966},
		{100756, 50346},
	}
	for _, tt := range tests {
		t.Run(string(tt.mass), func(t *testing.T) {
			if got := fullCalc(tt.mass); got != tt.want {
				t.Errorf("fullCalc() = %v, want %v", got, tt.want)
			}
		})
	}
}
