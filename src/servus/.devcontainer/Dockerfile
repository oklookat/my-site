FROM golang:1.18.2-alpine3.15

# build
WORKDIR /app
COPY . .
RUN apk update && apk add build-base && go mod download