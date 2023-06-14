ARG IMG_TAG=latest

# Compile the trstd binary
FROM golang:1.20-alpine AS trstd-builder
WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN rm -rf /src/app/go-cosmwasm
RUN rm -rf /src/app/x/registration
RUN rm -rf /src/app/x/compute
RUN rm -rf /src/app/x/item
RUN rm -rf /src/app/cmd/trstd/attestation.go
RUN rm -rf /src/app/cmd/trstd/genwasm.go

ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES
RUN CGO_ENABLED=0 make install

# Add to a distroless container
FROM cgr.dev/chainguard/static:$IMG_TAG
ARG IMG_TAG
COPY --from=trstd-builder /go/bin/trstd /usr/local/bin/
COPY --from=trstd-builder --chown=0:0 /src/app/ /src/app/

EXPOSE 26656 26657 1317 9090
USER 0

ENTRYPOINT ["trstd", "start"]
