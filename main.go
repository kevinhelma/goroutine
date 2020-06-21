package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"time"
	)

type Task struct {
	closed chan struct{}
	wg     sync.WaitGroup
}

func (t *Task) Run() {
	for {
		select {
		case <-t.closed:
			return
		default:
			handle()
			os.Exit(0)
		}
	}
}

func (t *Task) Stop() {
	close(t.closed)
	t.wg.Wait()
}

func handle() {
	// # Read integer
	var i int
	fmt.Print("Input jumlah loket: ")
	fmt.Scanf("%d", &i)

	if i < 1 || i > 5 {
		fmt.Println("Loket harus antara 1-5")
		os.Exit(0)
	}

	//queque in second
	queque := []int{1, 2, 4, 2, 3, 5, 2, 3, 1, 3}
	numJobs := len(queque)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= i; w++ {
		go worker(queque, w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		<-results
	}
}

func worker(queque []int, id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("counter", id, "started queue", j)
		time.Sleep(time.Duration(queque[j-1]) * time.Second)
		fmt.Println("counter", id, "finished queue", j, " in ", queque[j-1], " second")
		results <- j
	}
}

func main() {
	task := &Task{
		closed: make(chan struct{}),
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	task.wg.Add(1)
	go func() { defer task.wg.Done(); task.Run() }()

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
		<-c
		task.Stop()
	}
}