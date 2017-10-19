package main

import (
	"fmt"
	"sort"
)

type ByLenght []string

func (s ByLenght) Len() int {
	return len(s)
}
func (s ByLenght) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByLenght) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

func main() {
	fruits := []string{"peach", "banana", "kivi"}
	sort.Sort(ByLenght(fruits))
	fmt.Println((fruits))
}
