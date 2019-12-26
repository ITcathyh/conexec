package conexec

import (
	"testing"
	"time"
)

func TestTimeOut(t *testing.T) {
	c := NewActuator()
	c.WithTimeOut(time.Millisecond * 50)
	testTimeout(t, c)
}

func TestError(t *testing.T) {
	c := NewActuator()
	c.WithTimeOut(time.Second)
	testError(t, c)
}

func TestNormal(t *testing.T) {
	c := NewActuator()
	testNormal(t, c)
	c.WithTimeOut(time.Minute)
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
