package main

import (
	"fmt"

	"github.com/samarth8765/bloom-filters/bloom"
)

func main() {
	filter := bloom.NewBloomFilter(10, 2)
	arr := [...]string{"apple", "oranges", "bananas", "pears", "grapes", "mangoes"}

	for _, fruit := range arr {
		filter.Add([]byte(fruit))
	}

	testArray := [...]string{"apple", "oranges", "bananas", "pears", "grapes", "mangoes",
		"watermelon", "strawberries", "pineapples", "kiwi"}

	for _, fruit := range testArray {
		found, hashIdxs, err := filter.Check([]byte(fruit))
		if err != nil {
			panic(err)
		}
		if found {
			fmt.Printf("%s found at indexes %v\n", fruit, hashIdxs)
		} else {
			fmt.Printf("%s not found\n", fruit)
		}
	}

}
