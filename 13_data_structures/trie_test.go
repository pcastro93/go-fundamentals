package main

import "testing"

func TestInsertAndSearch(t *testing.T) {
	trie := NewTrie[int]()
	if err := trie.Insert("hello", 42); err != nil {
		t.Fatalf("unexpected error on insert: %v", err)
	}
	val, err := trie.Search("hello")
	if err != nil {
		t.Fatalf("unexpected error on search: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

func TestSearchPrefixNotInserted(t *testing.T) {
	trie := NewTrie[string]()
	trie.Insert("apple", "fruit")

	val, _ := trie.Search("app")
	if val != "" {
		t.Errorf("expected empty string for prefix-only key, got %q", val)
	}
}

func TestInsertDuplicateKeyReturnsError(t *testing.T) {
	trie := NewTrie[int]()
	if err := trie.Insert("dup", 1); err != nil {
		t.Fatalf("first insert should succeed, got: %v", err)
	}
	if err := trie.Insert("dup", 2); err == nil {
		t.Error("second insert of duplicate key should return an error")
	}
}

func TestSearchOnEmptyTrie(t *testing.T) {
	trie := NewTrie[int]()
	val, _ := trie.Search("anything")
	if val != 0 {
		t.Errorf("expected zero value on empty trie, got %d", val)
	}
}

func TestSharedPrefixKeysAreIndependent(t *testing.T) {
	trie := NewTrie[string]()
	trie.Insert("car", "vehicle")
	trie.Insert("card", "plastic")

	car, _ := trie.Search("car")
	card, _ := trie.Search("card")

	if car != "vehicle" {
		t.Errorf("expected 'vehicle' for 'car', got %q", car)
	}
	if card != "plastic" {
		t.Errorf("expected 'plastic' for 'card', got %q", card)
	}
}

func TestUpsertNewKeyStoresInitialValue(t *testing.T) {
	trie := NewTrie[int]()
	trie.Upsert("word", 1, func(n int) int { return n + 1 })

	val, _ := trie.Search("word")
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestUpsertExistingKeyAppliesUpdateFn(t *testing.T) {
	trie := NewTrie[int]()
	trie.Upsert("word", 1, func(n int) int { return n + 1 })
	trie.Upsert("word", 1, func(n int) int { return n + 1 })

	val, _ := trie.Search("word")
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
}

func TestUpsertAccumulatesCorrectly(t *testing.T) {
	trie := NewTrie[int]()
	inc := func(n int) int { return n + 1 }
	for range 5 {
		trie.Upsert("word", 1, inc)
	}

	val, _ := trie.Search("word")
	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}
}

func TestUpsertIndependentKeys(t *testing.T) {
	trie := NewTrie[int]()
	inc := func(n int) int { return n + 1 }

	for range 3 {
		trie.Upsert("foo", 1, inc)
	}
	for range 7 {
		trie.Upsert("bar", 1, inc)
	}

	foo, _ := trie.Search("foo")
	bar, _ := trie.Search("bar")
	if foo != 3 {
		t.Errorf("expected foo=3, got %d", foo)
	}
	if bar != 7 {
		t.Errorf("expected bar=7, got %d", bar)
	}
}

func TestUpsertSharedPrefixCountedIndependently(t *testing.T) {
	trie := NewTrie[int]()
	inc := func(n int) int { return n + 1 }

	for range 2 {
		trie.Upsert("car", 1, inc)
	}
	for range 4 {
		trie.Upsert("card", 1, inc)
	}

	car, _ := trie.Search("car")
	card, _ := trie.Search("card")
	if car != 2 {
		t.Errorf("expected car=2, got %d", car)
	}
	if card != 4 {
		t.Errorf("expected card=4, got %d", card)
	}
}

func TestUnicodeKey(t *testing.T) {
	trie := NewTrie[string]()
	trie.Insert("café", "drink")

	val, _ := trie.Search("café")
	if val != "drink" {
		t.Errorf("expected 'drink' for unicode key, got %q", val)
	}

	val, _ = trie.Search("cafe")
	if val != "" {
		t.Errorf("expected empty for ASCII variant of unicode key, got %q", val)
	}
}
