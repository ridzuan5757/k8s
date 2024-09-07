package main

import (
	"fmt"
	"time"
)

// When a program starts, its only goroutine is the one that calls the `main`
// function, so we call it the main goroutine.
// New goroutine are created by the `go` statement. Syntactically, a go
// statement is an ordinary function or method call prefixed by the keyword
// `go`.
// A `go` statement causes the function to be called in a newly created
// goroutine. The `go` statement itself complates immediately:
//
//    f()     // call f(); wait for it to return
//    go f()  // create a new goroutine that calls f(); do not wait
//
// After several animation, the `fib(45)` call returns and the main function
// prints its result.
//
// The `main` function then returns. When this happens, are abruptly terminated
// and the program exits. Other than by returning from the main or exiting the
// program, there is no programmatic way for one goroutine to stop another.

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibn := fib(n)
	fmt.Printf("\r %d", fibn)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	} else {
		return fib(x-1) + fib(x-2)
	}
}
