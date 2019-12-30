package conexec

import (
	"testing"
	"time"
)

func TestTimeOut(t *testing.T) {
	timeout := time.Millisecond*50
	opt := &Options{TimeOut:&timeout}
	c := NewActuator(opt)
	testTimeout(t, c)
}

func TestError(t *testing.T) {
	timeout := time.Second
	opt := &Options{TimeOut:&timeout}
	c := NewActuator(opt)
	testError(t, c)
}

func TestNormal(t *testing.T) {
	c := NewActuator()
	testNormal(t, c)

	timeout := time.Minute
	opt := &Options{TimeOut:&timeout}
	c = NewActuator(opt)
	testNormal(t, c)
}

func TestEmpty(t *testing.T) {
	c := NewActuator()
	testEmpty(t, c)
}

func TestPanic(t *testing.T) {
	c := NewActuator()
	testPanic(t, c)
}
