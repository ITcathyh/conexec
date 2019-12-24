package conexec

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// TimeOut...
var ErrorTimeOut = fmt.Errorf("TimeOut")

// Func wrapper
type Job func() error

// Base struct
type Actuator struct {
	timeOut *time.Duration
}

// NewActuator is used to create actuator instance
func NewActuator() *Actuator {
	return &Actuator{}
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

	var timeout time.Duration
	if c.timeOut != nil {
		timeout = *c.timeOut
	} else {
		timeout = time.Hour
	}

	ctx, cancel := context.WithCancel(ctx)
	resChan := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(l)

	// Make sure the jobs are completed
	// and channel is closed
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
					err := errors.New(fmt.Sprintf("conexec panic:%v, info:%s", r, string(debug.Stack())))
					resChan <- err
				}
			}()

			err := f()
			if err != nil {
				resChan <- err
			}
		}(job)
	}

	select {
	case <-time.After(timeout):
		return ErrorTimeOut
	case <-ctx.Done():
		return nil
	case err := <-resChan:
		return err
	}
}
