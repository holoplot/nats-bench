package client

import (
	"fmt"
)

type Config struct {
	NatsURL              string
	NumRealms            int
	NumConsumers         int
	NumRealmsPerConsumer int
	Suffixes             []string
}

func (c *Config) String() string {
	return fmt.Sprintf("Realms: %d, Consumers: %d, Realms / Consumer: %d",
		c.NumRealms, c.NumConsumers, c.NumRealmsPerConsumer)
}
