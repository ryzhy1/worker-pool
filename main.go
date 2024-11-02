package main

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id    int
	input chan string
	done  chan bool
}

type WorkerPool struct {
	workers []*Worker
	input   chan string
	mu      sync.Mutex
}

func NewWorkerPool() *WorkerPool {
	return &WorkerPool{
		input: make(chan string),
	}
}

func (p *WorkerPool) AddWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	workerID := len(p.workers) + 1
	worker := &Worker{
		id:    workerID,
		input: p.input,
		done:  make(chan bool),
	}
	p.workers = append(p.workers, worker)

	go worker.Start()
}

func (p *WorkerPool) RemoveWorker() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.workers) == 0 {
		return
	}

	worker := p.workers[len(p.workers)-1]
	p.workers = p.workers[:len(p.workers)-1]

	worker.done <- true
}

func (w *Worker) Start() {
	for {
		select {
		case data := <-w.input:
			fmt.Printf("Worker %d, data: %s\n", w.id, data)
			time.Sleep(5 * time.Millisecond) // Представим что у нас реально происходят какие-то тяжелые вычисления.
			// Я это делаю для того, чтобы избежать случая, когда один конкретный воркер забирает себе все
		case <-w.done:
			fmt.Printf("Worker %d stopping\n", w.id)
			return
		}
	}
}

func main() {
	wp := NewWorkerPool()

	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()

	for i := 0; i < 5; i++ {
		wp.input <- fmt.Sprintf("task %d", i)
	}

	wp.RemoveWorker()

	for i := 6; i <= 10; i++ {
		wp.input <- fmt.Sprintf("task %d", i)
	}
}
