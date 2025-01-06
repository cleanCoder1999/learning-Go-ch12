package main

import (
	"fmt"
	"math"
	"sync"
)

func main() {
	// ### - exercise 1.1: function creates 2 goroutines that write strings to a channel
	{
		ch := make(chan string, 20)

		writeStringsToChan(ch)

		for i := 0; i < 20; i++ {
			fmt.Println(<-ch)
		}
		close(ch)
		fmt.Println("")
	}

	// ### - exercise 1.2: function creates 3 goroutines; 2 write 10 numbers each; 1 reads all numbers and prints them out
	{
		var wg sync.WaitGroup
		wg.Add(1)

		handleGoRoutinesWithWaitGroups(&wg)
		wg.Wait()
		fmt.Println("")
	}

	// ### - exercise 1.3 (SOLUTION): function creates 3 goroutines; 2 write 10 numbers each; 1 reads all numbers and prints them out
	{
		ProcessDataExercise1()
		fmt.Println("")
	}

	// ### - exercise 2: function creates 2 goroutines; each writes 10 numbers in its channel; use a for-select statement to read from both channels
	{
		ProcessDataExercise2()
		fmt.Println("")
	}

	// ### - exercise 3: sync.OnceValue
	{
		for i := 0; i < 100_000; i++ {
			if i%1000 == 0 {
				sqrt := lookup(i)
				fmt.Printf("i=%d: %f\n", i, sqrt)
			}
		}
		fmt.Println("")
	}
}

func lookup(i int) float64 {
	// NOTE:
	//	due to sync.OnceValue() its parameter function is invoked only once;
	// 	after the first invocation, buildMapCached() returns the cached value that was saved/cached after the first invocation
	return buildMapCached()[i]
}

// OnceValue returns a function that invokes f only once
//
//	and returns the value returned by f
//
// NOTE: the value returned by f is cached for future invocations
var buildMapCached func() map[int]float64 = sync.OnceValue(buildMap)

func buildMap() map[int]float64 {
	fmt.Println("in buildMap()")

	m := make(map[int]float64, 100_000)

	for i := 0; i < 100_000; i++ {
		m[i] = math.Sqrt(float64(i))
	}

	return m
}

func ProcessDataExercise2() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for count := 0; count < 2; {
		select {
		case v, ok := <-ch1:
			if !ok {
				ch1 = nil // turn off case <-ch1
				count++
				continue
			}
			fmt.Println("g1:", v)
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil // turn off case <-ch2
				count++
				continue
			}
			fmt.Println("g2:", v)
		}
	}
}

func writeStringsToChan(ch chan string) {

	go func() {
		fmt.Println("#1 started")
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("go1: %d", i)
		}
	}()

	go func() {
		fmt.Println("#2 started")
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("go2: %d", i)
		}
	}()
}

func handleGoRoutinesWithWaitGroups(wg *sync.WaitGroup) {

	ch := make(chan string)
	done1 := make(chan struct{})
	done2 := make(chan struct{})

	// write 10 numbers
	go func() {
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("go1: %d", i)
		}
		close(done1)
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ch <- fmt.Sprintf("go2: %d", i)
		}
		close(done2)
	}()

	// read all numbers
	go func() {

		i := 0
		for count := 0; count < 2; {

			select {
			case v := <-ch:
				fmt.Printf("%d: %s\n", i, v)
				i++
			case _, ok := <-done1:
				if !ok {
					done1 = nil // turn off the case <-done1
					count++
					continue
				}
			case _, ok := <-done2:
				if !ok {
					done2 = nil // turn off the case <-done2
					count++
					continue
				}
			}
		}

		close(ch)
		wg.Done()
	}()
}

func ProcessDataExercise1() {
	ch := make(chan int)
	// use 2 waitgroups!
	// the 1st waitgroup controls when to close the channel
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- i*100 + 1
		}
	}()

	// launch this helper goroutine to close the channel when the two writing goroutines are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	// the second waitgroup signals when the reading goroutine is done
	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}()
	wg2.Wait()
}
