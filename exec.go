package conexec

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"
)

// wait waits for the notification of execution result
func wait(c TimedActuator, ctx context.Context, resChan chan error) error {
	if timeout := c.GetTimeout(); timeout != nil {
		return waitWithTimeout(ctx, resChan, *timeout)
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-resChan:
		return err
	}
}

// waitWithTimeout waits for the notification of execution result
// when the timeout is set
func waitWithTimeout(ctx context.Context, resChan chan error, timeout time.Duration) error {
	select {
	case <-time.After(timeout):
		return ErrorTimeOut
	case <-ctx.Done():
		return nil
	case err := <-resChan:
		return err
	}
}

// execTasks uses customized function to
// execute every task, such as using the simplyRun
func execTasks(c TimedActuator, parent context.Context,
	execFunc func(f func()), tasks ...Task) error {
	size := len(tasks)
	if size == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(parent)
	resChan := make(chan error, size)
	wg := &sync.WaitGroup{}
	wg.Add(size)

	// Make sure the tasks are completed and channel is closed
	go func() {
		wg.Wait()
		cancel()
		close(resChan)
	}()

	// Sadly we can not kill a goroutine manually
	// So when an error happens, the other tasks will continue
	// But the good news is that main progress
	// will know the error immediately
	for _, task := range tasks {
		child,_ := context.WithCancel(ctx)
		f := wrapperTask(child,task, wg, resChan)
		execFunc(f)
	}

	return wait(c, ctx, resChan)
}

// simplyRun uses a new goroutine to run the function
func simplyRun(f func()) {
	go f()
}

// Exec simply runs the tasks concurrently
// True will be returned is all tasks complete successfully
// otherwise false will be returned
func Exec(tasks ...Task) bool {
	var c int32
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))

	for _, t := range tasks {
		go func(task Task) {
			defer func() {
				if r := recover(); r != nil {
					atomic.StoreInt32(&c, 1)
					fmt.Printf("conexec panic:%v\n%s\n", r, string(debug.Stack()))
				}

				wg.Done()
			}()

			if err := task(); err != nil {
				atomic.StoreInt32(&c, 1)
			}
		}(t)
	}

	wg.Wait()
	return c == 0
}