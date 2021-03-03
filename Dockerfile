# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git gcc

# RUN go get github.com/swaggo/swag@v1.7.0
# RUN go get github.com/go-openapi/jsonreference@v0.19.5
# RUN go get github.com/swaggo/http-swagger@v1.0.0

ADD . /src
RUN cd /src && CGO_ENABLED=0 GOOS=linux go build -o rec-engine

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/rec-engine /app/
ENTRYPOINT ./rec-engine