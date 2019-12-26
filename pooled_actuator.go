package conexec

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

var (
	ErrorUsingActuator = fmt.Errorf("ErrorUsingActuator")
)

// Actuator which has a worker pool
// If the workerNum is zero or negative,
// base actuator will be used only
type PooledActuator struct {
	timeout *time.Duration

	workerNum    int32
	pool         *ants.Pool
	baseActuator *Actuator

	initOnce sync.Once
	opt      []*Options
}

// NewPooledActuator creates an PooledActuator instance
func NewPooledActuator(workerNum int32, opt ...*Options) *PooledActuator {
	c := &PooledActuator{
		workerNum: workerNum,
		opt:       opt,
	}
	setOptions(c, opt...)
	return c
}

// Exec is used to run tasks concurrently
func (c *PooledActuator) Exec(tasks ...Task) error {
	return c.ExecWithContext(context.Background(), tasks...)
}

// ExecWithContext is used to run tasks concurrently
// Return nil when tasks are all completed successfully,
// or return error when some exception happen such as timeout
func (c *PooledActuator) ExecWithContext(ctx context.Context, tasks ...Task) error {
	c.initOnce.Do(func() {
		c.initPooledActuator()
	})

	if c.workerNum == -1 {
		return ErrorUsingActuator
	} else if c.pool == nil {
		return c.baseActuator.ExecWithContext(ctx, tasks...)
	}

	return execTasks(c, ctx, c.runWithPool, tasks...)
}

// GetTimeout return the timeout set before
func (c *PooledActuator) GetTimeout() *time.Duration {
	return c.timeout
}

// initPooledActuator init the pooled actuator once while the runtime
func (c *PooledActuator) initPooledActuator() {
	if c.workerNum <= 0 {
		c.baseActuator = NewActuator(c.opt...)
		c.workerNum = 0
		c.opt = nil
		return
	}

	var err error
	c.pool, err = ants.NewPool(int(c.workerNum))

	if err != nil {
		c.workerNum = -1
		log.Fatal("initPooledActuator err")
		return
	}
}

// Release will release the pool
func (c *PooledActuator) Release() {
	if c.pool != nil {
		c.pool.Release()
	}
}

// runWithPool used the goroutine pool to execute the tasks
func (c *PooledActuator) runWithPool(f func()) {
	err := c.pool.Submit(f)
	if err != nil {
		log.Errorf("submit task err:%s", err.Error())
	}
}

// setTimeout sets the timeout
func (c *PooledActuator) setTimeout(timeout *time.Duration) {
	c.timeout = timeout
}
