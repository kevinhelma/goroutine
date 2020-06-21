package main

import "sync"

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