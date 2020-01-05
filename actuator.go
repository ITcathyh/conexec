package conexec

import (
	"context"
	"fmt"
	"time"
)

// BaseActuator is the actuator interface
type BaseActuator interface {
	Exec(tasks ...Task) error
	ExecWithContext(ctx context.Context, tasks ...Task) error
}

// TimedActuator is the actuator interface within timeout method
type TimedActuator interface {
	BaseActuator
	GetTimeout() *time.Duration
	setTimeout(timeout *time.Duration)
}

// ErrorTimeOut is the error when executes tasks timeout
var ErrorTimeOut = fmt.Errorf("TimeOut")

// Task Type
type Task func() error

// Actuator is the base struct
type Actuator struct {
	timeout *time.Duration
}

// NewActuator creates an Actuator instance
func NewActuator(opt ...*Options) *Actuator {
	c := &Actuator{}
	setOptions(c, opt...)
	return c
}

// Exec is used to run tasks concurrently
func (c *Actuator) Exec(tasks ...Task) error {
	return c.ExecWithContext(context.Background(), tasks...)
}

// ExecWithContext is used to run tasks concurrently
// Return nil when tasks are all completed successfully,
// or return error when some exception happen such as timeout
func (c *Actuator) ExecWithContext(ctx context.Context, tasks ...Task) error {
	return execTasks(ctx, c, simplyRun, tasks...)
}

// GetTimeout return the timeout set before
func (c *Actuator) GetTimeout() *time.Duration {
	return c.timeout
}

// setTimeout sets the timeout
func (c *Actuator) setTimeout(timeout *time.Duration) {
	c.timeout = timeout
}
