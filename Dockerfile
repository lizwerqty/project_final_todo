FROM golang:1.25.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

FROM alpine

WORKDIR /app

COPY --from=builder /app/app .

COPY --from=builder /app/web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./app"]