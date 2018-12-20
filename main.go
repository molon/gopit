package main

import "fmt"

var a = 100

func init() {
	// fmt.Println(a)
}

func main() {
	path := []string{"a", "b", "c", "d"}
	fType := "c"

fields:
	for i := 0; i < 10; i++ {
		fmt.Println(i)
		for j, e := range path {
			fmt.Println("   ", j, " ", e)
			if fType == e {
				fmt.Println("continue fields")
				continue fields
			}
		}
	}
}
