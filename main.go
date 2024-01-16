package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/holoplot/nats-bench/client"
	"github.com/holoplot/nats-bench/consumer"
)

func main() {
	approach := flag.String("approach", "multiple-filter-subjects", fmt.Sprintf("One of %v", consumer.Approaches()))
	numRealms := flag.Int("num-realms", 10000, "Number of items in the config realm")
	numConsumers := flag.Int("num-consumers", 10, "Number of consumers")
	numRealmsPerConsumer := flag.Int("num-realms-per-consumer", 10, "Number of subjects each consumer subscribes to")
	flag.Parse()

	config := client.Config{
		Approach:             consumer.NewApproach(*approach),
		NumRealms:            *numRealms,
		NumConsumers:         *numConsumers,
		NumRealmsPerConsumer: *numRealmsPerConsumer,
		Suffixes:             []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
	}

	if natsURL, ok := os.LookupEnv("NATS_URL"); ok {
		config.NatsURL = natsURL
	} else {
		config.NatsURL = "localhost:4222"
	}

	client := client.New(config)
	result := client.Run(context.Background())

	fmt.Printf("%s\n", config.String())
	fmt.Printf("%s\n", result.String())
}
