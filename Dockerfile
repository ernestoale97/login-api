FROM golang:1.21.5-alpine3.18 AS builder

LABEL maintainer="@ernestoale97"

ENV GO111MODULE=on
ENV PORT=8030

# Install git.
# Git is required for fetching the dependencies.
RUN apk add --update --no-cache \
  git \
  sqlite-dev \
  build-base

#RUN apt-get update -qq \
#    && apt-get install build-essential libsqlite3-dev

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build

FROM alpine:latest

WORKDIR /app

# Copy our static executable.
COPY --from=builder /app/login_api /app/login_api
COPY --from=builder /app/env /app/env

EXPOSE 8880
ENTRYPOINT ["/app/login_api"]