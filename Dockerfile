# Compile
FROM golang:alpine AS builder
WORKDIR /src/app/
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN for bin in cmd/*; do CGO_ENABLED=1 go build -o=/usr/local/bin/$(basename $bin) ./cmd/$(basename $bin); done


# Add to a  codistrolessntainer
FROM gcr.io/distroless/base
COPY --from=builder /usr/local/bin /usr/local/bin
USER nonroot:nonroot
CMD ["trstd start"]
