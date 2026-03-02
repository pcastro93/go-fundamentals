package main

import (
	"fmt"
	"sync"
	"time"
)

func goRoutinesWithoutGo() {
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(i)
		}()
	}
	wg.Wait()
}

func goRoutinesWithGo() {
	var wg sync.WaitGroup
	for i := range 5 {
		wg.Go(func() {
			fmt.Println(i)
		})
	}
	wg.Wait()
}

func channelsSimple() {
	ch := make(chan int)
	go func() {
		ch <- 1
	}()
	message := <-ch
	fmt.Println("Message is", message)
}

func multipleGoRoutines() {
	ch := make(chan int)
	var wg sync.WaitGroup
	for i := range 5 {
		tmpFunc := func() {
			ch <- i
		}
		wg.Go(tmpFunc)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for message := range ch {
		fmt.Println(message)
	}
}

func bufferedChannels() {
	prodSleepTime := 100
	conSleepTime := 400
	ch := make(chan int, 3)
	var wg sync.WaitGroup
	// producer
	produce := func() {
		fmt.Printf("%d: Starting producer\n", time.Now().UnixMilli())
		for i := range 12 {
			ch <- i
			fmt.Printf("%d: Produced: %d\n", time.Now().UnixMilli(), i)
			time.Sleep(time.Millisecond * time.Duration(prodSleepTime))
		}
		fmt.Printf("%d: Ending producer\n", time.Now().UnixMilli())
		close(ch)
		fmt.Println(time.Now().UnixMilli(), "Channel closed")
	}
	// consumer
	consume := func() {
		fmt.Printf("%d: Starting consumer\n", time.Now().UnixMilli())
		for message := range ch {
			fmt.Printf("%d: Consumed: %d\n", time.Now().UnixMilli(), message)
			time.Sleep(time.Millisecond * time.Duration(conSleepTime))
		}
		fmt.Printf("%d: Ending consumer\n", time.Now().UnixMilli())
	}
	wg.Go(produce)
	wg.Go(consume)
	wg.Wait()
}

func main() {
	goRoutinesWithoutGo()
	fmt.Println("==========")
	goRoutinesWithGo()
	fmt.Println("==========")
	channelsSimple()
	fmt.Println("==========")
	multipleGoRoutines()
	fmt.Println("==========")
	bufferedChannels()
}
