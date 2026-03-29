FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o scheduler-app main.go

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/scheduler-app .

COPY --from=builder /app/web ./web

COPY --from=builder /app/config.yaml .

RUN mkdir -p /data

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

CMD ["./scheduler-app"]