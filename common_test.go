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
	err := c.Exec(getTasks()...)

	Equal(t, err, ErrorTimeOut)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testError(t *testing.T, c TimedActuator) {
	st := time.Now().UnixNano()
	tasks, te := getErrorTask()
	err := c.Exec(tasks...)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testManyError(t *testing.T, c TimedActuator) {
	tasks := make([]Task, 0)
	tmp, te := getErrorTask()
	tasks = append(tasks, tmp...)

	for i := 0; i < 100; i++ {
		tmp, _ = getErrorTask()
		tasks = append(tasks, tmp...)
	}

	st := time.Now().UnixNano()
	err := c.Exec(tasks...)

	Equal(t, err, te)
	et := time.Now().UnixNano()
	t.Logf("used time:%dms", (et-st)/1000000)
	time.Sleep(time.Millisecond * 500)
}

func testNormal(t *testing.T, c TimedActuator) {
	Equal(t, c.Exec(getTasks()...), nil)
}

func testPanic(t *testing.T, c TimedActuator) {
	NotEqual(t, c.Exec(getPanicTask()), nil)
}

func testEmpty(t *testing.T, c TimedActuator) {
	Equal(t, c.Exec(), nil)
}

func getTasks() []Task {
	return []Task{
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
}

func getErrorTask() ([]Task, error) {
	te := errors.New("TestErr")

	tasks := getTasks()
	tasks = append(tasks,
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
		}, )

	return tasks, te
}

func getPanicTask() Task {
	return func() error {
		var i *int64
		num := *i + 1
		fmt.Println(num)
		return nil
	}
}
