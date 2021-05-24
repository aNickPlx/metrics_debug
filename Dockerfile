FROM docker.io/golang:1.15-buster AS builder

WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 go build -o metrics_debug github.com/aNickPlx/metrics_debug

FROM docker.io/alpine:3

COPY --from=builder /build/metrics_debug .

CMD ["./metrics_debug"]
