package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

// RadixNode is a node in a compressed radix tree (Patricia trie).
// The label field holds the edge string from the parent to this node.
type RadixNode[T any] struct {
	label    string
	data     T
	hasData  bool
	children map[rune]*RadixNode[T] // keyed by first rune of each child's label
}

// RadixTree is a compressed trie where each edge carries a string label
// instead of a single character, reducing node count for long keys.
type RadixTree[T any] struct {
	root *RadixNode[T]
}

func NewRadixTree[T any]() *RadixTree[T] {
	return &RadixTree[T]{root: &RadixNode[T]{}}
}

// commonPrefixBytes returns the byte length of the longest common UTF-8-aligned
// prefix shared by a and b.
func commonPrefixBytes(a, b string) int {
	i := 0
	for i < len(a) && i < len(b) {
		ra, sza := utf8.DecodeRuneInString(a[i:])
		rb, _ := utf8.DecodeRuneInString(b[i:])
		if ra != rb {
			break
		}
		i += sza // ra == rb implies sza == szb
	}
	return i
}

func (t *RadixTree[T]) Insert(key string, data T) error {
	if t == nil {
		return fmt.Errorf("radix tree is nil")
	}
	if t.root == nil {
		t.root = &RadixNode[T]{}
	}
	return radixInsert(t.root, key, data)
}

func radixInsert[T any](cur *RadixNode[T], remaining string, data T) error {
	if remaining == "" {
		if cur.hasData {
			return fmt.Errorf("key already exists")
		}
		cur.data = data
		cur.hasData = true
		return nil
	}

	firstRune, _ := utf8.DecodeRuneInString(remaining)
	if cur.children == nil {
		cur.children = make(map[rune]*RadixNode[T])
	}

	child, ok := cur.children[firstRune]
	if !ok {
		cur.children[firstRune] = &RadixNode[T]{label: remaining, data: data, hasData: true}
		return nil
	}

	cpLen := commonPrefixBytes(remaining, child.label)
	if cpLen == len(child.label) {
		// Edge label fully consumed — descend
		return radixInsert(child, remaining[cpLen:], data)
	}

	// Split the child edge at cpLen
	commonLabel := remaining[:cpLen]
	childSuffix := child.label[cpLen:]
	keySuffix := remaining[cpLen:]

	splitNode := &RadixNode[T]{
		label:    commonLabel,
		children: make(map[rune]*RadixNode[T]),
	}
	child.label = childSuffix
	childFirstRune, _ := utf8.DecodeRuneInString(childSuffix)
	splitNode.children[childFirstRune] = child
	cur.children[firstRune] = splitNode

	if keySuffix == "" {
		splitNode.data = data
		splitNode.hasData = true
	} else {
		keyFirstRune, _ := utf8.DecodeRuneInString(keySuffix)
		splitNode.children[keyFirstRune] = &RadixNode[T]{
			label:   keySuffix,
			data:    data,
			hasData: true,
		}
	}
	return nil
}

func (t *RadixTree[T]) Upsert(key string, initial T, update func(T) T) {
	if t.root == nil {
		t.root = &RadixNode[T]{}
	}
	radixUpsert(t.root, key, initial, update)
}

func radixUpsert[T any](cur *RadixNode[T], remaining string, initial T, update func(T) T) {
	if remaining == "" {
		if cur.hasData {
			cur.data = update(cur.data)
		} else {
			cur.data = initial
			cur.hasData = true
		}
		return
	}

	firstRune, _ := utf8.DecodeRuneInString(remaining)
	if cur.children == nil {
		cur.children = make(map[rune]*RadixNode[T])
	}

	child, ok := cur.children[firstRune]
	if !ok {
		cur.children[firstRune] = &RadixNode[T]{label: remaining, data: initial, hasData: true}
		return
	}

	cpLen := commonPrefixBytes(remaining, child.label)
	if cpLen == len(child.label) {
		radixUpsert(child, remaining[cpLen:], initial, update)
		return
	}

	// Split the child edge at cpLen
	commonLabel := remaining[:cpLen]
	childSuffix := child.label[cpLen:]
	keySuffix := remaining[cpLen:]

	splitNode := &RadixNode[T]{
		label:    commonLabel,
		children: make(map[rune]*RadixNode[T]),
	}
	child.label = childSuffix
	childFirstRune, _ := utf8.DecodeRuneInString(childSuffix)
	splitNode.children[childFirstRune] = child
	cur.children[firstRune] = splitNode

	if keySuffix == "" {
		splitNode.data = initial
		splitNode.hasData = true
	} else {
		keyFirstRune, _ := utf8.DecodeRuneInString(keySuffix)
		splitNode.children[keyFirstRune] = &RadixNode[T]{
			label:   keySuffix,
			data:    initial,
			hasData: true,
		}
	}
}

func (t *RadixTree[T]) Search(key string) (T, error) {
	var zero T
	if t == nil || t.root == nil {
		return zero, nil
	}
	return radixSearch(t.root, key)
}

func radixSearch[T any](cur *RadixNode[T], remaining string) (T, error) {
	var zero T
	if remaining == "" {
		if cur.hasData {
			return cur.data, nil
		}
		return zero, nil
	}

	firstRune, _ := utf8.DecodeRuneInString(remaining)
	if cur.children == nil {
		return zero, nil
	}
	child, ok := cur.children[firstRune]
	if !ok {
		return zero, nil
	}

	cpLen := commonPrefixBytes(remaining, child.label)
	if cpLen < len(child.label) {
		return zero, nil
	}
	return radixSearch(child, remaining[cpLen:])
}

// Size returns an estimate of the memory consumed by the radix tree in bytes.
func (t *RadixTree[T]) Size() uintptr {
	if t == nil || t.root == nil {
		return unsafe.Sizeof(*t)
	}
	return unsafe.Sizeof(*t) + radixNodeSize(t.root)
}

func radixNodeSize[T any](n *RadixNode[T]) uintptr {
	if n == nil {
		return 0
	}
	size := unsafe.Sizeof(*n)
	size += uintptr(len(n.label))

	v := reflect.ValueOf(n.data)
	if v.Kind() == reflect.String {
		size += uintptr(v.Len())
	}

	if len(n.children) > 0 {
		const mapHeaderBytes = 128
		const mapEntryBytes = 12
		size += mapHeaderBytes + uintptr(len(n.children))*mapEntryBytes
	}

	for _, child := range n.children {
		size += radixNodeSize(child)
	}
	return size
}
