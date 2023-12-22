FROM ${DOCKER_HUB}golang:1.21-alpine as builder
COPY . /opt/build
RUN cd /opt/build && go build -o nats-bench

FROM ${DOCKER_HUB}alpine:3.18
COPY --from=builder /opt/build/nats-bench /opt/

CMD /opt/nats-bench
