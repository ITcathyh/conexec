package conexec

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
)

func testTimeout(t *testing.T, c TimedActuator) {
	st := time.Now().UnixNano()
	err := c.Exec(
		func() error {
			fmt.Println(1)
			time.Sleep(time.Millisecond * 100)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Millisecond * 200)
			fmt.Println(3)
			return nil
		},
	)

	Equal(t, err, ErrorTimeOut)
	et := time.Now().UnixNano()
	t.Logf("used time:%ds", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testError(t *testing.T, c TimedActuator) {
	st := time.Now().UnixNano()
	te := errors.New("TestErr")
	err := c.Exec(
		func() error {
			fmt.Println(1)
			time.Sleep(time.Millisecond * 100)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Millisecond * 200)
			fmt.Println(3)
			return nil
		},
		func() error {
			fmt.Println("4")
			return te
		},
		func() error {
			time.Sleep(time.Millisecond * 300)
			fmt.Println("5")
			return te
		},
		func() error {
			time.Sleep(time.Second)
			fmt.Println("6")
			return te
		},
	)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%ds", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testNormal(t *testing.T, c TimedActuator) {
	fs := []Task{
		func() error {
			fmt.Println(1)
			time.Sleep(time.Millisecond * 100)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Millisecond * 200)
			fmt.Println(3)
			return nil
		},
	}

	Equal(t, c.Exec(fs...), nil)
}

func testPanic(t *testing.T, c TimedActuator) {
	NotEqual(t, c.Exec(
		func() error {
			var i *int64
			num := *i + 1
			fmt.Println(num)
			return nil
		}), nil)
}

func testEmpty(t *testing.T, c TimedActuator) {
	Equal(t, c.Exec(), nil)
}
