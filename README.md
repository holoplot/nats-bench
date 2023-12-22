# NATS Bench

A performance benchmark for the NATS server.

## Motivation

This benchmark focuses on the performance of JetStream consumers with multiple filter subjects. It does so because we have found that this use case performs worse than expected. This benchmark provides the means to reproduce this finding easily. It also captures a profile of the NATS server in order to identify code paths that could potentially be optimized.

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
docker-compose build
```

## Executing the benchmark

Set the NATS image to test in the `.env` file:

```shell
NATS_IMAGE=nats:2.10-alpine
```

or use the nightly image:

```shell
NATS_IMAGE=synadia/nats-server:nightly
```

Run the benchmark with default parameters:

```shell
docker-compose up
```

Press Ctrl-C, or run `docker-compose down`, when done.

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

## Profiling

The `pprof` container obtains a profile from the NATS server under test.
If your benchmark runs longer than 30 seconds, or if you wait long enough
at the end of the benchmark, a `nats.profile` will appear in the `profiles`
folder. This file can be used with the [pprof](https://pkg.go.dev/net/http/pprof)
tool:

```shell
go tool pprof profiles/nats.profile
```

## Use w/o Docker Compose

The use of `docker-compose` guarantees that the NATS server is freshly started for each run of the benchmark. If you want to make changes to the client, a faster development cycle may be desirable.

Start a NATS server and keep it running:

```shell
docker run -p 4222 nats:2.10-alpine nats-server --jetstream
```

Compile the client locally:

```shell
go build -o nats-bench
```

Run the client locally:

```shell
./nats-bench -num-realms 1000 -num-consumers 10 -num-realms-per-consumer 100
```
