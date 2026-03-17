package main

import "testing"

func TestRadixInsertAndSearch(t *testing.T) {
	rt := NewRadixTree[int]()
	if err := rt.Insert("hello", 42); err != nil {
		t.Fatalf("unexpected error on insert: %v", err)
	}
	val, err := rt.Search("hello")
	if err != nil {
		t.Fatalf("unexpected error on search: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %d", val)
	}
}

func TestRadixSearchPrefixNotInserted(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("apple", "fruit")

	val, _ := rt.Search("app")
	if val != "" {
		t.Errorf("expected empty string for prefix-only key, got %q", val)
	}
}

func TestRadixInsertDuplicateKeyReturnsError(t *testing.T) {
	rt := NewRadixTree[int]()
	if err := rt.Insert("dup", 1); err != nil {
		t.Fatalf("first insert should succeed, got: %v", err)
	}
	if err := rt.Insert("dup", 2); err == nil {
		t.Error("second insert of duplicate key should return an error")
	}
}

func TestRadixSearchOnEmpty(t *testing.T) {
	rt := NewRadixTree[int]()
	val, _ := rt.Search("anything")
	if val != 0 {
		t.Errorf("expected zero value on empty tree, got %d", val)
	}
}

func TestRadixSharedPrefixKeysAreIndependent(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("car", "vehicle")
	rt.Insert("card", "plastic")

	car, _ := rt.Search("car")
	card, _ := rt.Search("card")

	if car != "vehicle" {
		t.Errorf("expected 'vehicle' for 'car', got %q", car)
	}
	if card != "plastic" {
		t.Errorf("expected 'plastic' for 'card', got %q", card)
	}
}

func TestRadixUpsertNewKeyStoresInitialValue(t *testing.T) {
	rt := NewRadixTree[int]()
	rt.Upsert("word", 1, func(n int) int { return n + 1 })

	val, _ := rt.Search("word")
	if val != 1 {
		t.Errorf("expected 1, got %d", val)
	}
}

func TestRadixUpsertExistingKeyAppliesUpdateFn(t *testing.T) {
	rt := NewRadixTree[int]()
	rt.Upsert("word", 1, func(n int) int { return n + 1 })
	rt.Upsert("word", 1, func(n int) int { return n + 1 })

	val, _ := rt.Search("word")
	if val != 2 {
		t.Errorf("expected 2, got %d", val)
	}
}

func TestRadixUpsertAccumulatesCorrectly(t *testing.T) {
	rt := NewRadixTree[int]()
	inc := func(n int) int { return n + 1 }
	for range 5 {
		rt.Upsert("word", 1, inc)
	}

	val, _ := rt.Search("word")
	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}
}

func TestRadixUpsertIndependentKeys(t *testing.T) {
	rt := NewRadixTree[int]()
	inc := func(n int) int { return n + 1 }

	for range 3 {
		rt.Upsert("foo", 1, inc)
	}
	for range 7 {
		rt.Upsert("bar", 1, inc)
	}

	foo, _ := rt.Search("foo")
	bar, _ := rt.Search("bar")
	if foo != 3 {
		t.Errorf("expected foo=3, got %d", foo)
	}
	if bar != 7 {
		t.Errorf("expected bar=7, got %d", bar)
	}
}

func TestRadixUpsertSharedPrefixCountedIndependently(t *testing.T) {
	rt := NewRadixTree[int]()
	inc := func(n int) int { return n + 1 }

	for range 2 {
		rt.Upsert("car", 1, inc)
	}
	for range 4 {
		rt.Upsert("card", 1, inc)
	}

	car, _ := rt.Search("car")
	card, _ := rt.Search("card")
	if car != 2 {
		t.Errorf("expected car=2, got %d", car)
	}
	if card != 4 {
		t.Errorf("expected card=4, got %d", card)
	}
}

func TestRadixUnicodeKey(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("café", "drink")

	val, _ := rt.Search("café")
	if val != "drink" {
		t.Errorf("expected 'drink' for unicode key, got %q", val)
	}

	val, _ = rt.Search("cafe")
	if val != "" {
		t.Errorf("expected empty for ASCII variant of unicode key, got %q", val)
	}
}

// --- Radix-specific: node splitting ---

// TestRadixSplitMidEdge inserts two keys that share a common prefix shorter than
// either key, forcing a split in the middle of an existing edge.
func TestRadixSplitMidEdge(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("apple", "fruit")
	rt.Insert("apply", "verb")

	apple, _ := rt.Search("apple")
	apply, _ := rt.Search("apply")

	if apple != "fruit" {
		t.Errorf("expected 'fruit' for 'apple', got %q", apple)
	}
	if apply != "verb" {
		t.Errorf("expected 'verb' for 'apply', got %q", apply)
	}
}

// TestRadixShortKeyAfterLong inserts a longer key first, then a key that is
// a prefix of the longer one, triggering a split where the shorter key lands
// on the newly created intermediate node.
func TestRadixShortKeyAfterLong(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("testing", "gerund")
	rt.Insert("test", "noun")

	test, _ := rt.Search("test")
	testing_, _ := rt.Search("testing")

	if test != "noun" {
		t.Errorf("expected 'noun' for 'test', got %q", test)
	}
	if testing_ != "gerund" {
		t.Errorf("expected 'gerund' for 'testing', got %q", testing_)
	}
}

// TestRadixLongKeyAfterShort inserts the prefix key first, then the longer key,
// which should be stored as a child without any split.
func TestRadixLongKeyAfterShort(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("test", "noun")
	rt.Insert("testing", "gerund")

	test, _ := rt.Search("test")
	testing_, _ := rt.Search("testing")

	if test != "noun" {
		t.Errorf("expected 'noun' for 'test', got %q", test)
	}
	if testing_ != "gerund" {
		t.Errorf("expected 'gerund' for 'testing', got %q", testing_)
	}
}

// TestRadixMultipleSplits inserts three keys that cause two successive splits,
// verifying all values remain retrievable after the tree restructures.
func TestRadixMultipleSplits(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("apple", "fruit")
	rt.Insert("apply", "verb")
	rt.Insert("application", "program")

	cases := map[string]string{
		"apple":       "fruit",
		"apply":       "verb",
		"application": "program",
	}
	for key, want := range cases {
		got, _ := rt.Search(key)
		if got != want {
			t.Errorf("Search(%q): expected %q, got %q", key, want, got)
		}
	}
}

// TestRadixPrefixOfInsertedKeyNotFound verifies that a prefix of an inserted key
// returns the zero value, even after a split restructured the tree.
func TestRadixPrefixOfInsertedKeyNotFound(t *testing.T) {
	rt := NewRadixTree[string]()
	rt.Insert("apple", "fruit")
	rt.Insert("apply", "verb")

	// "appl" is the common split node but was never explicitly inserted
	val, _ := rt.Search("appl")
	if val != "" {
		t.Errorf("expected empty for un-inserted split prefix 'appl', got %q", val)
	}
}
