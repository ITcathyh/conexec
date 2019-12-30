package conexec

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
)

func TestTimeOutPtr(t *testing.T) {
	timeout := time.Minute
	Equal(t, timeout, *TimeOutPtr(timeout))
}
