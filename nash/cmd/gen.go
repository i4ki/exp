package main

import "fmt"

func main() {
	for i := 'a'; i < 'z'; i += 2 {
		fmt.Printf("%c%c\n", rune(i), rune(i+1))
	}
}
