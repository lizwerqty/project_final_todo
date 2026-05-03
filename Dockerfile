FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/app .

COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./app"]