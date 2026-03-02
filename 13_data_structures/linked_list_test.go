package main

import "testing"

func TestAppend(t *testing.T) {
	tests := []struct {
		name     string
		values   []int
		expected int
	}{
		{"empty", []int{}, 0},
		{"one element", []int{1}, 1},
		{"three elements", []int{1, 2, 3}, 3},
		{"six elements", []int{6, 5, 4, 3, 2, 1}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := NewLinkedList[int]()
			for _, v := range tt.values {
				l.Append(v)
			}
			if l.Size() != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, l.Size())
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		idx      int
		expected bool
	}{
		{"empty", []int{}, 0, false},
		{"single element", []int{1}, 0, true},
		{"two elements (head)", []int{1, 2}, 0, true},
		{"two elements (tail)", []int{1, 2}, 1, true},
		{"two elements (out of range positive)", []int{1, 2}, 2, false},
		{"two elements (out of range negative)", []int{1, 2}, -1, false},
		{"three elements (head)", []int{1, 2, 3}, 0, true},
		{"three elements (middle)", []int{1, 2, 3}, 1, true},
		{"three elements (tail)", []int{1, 2, 3}, 2, true},
		{"three elements (out of range)", []int{1, 2, 3}, 10, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := NewLinkedList[int]()
			for _, v := range tt.data {
				l.Append(v)
			}
			if ok := l.Delete(tt.idx); ok != tt.expected {
				t.Errorf("expected %t, got %t", tt.expected, ok)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		data     []int
		idx      int
		expected int
		ok       bool
	}{
		{"empty", []int{}, 0, 0, false},
		{"one element", []int{1}, 0, 1, true},
		{"one element (out of range positive)", []int{1}, 1, 0, false},
		{"one element (out of range negative)", []int{1}, -1, 0, false},
		{"two elements (head)", []int{1, 2}, 0, 1, true},
		{"two elements (tail)", []int{1, 2}, 1, 2, true},
		{"three elements (head)", []int{1, 2, 3}, 0, 1, true},
		{"three elements (middle)", []int{1, 2, 3}, 1, 2, true},
		{"three elements (tail)", []int{1, 2, 3}, 2, 3, true},
		{"three elements (out of range positive)", []int{1, 2, 3}, 4, 0, false},
		{"three elements (out of range negative)", []int{1, 2, 3}, -10, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			l := NewLinkedList[int]()
			for _, v := range tt.data {
				l.Append(v)
			}
			res, ok := l.Get(tt.idx)
			// validate both only if it data should have been returned
			if ok != tt.ok {
				t.Errorf("expected %t, got %t", tt.ok, ok)
			}
			if tt.ok && res != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, res)
			}
		})
	}
}
