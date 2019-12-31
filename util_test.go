package conexec

import (
	"testing"
	"time"

	. "github.com/go-playground/assert/v2"
)

func TestDurationPtr(t *testing.T) {
	timeout := time.Minute
	Equal(t, timeout, *DurationPtr(timeout))
}
