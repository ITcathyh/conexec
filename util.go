package conexec

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

// DurationPtr helps to make a duration ptr
func DurationPtr(t time.Duration) *time.Duration {
	return &t
}

// wrapperTask will wrapper the task in order to notice execution result
// to the main process
func wrapperTask(ctx context.Context, task Task,
	wg *sync.WaitGroup, resChan chan error) func() {
	return func() {
		tmpChan := make(chan error, 1)
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("conexec panic:%v\n%s", r, string(debug.Stack()))
				resChan <- err
			}

			wg.Done()
			close(tmpChan)
		}()

		select {
		case <-ctx.Done():
			return // fast return
		case tmpChan <- task():
			if err := <-tmpChan; err != nil {
				resChan <- err
			}
		}

		// if err := task(); err != nil {
		// 	resChan <- err
		// }
	}
}

// setOptions set the options for actuator
func setOptions(c TimedActuator, options ...*Options) {
	if options == nil || len(options) == 0 || options[0] == nil {
		return
	}

	opt := options[0]
	if opt.TimeOut != nil {
		c.setTimeout(opt.TimeOut)
	}
}
