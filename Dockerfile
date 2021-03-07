# build stage
FROM golang:stretch AS build-env

RUN go version

ADD . /src
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -o rec-engine

# final stage
FROM golang:alpine
WORKDIR /app
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /src/rec-engine /app/
ENTRYPOINT ./rec-engine