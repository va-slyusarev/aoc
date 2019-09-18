package main

import (
	"reflect"
	"testing"
)

func Test_item_next(t *testing.T) {
	tests := []struct {
		it      item
		want    string
		wantErr bool
	}{
		{map[string]dependency{"C": dependency{hold: false, dep: nil}}, "C", false},
		{map[string]dependency{"C": dependency{hold: true, dep: nil}, "A": dependency{hold: false, dep: nil}}, "A", false},
		{map[string]dependency{"C": dependency{hold: true, dep: nil}, "A": dependency{hold: true, dep: nil}}, "", true},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			got, err := tt.it.next()
			if (err != nil) != tt.wantErr {
				t.Errorf("item.next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("item.next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_item_delDependency(t *testing.T) {
	tests := []struct {
		it   item
		dep  string
		want item
	}{
		{
			map[string]dependency{
				"C": dependency{hold: true, dep: nil},
				"A": dependency{hold: false, dep: []string{"C", "D"}}},
			"C",
			map[string]dependency{
				"A": dependency{hold: false, dep: []string{"D"}}},
		},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			if tt.it.delDependency(tt.dep); !reflect.DeepEqual(tt.it, tt.want) {
				t.Errorf("worker.delDependency() = %v, want %v", tt.it, tt.want)
			}
		})
	}
}

func Test_worker_tick(t *testing.T) {
	tests := []struct {
		w    worker
		want []string
	}{
		{map[string]task{"1": task{time: 1, item: "A"}}, []string{"A"}},
		{map[string]task{"1": task{time: 1, item: "A"}, "2": task{time: 12, item: "B"}}, []string{"A"}},
		{map[string]task{"1": task{time: 0, item: "A"}, "2": task{time: 12, item: "B"}}, nil},
	}
	for i, tt := range tests {
		t.Run(string(i), func(t *testing.T) {
			if got := tt.w.tick(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("worker.tick() = %v, want %v", got, tt.want)
			}
		})
	}
}
