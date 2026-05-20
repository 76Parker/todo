FROM golang:1.26.2 AS builder

WORKDIR /app
COPY . .
RUN go build -o app ./cmd/

FROM scratch

COPY --from=builder /app/migrations /migrations
COPY --from=builder /app/config.yaml /app/config.yaml
COPY --from=builder /app/app /app/app

CMD ["./app"]
