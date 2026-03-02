package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func indexing() {
	toPrint := "It should support eñes and emojis 👍"
	fmt.Println("Total bytes", len(toPrint))
	fmt.Println("Total runes", utf8.RuneCountInString(toPrint))

	fmt.Println("==========")
	fmt.Println("Bytes")
	fmt.Println("==========")
	for i := range len(toPrint) {
		fmt.Printf("%02X ", toPrint[i])
	}
	fmt.Println()
	// Print each rune
	fmt.Println("==========")
	fmt.Println("Runes 1")
	fmt.Println("==========")
	for _, c := range toPrint {
		fmt.Print(string(c))
	}
	fmt.Println()
	// Convert to slice of runes
	fmt.Println("==========")
	fmt.Println("Runes 2")
	fmt.Println("==========")
	runes := []rune(toPrint)
	for _, ri := range runes {
		fmt.Print(string(ri))
	}
	fmt.Println()
}

func reverse() {
	toReverse := `A leader’s job is to have all the questions.
	You have to be incredibly comfortable looking like the dumbest person
	in the room.`
	runes := []rune(toReverse)
	slices.Reverse(runes)
	fmt.Println("Normal:", toReverse)
	fmt.Println("Reversed:", string(runes))
}

func palindromeCheck() {
	toCheck := []struct {
		s        string
		expected bool
	}{
		{"abcdcba", true},
		{"abc", false},
		{"a", true},
		{"ab", false},
	}

	check := func(s string) bool {
		runes := []rune(s)
		for i := 0; i < len(runes)/2; i++ {
			if runes[i] != runes[len(runes)-i-1] {
				return false
			}
		}
		return true
	}
	for _, tci := range toCheck {
		if res := check(tci.s); res != tci.expected {
			panic(fmt.Sprintf("Palindrome validation failed: %s got %t expected %t", tci.s, res, tci.expected))
		}
	}
}

func wordCount() {
	s := "word1 2 word3 hello world"
	fields := strings.Fields(s)
	fmt.Println("Total words", len(fields))
}

func anagramCheck() {
	example1 := "abc"
	example2 := "bca"
	getFreq := func(s string) map[rune]int {
		freqMap := map[rune]int{}
		for _, r := range s {
			freqMap[r]++
		}
		return freqMap
	}
	check := func(s1 string, s2 string) bool {
		m1 := getFreq(s1)
		m2 := getFreq(s2)
		for k, v := range m1 {
			if (m2)[k] != v {
				return false
			}
		}
		return true
	}
	fmt.Println("Is anagram?", check(example1, example2))
}

func stringBuilder() {
	iterations := 50_000
	// string
	strStart := time.Now()
	s := ""
	for i := range iterations {
		s += strconv.Itoa(i)
	}
	strDuration := time.Since(strStart).Milliseconds()

	fmt.Printf("Time using string: %d (ms)\n", strDuration)
	// Builder
	builderStart := time.Now()
	var sb strings.Builder
	for i := range iterations {
		sb.WriteString(strconv.Itoa(i))
	}
	builderDuration := time.Since(builderStart).Milliseconds()
	fmt.Printf("Time using builder: %d (ms)\n", builderDuration)
}

func splitting() {
	s := "  /home/user/hello/world"
	from := strings.Index(s, "/") + 1
	splitted := strings.SplitSeq(s[from:], "/")
	for si := range splitted {
		fmt.Println(si)
	}
}

func main() {
	indexing()
	reverse()
	palindromeCheck()
	wordCount()
	anagramCheck()
	stringBuilder()
	splitting()
}
