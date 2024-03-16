# syntax = docker/dockerfile:1

ARG GO_VERSION="1.20"
ARG RUNNER_IMAGE="alpine:3.16"

FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /opt
RUN apk add --no-cache make git gcc musl-dev openssl-dev linux-headers

COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download

# Copy the remaining files
COPY . .

RUN LINK_STATICALLY=true make build

# Add to a distroless container
FROM ${RUNNER_IMAGE}

COPY --from=builder /opt/build/intentod /usr/local/bin/intentod
RUN apk add bash vim sudo dasel \
    && addgroup -g 1000 intento \
    && adduser -S -h /home/intento -D intento -u 1000 -G intento 

RUN mkdir -p /etc/sudoers.d \
    && echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel \
    && echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers \
    && adduser intento wheel 

USER 1000
ENV HOME /home/intento
WORKDIR $HOME

EXPOSE 26657 26656 1317 9090

CMD ["intentod", "start"]