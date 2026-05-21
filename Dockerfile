FROM golang:1.26.2 AS builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/

FROM scratch


WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.yaml ./config.yaml
COPY --from=builder /app/app .

CMD ["./app"]
