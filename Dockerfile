FROM ${DOCKER_HUB}golang:1.22-alpine AS builder
COPY ["go.mod", "go.sum", "./"]
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
COPY . /opt/build
RUN --mount=type=cache,target=/root/.cache/go-build cd /opt/build && go build -o nats-bench

FROM ${DOCKER_HUB}alpine:3.20
COPY --from=builder /opt/build/nats-bench /opt/

CMD /opt/nats-bench
