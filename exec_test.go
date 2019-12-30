package conexec

import (
	"testing"

	. "github.com/go-playground/assert/v2"
)

func TestExec(t *testing.T) {
	Equal(t, Exec(getTasks()...), true)
	tasks, _ := getErrorTask()
	Equal(t, Exec(tasks...), false)
	Equal(t, Exec(getPanicTask()), false)
	Equal(t, Exec(), true)
}
