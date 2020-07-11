package main

import (
	"fmt"
	"sync"
)

// Parallelsum does parallel vector sum,
// in each loop,buffered input capacity will be cut half
// for the next loop to Sum goroutine to consume.
// it terminates untill the input remains one.
func ParallelSum(slcs ...[]int) []int {
	input := make(chan []int, len(slcs))
	output := make(chan []int)
	var result []int
	go func(input chan []int) {
		for _, slc := range slcs {
			input <- slc
		}
		close(input)
	}(input)

	for forLoop := 0; ; forLoop++ {
		var wg sync.WaitGroup
		wg.Add(cap(input) / 2)
		for i := 0; i < cap(input)/2; i++ {
			out := Sum(input, i, forLoop)
			go func() {
				defer wg.Done()
				for o := range out {
					output <- o
				}
			}()
		}
		go func(output chan []int) {
			wg.Wait()
			close(output)
		}(output)
		input = make(chan []int, cap(input)/2)
		if cap(input) < 2 {
			result = <-output
			break
		}
		for o := range output {
			input <- o
		}
		output = make(chan []int)
		close(input)
	}
	return result
}

//TODO:complete the Sum for the parallel sum function.
// ParallelSum的策略是每次取一半的Goroutine，每个goroutine都自由地从input中拿vector出来求和。
// 并且把结果加到新的input中，直到老input被取光。之后开始下一个循环
func Sum(in chan []int, iter int, forLoop int) (out chan []int) {
	out = make(chan []int)
	go func() {
		acc := <-in
		cnt := 1
		for vec := range in {
			if len(vec) == 0 {
				break
			}
			for i, v := range vec {
				acc[i] += v
			}
			cnt++
		}
		fmt.Printf("==forLoop %d iter %d consume: %d; output: %v==\n", forLoop, iter, cnt, acc)
		out <- acc
		defer close(out)
	}()
	return out
}

func main() {
	fmt.Println(`Please edit main.go,and complete the 'Sum' function for the parallel sum to pass the test.
Concurrency is the most important feature of Go,and the principle is
'Do not communicate by sharing memory; instead, share memory by communicating.'
In this exercise you need to catch many features of channels.This is a tour for you to figure out!
Because here the focus is pipleline model (link:http://blog.golang.org/pipelines).
It's different from the custom parallel vector sum in which sum numer at every index of the vectors in a goroutine.
In this exercies,vector is just a abstract,you can change it to a struct or any thing else that can be sumed up.`)
}
