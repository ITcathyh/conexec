package conexec

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// TimeOut
var ErrorTimeOut = fmt.Errorf("TimeOut")

// Job Type
type Job func() error

// Base struct
type Actuator struct {
	timeOut *time.Duration
}

// NewActuator is used to create actuator instance
func NewActuator() *Actuator {
	return &Actuator{

	}
}

// WithTimeOut is used to set timeout
func (c *Actuator) WithTimeOut(t time.Duration) *Actuator {
	c.timeOut = &t
	return c
}

// Exec is used to run jobs concurrently
func (c *Actuator) Exec(jobs ...Job) error {
	return c.ExecWithContext(context.Background(), jobs...)
}

// ExecWithContext is used to run jobs concurrently
// Return nil when jobs are all completed successfully,
// or return error when some exception happen such as timeout
func (c *Actuator) ExecWithContext(ctx context.Context, jobs ...Job) error {
	l := len(jobs)
	if l == 0 {
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	resChan := make(chan error, l)
	wg := &sync.WaitGroup{}
	wg.Add(l)

	// Make sure the jobs are completed and channel is closed
	go func() {
		wg.Wait()
		cancel()
		close(resChan)
	}()

	// Sadly we can not kill a goroutine manually
	// So when an error happens, the other jobs will continue
	// But the good news is that main progress
	// will know the error immediately
	for _, job := range jobs {
		go func(f Job) {
			defer func() {
				wg.Done()

				if r := recover(); r != nil {
					err := fmt.Errorf("conexec panic:%v, info:%s", r, string(debug.Stack()))
					resChan <- err
				}
			}()

			if err := f(); err != nil {
				resChan <- err
			}
		}(job)
	}

	return c.wait(ctx, resChan)
}

// wait waits for the notification of execution result
func (c *Actuator) wait(ctx context.Context, resChan chan error) error {
	if c.timeOut != nil {
		return c.waitWithTimeout(ctx, resChan)
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-resChan:
		return err
	}
}

// waitWithTimeout is used to waits for the notification of execution result
// when the timeout is set
func (c *Actuator) waitWithTimeout(ctx context.Context, resChan chan error) error {
	select {
	case <-time.After(*c.timeOut):
		return ErrorTimeOut
	case <-ctx.Done():
		return nil
	case err := <-resChan:
		return err
	}
}
