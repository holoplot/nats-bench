package consumer

import (
	"strings"

	"github.com/nats-io/nats.go/jetstream"
)

type Config struct {
	NatsURL  string
	Stream   string
	ClientID string
	Approach Approach
	Realms   []string
	Suffixes []string
}

func (c *Config) consumerConfig() jetstream.OrderedConsumerConfig {
	config := jetstream.OrderedConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
	}

	switch c.Approach {
	case MultipleFilterSubjects:
		for _, realm := range c.Realms {
			subject := strings.Join([]string{"config", realm, ">"}, ".")
			config.FilterSubjects = append(config.FilterSubjects, subject)
		}
	case ManyConsumers:
		if len(c.Realms) != 1 {
			panic("only a single realm is allowed with this approach")
		}
		subject := strings.Join([]string{"config", c.Realms[0], ">"}, ".")
		config.FilterSubjects = []string{subject}
	case Wildcard:
		if len(c.Realms) != 0 {
			panic("specifying realms is not allowed with this approach")
		}
		subject := "config.*.>"
		config.FilterSubjects = []string{subject}
	}

	return config
}
