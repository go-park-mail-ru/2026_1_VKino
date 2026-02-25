FROM golang:1.26.0-alpine3.23 AS builder

WORKDIR /app

COPY ../../go.mod  ./
RUN go mod download

COPY ../.. .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o myapp ./cmd/api

FROM golang:1.26.0-alpine3.23

WORKDIR /root/

COPY --from=builder /app/myapp .

EXPOSE 8080

ENTRYPOINT ./api

CMD ["--config"]