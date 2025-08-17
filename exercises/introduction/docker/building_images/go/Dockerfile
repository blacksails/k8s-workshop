FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN ls /app

RUN go build -o my-binary

RUN ls /app

FROM alpine

WORKDIR /app

COPY --from=builder /app/my-binary /usr/local/bin/my-binary

RUN ls /app

CMD ["my-binary"]
