# syntax = docker/dockerfile:1

ARG GO_VERSION="1.22.7"
ARG RUNNER_IMAGE_VERSION="3.19"

FROM golang:${GO_VERSION}-alpine${RUNNER_IMAGE_VERSION} AS builder

WORKDIR /opt
RUN apk add --no-cache make git gcc musl-dev openssl-dev linux-headers ca-certificates build-base

COPY go.mod .
COPY go.sum .

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/root/go/pkg/mod \
    go mod download




ARG WASMVM_VERSION=v2.1.3
ADD https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/${WASMVM_VERSION}/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
    
RUN cp "/lib/libwasmvm_muslc.$(uname -m).a" /lib/libwasmvm_muslc.a
COPY . .
RUN BUILD_TAGS=muslc LINK_STATICALLY=true make build

# Add to a distroless container
FROM alpine:${RUNNER_IMAGE_VERSION}

COPY --from=builder /opt/build/intentod /usr/local/bin/intentod
RUN apk add bash vim sudo dasel jq curl \
    && addgroup -g 1000 intento \
    && adduser -S -h /home/intento -D intento -u 1000 -G intento 

RUN mkdir -p /etc/sudoers.d \
    && echo '%wheel ALL=(ALL) ALL' > /etc/sudoers.d/wheel \
    && echo "%wheel ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers \
    && adduser intento wheel 

USER 1000
ENV HOME=/home/intento
WORKDIR $HOME

EXPOSE 26657 26656 1317 9090

CMD ["intentod", "start"]