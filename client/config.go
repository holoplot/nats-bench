package client

import (
	"fmt"

	"github.com/holoplot/nats-bench/consumer"
)

type Config struct {
	NatsURL              string
	Approach             consumer.Approach
	NumRealms            int
	NumConsumers         int
	NumRealmsPerConsumer int
	Suffixes             []string
}

func (c *Config) String() string {
	return fmt.Sprintf("Realms: %d, Consumers: %d, Realms / Consumer: %d, Approach: %s",
		c.NumRealms, c.NumConsumers, c.NumRealmsPerConsumer, c.Approach)
}

// runConfig returns a Config adjusted for the approach taken
func (c *Config) runConfig() Config {
	runConfig := Config{
		NatsURL:   c.NatsURL,
		Approach:  c.Approach,
		Suffixes:  c.Suffixes,
		NumRealms: c.NumRealms,
	}

	switch c.Approach {
	case consumer.MultipleFilterSubjects:
		runConfig.NumConsumers = c.NumConsumers
		runConfig.NumRealmsPerConsumer = c.NumRealmsPerConsumer
	case consumer.ManyConsumers:
		runConfig.NumConsumers = c.NumConsumers * c.NumRealmsPerConsumer
		runConfig.NumRealmsPerConsumer = 1
	case consumer.Wildcard:
		runConfig.NumConsumers = c.NumConsumers
		runConfig.NumRealmsPerConsumer = 0
	}

	return runConfig
}
