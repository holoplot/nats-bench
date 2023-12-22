package client

import (
	"fmt"
	"time"
)

type Result struct {
	ConsumedMessages int
	Elapsed          time.Duration
}

func (r *Result) String() string {
	return fmt.Sprintf("Consumed %d messages in %s", r.ConsumedMessages, r.Elapsed)
}
