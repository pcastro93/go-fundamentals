package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

type TrieNode[T any] struct {
	data     T
	hasData  bool
	children map[rune]*TrieNode[T]
}

type Trie[T any] struct {
	root *TrieNode[T]
}

func NewTrie[T any]() *Trie[T] {
	return &Trie[T]{root: &TrieNode[T]{}}
}

func (t *Trie[T]) Insert(key string, data T) error {
	if t == nil {
		return fmt.Errorf("trie is nil")
	}
	if t.root == nil {
		t.root = &TrieNode[T]{}
	}
	cur := t.root
	for _, r := range key {
		if cur.children == nil {
			cur.children = make(map[rune]*TrieNode[T])
		}
		if _, ok := cur.children[r]; !ok {
			cur.children[r] = &TrieNode[T]{}
		}
		cur = cur.children[r]
	}
	if cur.hasData {
		return fmt.Errorf("key already exists")
	}
	cur.data = data
	cur.hasData = true
	return nil
}

func (t *Trie[T]) Upsert(key string, initial T, update func(T) T) {
	cur := t.root
	for _, r := range key {
		if cur.children == nil {
			cur.children = make(map[rune]*TrieNode[T])
		}
		if _, ok := cur.children[r]; !ok {
			cur.children[r] = &TrieNode[T]{}
		}
		cur = cur.children[r]
	}
	if cur.hasData {
		cur.data = update(cur.data)
	} else {
		cur.data = initial
		cur.hasData = true
	}
}

// Size returns an estimate of the memory consumed by the trie in bytes.
// Map internals are approximated: ~128 bytes of header overhead per node
// that has children, plus 12 bytes per child entry (4-byte rune key + 8-byte pointer).
// For string data, the backing array bytes are added on top of the header.
func (t *Trie[T]) Size() uintptr {
	if t == nil || t.root == nil {
		return unsafe.Sizeof(*t)
	}
	return unsafe.Sizeof(*t) + trieNodeSize(t.root)
}

func trieNodeSize[T any](n *TrieNode[T]) uintptr {
	if n == nil {
		return 0
	}

	size := unsafe.Sizeof(*n)

	// For string data, unsafe.Sizeof only counts the string header (16 bytes).
	// Add the length of the backing array.
	v := reflect.ValueOf(n.data)
	if v.Kind() == reflect.String {
		size += uintptr(v.Len())
	}

	// Map header + per-entry cost (rune key + pointer value).
	if len(n.children) > 0 {
		const mapHeaderBytes = 128
		const mapEntryBytes = 12 // unsafe.Sizeof(rune(0)) + unsafe.Sizeof((*TrieNode[T])(nil))
		size += mapHeaderBytes + uintptr(len(n.children))*mapEntryBytes
	}

	for _, child := range n.children {
		size += trieNodeSize(child)
	}

	return size
}

func (t *Trie[T]) Search(key string) (T, error) {
	var zero T
	if t == nil {
		return zero, nil
	}
	if t.root == nil {
		return zero, nil
	}
	cur := t.root
	for _, r := range key {
		if cur.children == nil {
			return zero, nil
		}
		if _, ok := cur.children[r]; !ok {
			return zero, nil
		}
		cur = cur.children[r]
	}
	if !cur.hasData {
		return zero, nil
	}
	return cur.data, nil
}
