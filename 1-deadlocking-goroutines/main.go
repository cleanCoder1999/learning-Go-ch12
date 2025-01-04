package main

import "fmt"

// main consists two examples that show the difference between goroutines,
// ones that are deadlocking
// and use the select statement to avoid deadlocks, respectively
func main() {
	// ### - deadlocking goroutines
	{
		ch1 := make(chan int)
		ch2 := make(chan int)
		go func() {
			inGoroutine := 1
			// write pauses goroutine until value is read from channel
			ch1 <- inGoroutine // writing to channel
			fromMain := <-ch2  // reading from channel
			fmt.Println("goroutine:", inGoroutine, fromMain)
		}()

		inMain := 2
		// write pauses goroutine until value is read from channel
		ch2 <- inMain          // writing to channel
		fromGoroutine := <-ch1 // reading from channel
		fmt.Println("main:", inMain, fromGoroutine)
	}

	// ### - using select to avoid deadlocking
	//       (not an ideal solution: the main goroutine exits before the explicitly launched goroutine finishes => leaking)
	{
		ch1 := make(chan int)
		ch2 := make(chan int)
		go func() {
			inGoroutine := 1
			// write pauses goroutine until value is read from channel
			ch1 <- inGoroutine // writing to channel
			fromMain := <-ch2  // reading from channel
			fmt.Println("goroutine:", inGoroutine, fromMain)
		}()

		var fromGoroutine int
		inMain := 2

		// #select statement
		// each case in a select statement is a read or write to a channel
		//
		// if a read or write is possible for a case,
		// it is executed along with the body of the case.
		//
		// if multiple cases have channels that can be read or written,
		// the select algorithm picks randomly from any of its cases that can go forward
		select {
		case v := <-ch1: // reading from channel
			fromGoroutine = v
		case ch2 <- inMain: // writing to channel
		}
		fmt.Println("main:", inMain, fromGoroutine)
	}

}
