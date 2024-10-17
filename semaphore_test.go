package semaphore_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/michaeldavidkelley/go-semaphore"
)

func TestSemaphore_Run(t *testing.T) {
	t.Parallel()

	const maxConcurrent = 3
	sem := semaphore.New(maxConcurrent)

	var running int32
	var maxSeen int32
	totalTasks := 10

	for i := 0; i < totalTasks; i++ {
		sem.Run(func() {
			atomic.AddInt32(&running, 1)
			currentRunning := atomic.LoadInt32(&running)
			if currentRunning > maxSeen {
				atomic.StoreInt32(&maxSeen, currentRunning)
			}
			time.Sleep(100 * time.Millisecond) // simulate work
			atomic.AddInt32(&running, -1)
		})
	}

	sem.Wait()

	if maxSeen > int32(maxConcurrent) {
		t.Errorf("expected max concurrent operations to be %d, but got %d", maxConcurrent, maxSeen)
	}
}

func TestSemaphore_Wait(t *testing.T) {
	t.Parallel()

	const maxConcurrent = 2
	sem := semaphore.New(maxConcurrent)

	totalTasks := 5
	var completedTasks int32

	for i := 0; i < totalTasks; i++ {
		sem.Run(func() {
			time.Sleep(50 * time.Millisecond) // simulate work
			atomic.AddInt32(&completedTasks, 1)
		})
	}

	sem.Wait()

	if completedTasks != int32(totalTasks) {
		t.Errorf("expected %d completed tasks, but got %d", totalTasks, completedTasks)
	}
}

func TestSemaphore_AcquireRelease(t *testing.T) {
	t.Parallel()

	const maxConcurrent = 1
	sem := semaphore.New(maxConcurrent)

	acquired := make(chan struct{})
	released := make(chan struct{})

	// Run a goroutine to acquire the semaphore
	go func() {
		sem.Acquire()
		acquired <- struct{}{}
		// Hold the semaphore for a bit
		time.Sleep(100 * time.Millisecond)
		sem.Release()
		released <- struct{}{}
	}()

	// Ensure the semaphore was acquired
	select {
	case <-acquired:
		// Success, semaphore was acquired
	case <-time.After(50 * time.Millisecond):
		t.Error("expected semaphore to be acquired, but timed out")
	}

	// Ensure the semaphore was released
	select {
	case <-released:
		// Success, semaphore was released
	case <-time.After(150 * time.Millisecond):
		t.Error("expected semaphore to be released, but timed out")
	}
}
