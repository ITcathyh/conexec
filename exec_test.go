package conexec

import (
	"fmt"
	"strings"
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

func TestExecWithError(t *testing.T) {
	Equal(t, ExecWithError(getTasks()...), nil)
	err := fmt.Errorf("TestErr")
	tasks, _ := getErrorTask()
	Equal(t, ExecWithError(tasks...), err)
	Equal(t, strings.Contains(ExecWithError(getPanicTask()).Error(), "panic"), true)
	Equal(t, ExecWithError(), nil)
}
