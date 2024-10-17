package semaphore

import (
	"sync"
)

// Semaphore is a struct that controls access to a finite number of resources
// or limits the number of concurrent operations. It utilizes a buffered channel
// to manage the number of allowed concurrent executions and a WaitGroup to wait
// for all operations to complete.
type Semaphore struct {
	channel chan struct{}
	wg      sync.WaitGroup
}

// New creates a new Semaphore with a specified size, which determines the maximum
// number of concurrent operations allowed.
func New(size int) *Semaphore {
	return &Semaphore{
		channel: make(chan struct{}, size),
	}
}

// Run accepts a function and executes it concurrently, while ensuring that the
// maximum concurrency limit is respected. It acquires a semaphore before running
// the function and releases it after the function completes.
func (s *Semaphore) Run(fn func()) {
	s.wg.Add(1)

	go func() {
		defer s.wg.Done()

		s.Acquire()

		fn()

		s.Release()
	}()
}

// Wait blocks until all operations executed by Run have completed.
func (s *Semaphore) Wait() {
	s.wg.Wait()
}

// Acquire blocks until the semaphore has capacity, at which point it decrements
// the available slots by one, allowing a new operation to run.
func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

// Release increments the available slots in the semaphore by one, allowing another
// operation to run.
func (s *Semaphore) Release() {
	<-s.channel
}
