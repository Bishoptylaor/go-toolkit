package main

import (
	"fmt"
	"github.com/Bishoptylaor/go-toolkit/xslice"
)

func main() {
	// zsets
	a := []int{1, 2, 3, 4, 5}
	b := []int{4, 5, 6, 7, 8}

	zs := xslice.NewZSets[int](a, b)
	fmt.Println("Equal: ", zs.Equal())                 // false
	fmt.Println("AExceptB: ", zs.AExceptB())           // [1 2 3]
	fmt.Println("BExceptA: ", zs.BExceptA())           // [6 7 8]
	fmt.Println("Intersection: ", zs.Intersection())   // [4 5]
	fmt.Println("Union: ", zs.Union())                 // [1 2 3 4 5 4 5 6 7 8]
	fmt.Println("DistinctUnion: ", zs.DistinctUnion()) // [1 2 3 4 5 6 7 8]

	// zslice
	c := []int{1, 2, 3, 4, 5}
	d := []int{1, 2, 3, 3, 4}
	xslice.Reverse(c)
	fmt.Println("Reverse a: ", c)
	fmt.Println("Distinct: ", xslice.Distinct(d))
	fmt.Println("Contains: ", xslice.Contains(d, 5))
	fmt.Println("Cut: ", xslice.Cut(d, 5, 10))
}
