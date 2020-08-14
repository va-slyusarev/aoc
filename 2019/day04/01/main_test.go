package main

import (
	"fmt"
	"testing"
)

func Test_sameDigits(t *testing.T) {
	tests := []struct {
		d    digits
		want bool
	}{
		{digits{0, 0, 0, 0, 0, 0}, true},
		{digits{1, 2, 3, 0, 5, 5}, true},
		{digits{1, 2, 3, 4, 5, 6}, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.d), func(t *testing.T) {
			if got := sameDigits(tt.d); got != tt.want {
				t.Errorf("sameDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_neverDecrease(t *testing.T) {
	tests := []struct {
		d    digits
		want bool
	}{
		{digits{0, 0, 0, 0, 0, 0}, true},
		{digits{1, 2, 3, 4, 5, 6}, true},
		{digits{1, 2, 3, 4, 5, 4}, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.d), func(t *testing.T) {
			if got := neverDecrease(tt.d); got != tt.want {
				t.Errorf("neverDecrease() = %v, want %v", got, tt.want)
			}
		})
	}
}
