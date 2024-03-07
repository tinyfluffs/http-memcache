FROM golang:1.22 as build

WORKDIR /app
COPY / ./
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o memcache cmd/memcache.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    adduser -D -h /memcache memcache

COPY --from=build /app/memcache /bin/memcache

USER memcache
WORKDIR /memcache

ENV ADDRESS=0.0.0.0:9000
ENV EXPIRATION=30m
ENV GC_INTERVAL=1m
ENV CHUNK_COUNT=2048

EXPOSE 9000
CMD ["/bin/memcache"]
