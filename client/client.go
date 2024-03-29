package client

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/holoplot/nats-bench/consumer"
	"github.com/holoplot/nats-bench/utils"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName        = "config"
	streamDescription = "NATS Bench Stream"
)

type Client struct {
	config Config
}

func New(config Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Run(ctx context.Context) Result {
	runConfig := c.config.runConfig()

	realms := make([]string, runConfig.NumRealms)
	for i := 0; i < runConfig.NumRealms; i++ {
		realms[i] = uuid.New().String()
	}

	preparedWg := sync.WaitGroup{}
	preparedWg.Add(1)

	onConnected := func(ctx context.Context, nc *nats.Conn) {
		fmt.Printf("Connected to NATS server v%s\n", nc.ConnectedServerVersion())

		js, err := jetstream.New(nc)
		if err != nil {
			panic(err)
		}

		js.DeleteStream(ctx, streamName)

		if _, err := js.CreateStream(ctx, streamConfig()); err != nil {
			panic(err)
		}

		n := 0
		for _, realm := range realms {
			for _, suffix := range runConfig.Suffixes {
				subject := strings.Join([]string{"config", realm, suffix}, ".")
				s := uuid.New()

				if _, err := js.PublishAsync(subject, s[:]); err != nil {
					panic(err)
				}

				n += 1
			}

			<-js.PublishAsyncComplete()
		}

		fmt.Printf("Published %d messages\n", n)

		preparedWg.Done()
	}

	fmt.Printf("Connecting to NATS (%s)\n", runConfig.NatsURL)
	if err := utils.NatsConnect(ctx, runConfig.NatsURL, "client", onConnected); err != nil {
		panic(err)
	}

	preparedWg.Wait()

	fmt.Printf("Testing approach %s\n", runConfig.Approach)

	consumers := make([]*consumer.Consumer, runConfig.NumConsumers)
	n := 0

	for i := 0; i < runConfig.NumConsumers; i++ {
		consumerConfig := consumer.Config{
			NatsURL:  runConfig.NatsURL,
			Suffixes: runConfig.Suffixes,
			Stream:   streamName,
			ClientID: fmt.Sprintf("consumer-%d", i),
			Approach: runConfig.Approach,
		}

		consumerRealms := make(map[string]struct{}, runConfig.NumRealmsPerConsumer)
		for j := 0; j < runConfig.NumRealmsPerConsumer; j++ {
			consumerRealms[realms[rand.Intn(len(realms))]] = struct{}{}
		}

		for realm := range consumerRealms {
			consumerConfig.Realms = append(consumerConfig.Realms, realm)
		}

		consumers[i] = consumer.New(consumerConfig, runConfig.NumRealms)
		n += consumers[i].NumMessages()
	}

	consumerWg := sync.WaitGroup{}
	consumerWg.Add(len(consumers))

	fmt.Printf("Creating %d consumers for a total of %d messages\n", len(consumers), n)

	start := time.Now()

	for _, c := range consumers {
		go func(c *consumer.Consumer) {
			if err := c.Run(ctx); err != nil {
				panic(err)
			}

			consumerWg.Done()
		}(c)
	}

	consumerWg.Wait()

	return Result{
		ConsumedMessages: n,
		Elapsed:          time.Since(start),
	}
}

func streamConfig() jetstream.StreamConfig {
	return jetstream.StreamConfig{
		Name:              streamName,
		Description:       streamDescription,
		Subjects:          []string{"config.>"},
		Storage:           jetstream.FileStorage,
		Retention:         jetstream.LimitsPolicy,
		MaxAge:            time.Hour,
		Duplicates:        10 * time.Second,
		Discard:           jetstream.DiscardOld,
		NoAck:             false,
		MaxMsgs:           -1,
		MaxBytes:          -1,
		MaxConsumers:      -1,
		Replicas:          1,
		MaxMsgsPerSubject: 1,
	}
}
