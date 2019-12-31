package conexec

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
)

const (
	DefaultWorkerNum = 4
)

var (
	ErrorUsingActuator = fmt.Errorf("ErrorUsingActuator")
)

// Actuator which has a worker pool
type PooledActuator struct {
	timeout *time.Duration

	workerNum int
	pool      *ants.Pool

	initOnce sync.Once
}

// NewPooledActuator creates an PooledActuator instance
func NewPooledActuator(workerNum int, opt ...*Options) *PooledActuator {
	c := &PooledActuator{
		workerNum: workerNum,
	}
	setOptions(c, opt...)
	return c
}

// Exec is used to run tasks concurrently
func (c *PooledActuator) Exec(tasks ...Task) error {
	return c.ExecWithContext(context.Background(), tasks...)
}

// ExecWithContext uses goroutine pool to run tasks concurrently
// Return nil when tasks are all completed successfully,
// or return error when some exception happen such as timeout
func (c *PooledActuator) ExecWithContext(ctx context.Context, tasks ...Task) error {
	c.initOnce.Do(func() {
		c.initPooledActuator()
	})

	if c.workerNum == -1 {
		return ErrorUsingActuator
	}

	return execTasks(c, ctx, c.runWithPool, tasks...)
}

// GetTimeout return the timeout set before
func (c *PooledActuator) GetTimeout() *time.Duration {
	return c.timeout
}

// Release will release the pool
func (c *PooledActuator) Release() {
	if c.pool != nil {
		c.pool.Release()
	}
}

// initPooledActuator init the pooled actuator once while the runtime
// If the workerNum is zero or negative,
// default worker num will be used
func (c *PooledActuator) initPooledActuator() {
	if c.workerNum <= 0 {
		c.workerNum = DefaultWorkerNum
	}

	var err error
	c.pool, err = ants.NewPool(c.workerNum)

	if err != nil {
		c.workerNum = -1
		log.Errorf("initPooledActuator err")
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
