#first stage - builder
FROM golang:1.11.0-stretch as builder
COPY . /LoadLocation
WORKDIR /LoadLocation

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o rec-engine

#second stage
FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache tzdata
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /rec-engine .
CMD ["./rec-engine"]