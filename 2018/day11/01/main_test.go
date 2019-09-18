package main

import "testing"

func Test_point_power(t *testing.T) {
	tests := []struct {
		name   string
		p      point
		serial int
		want   int
	}{
		{"input 0", point{3, 5}, 8, 4},
		{"input 1", point{122, 79}, 57, -5},
		{"input 2", point{217, 196}, 39, 0},
		{"input 3", point{101, 153}, 71, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.power(tt.serial); got != tt.want {
				t.Errorf("power() = %v, want %v", got, tt.want)
			}
		})
	}
}
