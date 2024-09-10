package client

import (
	"fmt"
	"time"
)

type Result struct {
	ConsumedMessages int
	Elapsed          time.Duration
}

func (r *Result) Rate() float64 {
	return float64(r.ConsumedMessages) / r.Elapsed.Seconds()
}

func (r *Result) String() string {
	return fmt.Sprintf("Consumed %d messages in %s (%.1f msgs/sec)", r.ConsumedMessages, r.Elapsed, r.Rate())
}
