package conexec

import (
	"testing"
	"time"
)

func TestPooledTimeOut(t *testing.T) {
	timeout := time.Millisecond*50
	opt := &Options{TimeOut:&timeout}

	c := NewPooledActuator(5, opt)
	testTimeout(t, c)
	c = NewPooledActuator(-1, opt)
	testTimeout(t, c)
}

func TestPooledError(t *testing.T) {
	timeout := time.Second
	opt := &Options{TimeOut:&timeout}

	c := NewPooledActuator(5, opt)
	testError(t, c)
}

func TestPooledNormal(t *testing.T) {
	c := NewPooledActuator(5)
	testNormal(t, c)

	timeout := time.Minute
	opt := &Options{TimeOut:&timeout}
	c = NewPooledActuator(5, opt)
	testNormal(t, c)
}

func TestPooledEmpty(t *testing.T) {
	c := NewPooledActuator(5)
	testEmpty(t, c)
}

func TestPooledPanic(t *testing.T) {
	c := NewPooledActuator(5)
	testPanic(t,c)
}