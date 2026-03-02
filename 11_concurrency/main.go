package main

import (
	"fmt"
	"sync"
	"time"
)

func consumer(id int, ch <-chan int, results chan<- int) {
	conSleeptime := 500
	fmt.Println(time.Now().UnixMilli(), "Starting consumer", id)
	for msg := range ch {
		fmt.Println(time.Now().UnixMilli(), id, "Processing: ", msg)
		time.Sleep(time.Millisecond * time.Duration(conSleeptime))
		results <- msg * 100
		fmt.Println(time.Now().UnixMilli(), id, "Processed: ", msg)
	}
	fmt.Println(time.Now().UnixMilli(), "Ending consumer", id)
}

func producer(tasks chan<- int) {
	proSleepTime := 20
	// produce data
	for i := range 5000 {
		tasks <- i
		fmt.Println(time.Now().UnixMilli(), "Produced", i)
		time.Sleep(time.Millisecond * time.Duration(proSleepTime))
	}
	fmt.Println(time.Now().UnixMilli(), "Finished producing")
	close(tasks)
}

func main() {
	var resultsSlice []int
	// fan out
	tasks := make(chan int, 200)
	results := make(chan int, 100)
	totalConsumers := 5
	var conWg sync.WaitGroup
	// consumers
	for i := range totalConsumers {
		conWg.Go(func() {
			consumer(i, tasks, results)
		})
	}

	// producer
	go producer(tasks)

	// close results channel once all consumers have finished
	go func() {
		conWg.Wait()
		close(results)
	}()

	// fan in
	fanInSleepTime := 1
	fmt.Println(time.Now().UnixMilli(), "Starting fan in")
	for msg := range results {
		fmt.Println(time.Now().UnixMilli(), "grouping result", msg)
		resultsSlice = append(resultsSlice, msg)
		time.Sleep(time.Millisecond * time.Duration(fanInSleepTime))
	}
	fmt.Println(time.Now().UnixMilli(), "All consumers finished")
	fmt.Println(time.Now().UnixMilli(), "Results", resultsSlice)
}
