package log

import (
	"testing"
)

func TestLog(t *testing.T) {
	Err("Hello %s !", "123")
}
