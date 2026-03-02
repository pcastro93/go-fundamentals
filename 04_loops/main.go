package main

import "fmt"

func letter_counter() {
	s := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean
	convallis, nisi vel cursus sodales, nisi urna feugiat ligula, quis accumsan
 erat ligula at justo. Nullam at ligula sit amet nisl ornare ornare vel sed
 tellus. Curabitur sed risus neque. Fusce at sem quis velit vulputate feugiat.
 Vivamus id elit magna. Maecenas in efficitur dui. Curabitur sem est, malesuada
 at odio a, gravida tincidunt dui. Nullam ac commodo elit`
	counters := make(map[rune]int, 30)
	for _, letter := range s {
		if letter == '\n' {
			continue
		}
		_, ok := counters[letter]
		if ok {
			counters[letter]++
		} else {
			counters[letter] = 1
		}
	}
	for letter, count := range counters {
		fmt.Printf("letter=%c:counter=%d\n", letter, count)
	}
}

func main() {
	letter_counter()
}
