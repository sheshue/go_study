package main

import (
	"fmt"
)

func classChecker(items ...interface{}) {
	for i, x := range items {
		switch x := x.(type) {
		case bool:
			fmt.Printf("prame #%d is a bool, value: %t\n", i, x)
		case float64:
			fmt.Printf("param #%d is a float64, value: %f\n", i, x)
		case int, int8, int16, int32, int64:
			fmt.Printf("param #%d is a int, value: %d\n", i, x)
		case uint, uint8, uint32, uint64:
			fmt.Printf("param #%d is a unit, value: %d\n", i, x)
		case nil:
			fmt.Printf("param #%d is a nil\n", i)
		case string:
			fmt.Printf("param #%d is a string, value: %s\n", i, x)
		default:
			fmt.Printf("param #%d's type is unknow\n", i)
		}
	}
}

func main() {
	classChecker(5, -17.98, "aiden", nil, true, complex(1, 1))
}
