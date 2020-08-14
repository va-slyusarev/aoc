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
		{digits{1, 1, 2, 2, 3, 3}, true},
		{digits{1, 2, 3, 4, 4, 4}, false},
		{digits{1, 1, 1, 1, 2, 2}, true},
		{digits{1, 1, 1, 2, 2, 2}, false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.d), func(t *testing.T) {
			if got := sameDigits(tt.d); got != tt.want {
				t.Errorf("sameDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}
