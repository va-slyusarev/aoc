package main

import (
	"reflect"
	"testing"
)

func Test_tree_len(t *testing.T) {
	tests := []struct {
		t    tree
		want int
	}{
		{},
		{nil, 0},
		{tree{&node{}}, 1},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			if got := tt.t.len(); got != tt.want {
				t.Errorf("len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tree_pop(t *testing.T) {
	tests := []struct {
		t       tree
		want    *node
		wantErr bool
	}{
		{tree{}, nil, true},
		{tree{&node{}}, &node{}, false},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			got, err := tt.t.pop()
			if (err != nil) != tt.wantErr {
				t.Errorf("pop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pop() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tree_push(t *testing.T) {
	tests := []struct {
		t    tree
		n    *node
		want tree
	}{
		{tree{}, &node{}, tree{&node{}}},
		{tree{&node{child: 1}}, &node{child: 2}, tree{&node{child: 1}, &node{child: 2}}},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			if tt.t.push(tt.n); !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("push() = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func Test_tree_touch(t *testing.T) {
	tests := []struct {
		t    tree
		want *node
	}{
		{},
		{tree{}, nil},
		{tree{&node{}}, &node{}},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			if got := tt.t.touch(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("touch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tree_amend(t *testing.T) {
	tests := []struct {
		t    tree
		n    *node
		want tree
	}{
		{tree{}, &node{}, tree{&node{}}},
		{tree{&node{child: 3}}, &node{child: 2}, tree{&node{child: 2}}},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			if tt.t.amend(tt.n); !reflect.DeepEqual(tt.t, tt.want) {
				t.Errorf("amend() = %v, want %v", tt.t, tt.want)
			}
		})
	}
}

func Test_licbytes_popNode(t *testing.T) {
	tests := []struct {
		name    string
		l       licbytes
		want    byte
		want1   byte
		wantErr bool
	}{
		{"0", licbytes{0, 1, 2}, 0, 1, false},
		{"1", licbytes{0}, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := tt.l.popNode()
			if (err != nil) != tt.wantErr {
				t.Errorf("popNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("popNode() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("popNode() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_licbytes_popSum(t *testing.T) {
	type args struct {
		n byte
	}
	tests := []struct {
		name    string
		l       licbytes
		args    args
		want    int
		wantErr bool
	}{
		{"0", licbytes{1, 2, 3}, args{3}, 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.l.popSum(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("popSum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("popSum() got = %v, want %v", got, tt.want)
			}
		})
	}
}
