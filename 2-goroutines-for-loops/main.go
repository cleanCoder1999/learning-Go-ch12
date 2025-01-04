package main

import "fmt"

func main() {
	// ### - before Go 1.22, each iteration of a for-loop used the same index and value variable
	//       which caused issues when working with closures
	{
		a := []int{2, 4, 6, 8, 10}
		ch := make(chan int, len(a))

		// when running with 1.21 or older, Go uses the same variable for v in each loop iteration
		// when running with 1.22 or later, Go creates a new variable for v in each loop iteration
		// to check the difference, change Go version in the go.mod file
		//
		// when using variables within a closure that were captured, this is difference is particularly important
		//
		// to prevent this being an issue regardless of the Go version in use,
		// (1) pass variables that are used within a closure as parameters (Go automatically makes a copy since it is a "call by value" language)
		// (2) OR make sure to make a shadowed copy of the captured value before the closure is entered
		for _, v := range a {
			go func() {
				fmt.Println("memory address of v:", &v)
				ch <- v * 2
			}()
		}

		for i := 0; i < len(a); i++ {
			fmt.Println(<-ch)
		}
	}
}
