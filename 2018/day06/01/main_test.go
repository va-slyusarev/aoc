package main

import (
	"fmt"
	"testing"
)

func Test_point_takeLess(t *testing.T) {
	tests := []struct {
		a, b point
		want point
	}{
		{point{0, 0}, point{0, 0}, point{0, 0}},
		{point{10, 10}, point{0, 0}, point{0, 0}},
		{point{10, 10}, point{5, 4}, point{5, 4}},
		{point{10, 3}, point{5, 4}, point{5, 3}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("it - %d", i), func(t *testing.T) {
			if got := tt.a.takeLess(tt.b); got != tt.want {
				t.Errorf("point.takeLess() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_point_takeMore(t *testing.T) {
	tests := []struct {
		a, b point
		want point
	}{
		{point{0, 0}, point{0, 0}, point{0, 0}},
		{point{10, 10}, point{0, 0}, point{10, 10}},
		{point{10, 10}, point{5, 4}, point{10, 10}},
		{point{10, 3}, point{5, 4}, point{10, 4}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("it - %d", i), func(t *testing.T) {
			if got := tt.a.takeMore(tt.b); got != tt.want {
				t.Errorf("point.takeMore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_point_manhattanDist(t *testing.T) {
	tests := []struct {
		a, b point
		want int
	}{
		{point{0, 0}, point{0, 0}, 0},
		{point{1, 2}, point{0, 0}, 3},
		{point{0, 0}, point{10, 10}, 20},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("it - %d", i), func(t *testing.T) {
			if got := tt.a.manhattanDist(tt.b); got != tt.want {
				t.Errorf("point.manhattanDist() = %v, want %v", got, tt.want)
			}
		})
	}
}

