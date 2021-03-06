package main

import (
	"fmt"
	"strconv"
	"sync"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	v   map[int]bool
	mux sync.Mutex
}

func main() {
	var wg sync.WaitGroup

	c := SafeCounter{v: make(map[int]bool)}

	lowerBound := 264360
	upperBound := 746325

	//RULE 2: The value is within the range given in your puzzle input.
	for i := lowerBound; i <= upperBound; i++ {
		wg.Add(1)
		go validateNumber(&wg, &c, i)

	}

	// Waiti for workers to finish
	wg.Wait()
	fmt.Printf("Number of passwords found: %d", len(c.v))
}

func validateNumber(wg *sync.WaitGroup, c *SafeCounter, id int) {
	defer wg.Done()

	s := strconv.Itoa(id)

	// Rule 1: It is a six-digit number.
	if len(s) != 6 {
		return
	}

	previousDigit := 0
	doubleFound := false
	for _, v := range s {
		currentDigit := int(v - 48)

		// Rule 4: Going from left to right, the digits never decrease; they only ever increase or stay the same (like 111123 or 135679).
		if currentDigit < previousDigit {
			return
		}

		// Rule 3: Two adjacent digits are the same (like 22 in 122345).
		if previousDigit == currentDigit {
			doubleFound = true
		}

		previousDigit = currentDigit
	}

	// Rule 3: Two adjacent digits are the same (like 22 in 122345).
	if !doubleFound {
		return
	}

	// add it to the list of found numbers
	c.Add(id)
}

// Safely ddd the number to our result set.
func (c *SafeCounter) Add(key int) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.v.
	c.v[key] = true
	c.mux.Unlock()
}
