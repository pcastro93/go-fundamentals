package main

import (
	"os"
	"strings"
	"testing"
)

// ----------------------------------------------------------------------------
// Shared corpus — built once for all benchmarks via TestMain.
// ----------------------------------------------------------------------------

var (
	corpus     []string   // flat list of words (with repetitions, like real text)
	benchTrie  *Trie[int] // pre-built for search benchmarks
	benchRadix *RadixTree[int]
)

// TestMain runs once per test binary. It builds the corpus and the two
// pre-loaded structures used by search benchmarks, then hands control to
// m.Run() which executes all tests and benchmarks.
func TestMain(m *testing.M) {
	corpus = buildCorpus()

	inc := func(n int) int { return n + 1 }

	benchTrie = NewTrie[int]()
	for _, w := range corpus {
		benchTrie.Upsert(w, 1, inc)
	}

	benchRadix = NewRadixTree[int]()
	for _, w := range corpus {
		benchRadix.Upsert(w, 1, inc)
	}

	os.Exit(m.Run())
}

// buildCorpus produces a synthetic word list that mimics natural-language text:
//   - prefix × root × suffix combinations create thousands of unique words with
//     rich shared-prefix structure, which stresses both data structures differently.
//   - Each unique word is repeated 10 times to simulate the Upsert
//     frequency-counting pattern from main.go.
func buildCorpus() []string {
	prefixes := []string{
		"", "un", "re", "pre", "over", "under", "out", "up",
	}
	roots := []string{
		"do", "make", "take", "come", "go", "see", "get", "give",
		"know", "think", "say", "tell", "hear", "find", "put",
		"stand", "keep", "let", "begin", "show", "feel", "turn",
		"call", "try", "ask", "need", "seem", "leave", "move",
		"live", "bring", "write", "provide", "sit", "lose", "pay",
		"meet", "include", "continue", "set", "learn", "change",
		"lead", "watch", "follow", "stop", "create", "speak",
		"read", "spend", "grow", "open", "walk", "win", "offer",
		"remember", "love", "consider", "appear", "buy", "wait",
		"serve", "die", "send", "expect", "build", "stay", "fall",
		"cut", "reach", "remain", "suggest", "raise", "pass",
		"sell", "require", "report", "decide", "pull",
	}
	suffixes := []string{
		"", "s", "ed", "ing", "er", "tion", "ment", "ness", "able",
	}

	seen := map[string]bool{}
	var unique []string
	for _, p := range prefixes {
		for _, r := range roots {
			for _, s := range suffixes {
				w := p + r + s
				if !seen[w] {
					seen[w] = true
					unique = append(unique, w)
				}
			}
		}
	}

	// 10 repetitions per unique word ≈ 50 k words total.
	out := make([]string, 0, len(unique)*10)
	for _, w := range unique {
		for range 10 {
			out = append(out, w)
		}
	}
	return out
}

// sinkInt prevents the compiler from dead-code-eliminating search calls.
var sinkInt int

// ----------------------------------------------------------------------------
// Build benchmarks — measure the cost of constructing each structure from
// scratch. b.ReportAllocs() enables the B/op and allocs/op columns.
// b.ReportMetric reports the final tree size as a custom "bytes/tree" column
// so you can compare memory footprint in the same run.
// ----------------------------------------------------------------------------

func BenchmarkTrieBuild(b *testing.B) {
	inc := func(n int) int { return n + 1 }
	b.ReportAllocs()
	var trie *Trie[int]
	for b.Loop() {
		trie = NewTrie[int]()
		for _, w := range corpus {
			trie.Upsert(w, 1, inc)
		}
	}
	b.ReportMetric(float64(trie.Size()), "bytes/tree")
}

func BenchmarkRadixBuild(b *testing.B) {
	inc := func(n int) int { return n + 1 }
	b.ReportAllocs()
	var rt *RadixTree[int]
	for range b.N {
		rt = NewRadixTree[int]()
		for _, w := range corpus {
			rt.Upsert(w, 1, inc)
		}
	}
	b.ReportMetric(float64(rt.Size()), "bytes/tree")
}

// ----------------------------------------------------------------------------
// Search benchmarks — measure lookup cost on a pre-built structure so build
// time does not contaminate the result.
// ----------------------------------------------------------------------------

// hit: the key exists
func BenchmarkTrieSearchHit(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		v, _ := benchTrie.Search("reading")
		sinkInt = v
	}
}

func BenchmarkRadixSearchHit(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		v, _ := benchRadix.Search("reading")
		sinkInt = v
	}
}

// miss: the key is absent — tests early-exit behaviour
func BenchmarkTrieSearchMiss(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		v, _ := benchTrie.Search("zzzmissing")
		sinkInt = v
	}
}

func BenchmarkRadixSearchMiss(b *testing.B) {
	b.ReportAllocs()
	for range b.N {
		v, _ := benchRadix.Search("zzzmissing")
		sinkInt = v
	}
}

// long key: exercises deeper traversal
func BenchmarkTrieSearchLong(b *testing.B) {
	b.ReportAllocs()
	key := strings.Repeat("re", 6) + "reading" // deep path
	for range b.N {
		v, _ := benchTrie.Search(key)
		sinkInt = v
	}
}

func BenchmarkRadixSearchLong(b *testing.B) {
	b.ReportAllocs()
	key := strings.Repeat("re", 6) + "reading"
	for range b.N {
		v, _ := benchRadix.Search(key)
		sinkInt = v
	}
}
