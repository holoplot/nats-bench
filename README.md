# NATS Bench

A performance benchmark for the [NATS server](https://github.com/nats-io/nats-server).

## Motivation

This benchmark focuses on the performance of [JetStream](https://docs.nats.io/nats-concepts/jetstream) consumers with multiple filter subjects. It does so because we have found that this use case performs worse than expected. This benchmark provides the means to reproduce this finding easily. It also captures a profile of the NATS server in order to identify code paths that could potentially be optimized.

### Stream Setup

The stream used by this benchmark provides configuration for a large set of realms identified by UUIDs. Each realm gets 10 configuration messages:

- `config.4b6bf71a-a097-11ee-8ef1-b445062a7e37.a`
- `config.4b6bf71a-a097-11ee-8ef1-b445062a7e37.b`
- `config.4b6bf71a-a097-11ee-8ef1-b445062a7e37.c`
- ...
- `config.4b6bf71a-a097-11ee-8ef1-b445062a7e37.j`

The messages itself are very small with a payload of 32 bytes.

### Consumer Setup

A configurable number of consumers fetches configuration messages. Each consumer is interested in a configurable number of realms. For each realm, identified by a UUID, it installs a filter subject that matches all ten configuration messages:

- `config.0b69a102-a098-11ee-8389-b445062a7e37.>`
- `config.4b6bf71a-a097-11ee-8ef1-b445062a7e37.>`
- `config.95a45516-a097-11ee-9fb4-b445062a7e37.>`
- `config.e69b9470-a097-11ee-90d4-b445062a7e37.>`

## Prerequisites

- Docker 18.06+
- Docker Compose 1.24.0+ (format: 3.7+)

## Compiling the benchmark

Build the benchmark client

```shell
docker compose build
```

## Executing the benchmark

Set the NATS image to test in the `.env` file:

```shell
NATS_IMAGE=nats:2.10-alpine
```

Or use the nightly image:

```shell
NATS_IMAGE=synadia/nats-server:nightly
```

Run the benchmark with default parameters:

```shell
docker compose up
```

Press Ctrl-C, or run `docker compose down`, when done.

### Parameters

Parameters can be specifed in the `.env` file:

```shell
NUM_REALMS=10000
NUM_CONSUMERS=200
NUM_REALMS_PER_CONSUMER=50
```

or as environment variables passed to `docker-compose`:

```shell
NUM_CONSUMERS=100 docker-compose up
```

### Approaches

#### Multiple filter subjects

This benchmark defaults to the approach where each client creates one consumer that has multiple filter subjects. This minimizes the number of consumers and it minimizes network traffic. Each client only receives the data it needs, the filtering is done on the NATS server.

#### Many consumers

An alternative approach is to let each client install many consumers, each with only a single filter subject. This also results in each client only receiving the data it needs, again the filtering is done on the NATS server. However the number of consumers is much larger. You can also benchmark this approach:

```shell
APPROACH=many-consumers docker-compose up
```

#### Wildcard subscription

Another alternative is to let each client install a single consumer using a wildcard subscription that matches all messages. Each client will receive all the data, also the data it doesn't actually need. This results in much higher network traffic, the filtering then needs to be done on the client side. You can also benchmark this approach:

```shell
APPROACH=wildcard docker-compose up
```


## Profiling

The `pprof` container obtains a profile from the NATS server under test. If your benchmark runs longer than 30 seconds, or if you wait long enough at the end of the benchmark, a `nats.profile` will appear in the `profiles`
folder. This file can be used with the [pprof](https://pkg.go.dev/net/http/pprof) tool:

```shell
go tool pprof profiles/nats.profile
```

## Use w/o Docker Compose

The use of `docker-compose` guarantees that the NATS server is freshly started for each run of the benchmark. If you want to make changes to the client, a faster development cycle may be desirable.

Start a NATS server and keep it running:

```shell
docker run -p 4222:4222 nats:2.10-alpine nats-server --jetstream
```

Compile the client locally:

```shell
go build -o nats-bench
```

Run the client locally:

```shell
./nats-bench -num-realms 1000 -num-consumers 10 -num-realms-per-consumer 100
```
