package conexec

import (
	"fmt"
	"runtime/debug"
	"sync"
)

// wrapperTask will wrapper the task in order to notice execution result
// to the main process
func wrapperTask(task Task, wg *sync.WaitGroup, resChan chan error) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("conexec panic:%v\n%s", r, string(debug.Stack()))
				resChan <- err
			}

			wg.Done()
		}()

		if err := task(); err != nil {
			resChan <- err
		}
	}
}

// setOptions set the options for actuator
func setOptions(c TimedActuator, options ...*Options) {
	if options == nil || len(options) == 0 {
		return
	}

	opt := options[0]
	if opt.TimeOut != nil {
		c.setTimeout(opt.TimeOut)
	}
}