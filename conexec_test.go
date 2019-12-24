package conexec

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestTimeOut(t *testing.T) {
	c := NewActuator()
	c.WithTimeOut(time.Millisecond * 500)
	st := time.Now().UnixNano()
	t.Log(c.Exec(
		func() error {
			fmt.Println(1)
			time.Sleep(time.Second * 2)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Second * 2)
			fmt.Println(3)
			return nil
		},
	))
	et := time.Now().UnixNano()
	t.Logf("used time:%ds", (et-st)/1000000)
	time.Sleep(time.Second * 5)
}

func TestError(t *testing.T) {
	c := NewActuator()
	c.WithTimeOut(time.Millisecond * 500)
	st := time.Now().UnixNano()
	t.Log(c.Exec(
		func() error {
			fmt.Println(1)
			time.Sleep(time.Second * 2)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Second * 2)
			fmt.Println(3)
			return nil
		},
		func() error {
			fmt.Println("4")
			return errors.New("TestErr")
		},
	))
	et := time.Now().UnixNano()
	t.Logf("used time:%ds", (et-st)/1000000)
	time.Sleep(time.Second * 5)
}

func TestNormal(t *testing.T) {
	c := NewActuator()
	st := time.Now().UnixNano()
	t.Log(c.Exec(
		func() error {
			fmt.Println(1)
			time.Sleep(time.Second * 2)
			return nil
		},
		func() error {
			fmt.Println(2)
			return nil
		},
		func() error {
			time.Sleep(time.Second * 1)
			fmt.Println(3)
			return nil
		},
	))
	et := time.Now().UnixNano()
	t.Logf("used time:%ds", (et-st)/1000000)
}
