package configure

import (
	"time"
)

// TCP ...
type TCP struct {
	Addr    string
	TimeOut time.Duration
}

// Format ...
func (t *TCP) Format(basePath string) (e error) {
	t.TimeOut *= time.Millisecond
	if t.TimeOut < 30*time.Second {
		t.TimeOut = 30 * time.Second
	}
	return
}
