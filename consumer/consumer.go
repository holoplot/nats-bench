package consumer

import (
	"context"
	"strings"
	"sync"

	"github.com/holoplot/nats-bench/utils"
	nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Config struct {
	NatsURL  string
	Stream   string
	ClientID string
	Realms   []string
	Suffixes []string
}

type Consumer struct {
	config Config
}

func New(config Config) *Consumer {
	return &Consumer{
		config: config,
	}
}

func (c *Consumer) NumMessages() int {
	return len(c.config.Realms) * len(c.config.Suffixes)
}

func (c *Consumer) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(c.NumMessages())

	onConnected := func(ctx context.Context, nc *nats.Conn) {
		js, err := jetstream.New(nc)
		if err != nil {
			panic(err)
		}

		config := jetstream.OrderedConsumerConfig{
			DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		}

		for _, realm := range c.config.Realms {
			subject := strings.Join([]string{"config", realm, ">"}, ".")
			config.FilterSubjects = append(config.FilterSubjects, subject)
		}

		cons, err := js.OrderedConsumer(ctx, c.config.Stream, config)
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
