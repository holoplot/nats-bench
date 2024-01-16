package consumer

import (
	"context"
	"sync"

	"github.com/holoplot/nats-bench/utils"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Consumer struct {
	config      Config
	numMessages int
}

func New(config Config, totalRealms int) *Consumer {
	n := 0

	switch config.Approach {
	case MultipleFilterSubjects, ManyConsumers:
		n = len(config.Realms)
	case Wildcard:
		n = totalRealms
	}

	return &Consumer{
		config:      config,
		numMessages: n * len(config.Suffixes),
	}
}

func (c *Consumer) NumMessages() int {
	return c.numMessages
}

func (c *Consumer) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(c.NumMessages())

	onConnected := func(ctx context.Context, nc *nats.Conn) {
		js, err := jetstream.New(nc)
		if err != nil {
			panic(err)
		}

		cons, err := js.OrderedConsumer(ctx, c.config.Stream, c.config.consumerConfig())
		if err != nil {
			panic(err)
		}

		consContext, err := cons.Consume(func(msg jetstream.Msg) {
			wg.Done()
		})
		if err != nil {
			panic(err)
		}

		go func() {
			<-ctx.Done()

			consContext.Stop()
		}()
	}

	if err := utils.NatsConnect(ctx, c.config.NatsURL, c.config.ClientID, onConnected); err != nil {
		return err
	}

	wg.Wait()

	return nil
}
