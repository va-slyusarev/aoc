package main

import "testing"

func Test_process(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"1,0,0,0,99", []int{1, 0, 0, 0, 99}, 2},
		{"2,3,0,3,99", []int{2, 3, 0, 3, 99}, 2},
		{"2,4,4,5,99,0", []int{2, 4, 4, 5, 99, 0}, 2},
		{"1,1,1,4,99,5,6,0,99", []int{1, 1, 1, 4, 99, 5, 6, 0, 99}, 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := process(tt.nums); got != tt.want {
				t.Errorf("process() = %v, want %v", got, tt.want)
			}
		})
	}
}
