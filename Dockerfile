FROM golang:1.26.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o scheduler .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

RUN mkdir -p /data

CMD ["./scheduler"]