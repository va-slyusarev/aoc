package main

import (
	"reflect"
	"testing"
)

func Test_pots_zoom(t *testing.T) {
	tests := []struct {
		name string
		p    pots
		size int
		want pots
	}{
		{"one",
			pots{
				{order: 0, value: "#"},
			},
			2,
			pots{
				{order: -2, value: "."},
				{order: -1, value: "."},
				{order: 0, value: "#"},
				{order: 1, value: "."},
				{order: 2, value: "."},
			},
		},
		{"two",
			pots{
				{order: 4, value: "#"},
				{order: 5, value: "#"},
			},
			3,
			pots{
				{order: 1, value: "."},
				{order: 2, value: "."},
				{order: 3, value: "."},
				{order: 4, value: "#"},
				{order: 5, value: "#"},
				{order: 6, value: "."},
				{order: 7, value: "."},
				{order: 8, value: "."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.p.zoom(tt.size); !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("pots.zoom() =\n%v\nwant\n%v", tt.p, tt.want)
			}
		})
	}
}

func Test_pots_trim(t *testing.T) {
	tests := []struct {
		name string
		p    pots
		want pots
	}{
		{"one",
			pots{
				{order: 0, value: "#"},
			},
			pots{
				{order: 0, value: "#"},
			},
		},
		{"two",
			pots{
				{order: -2, value: "."},
				{order: -1, value: "."},
				{order: 0, value: "#"},
				{order: 1, value: "."},
				{order: 2, value: "."},
				{order: 3, value: "#"},
			},
			pots{
				{order: 0, value: "#"},
				{order: 1, value: "."},
				{order: 2, value: "."},
				{order: 3, value: "#"},
			},
		},
		{"three",
			pots{
				{order: 0, value: "."},
			},
			pots{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.p.trim(); !reflect.DeepEqual(tt.p, tt.want) {
				t.Errorf("pots.trim() =\n%v\nwant\n%v", tt.p, tt.want)
			}
		})
	}
}
