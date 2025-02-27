FROM golang:1.22-alpine as builder
RUN apk update && apk add gcc musl-dev

WORKDIR /usr/src/app

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY cmd/producer/ cmd/producer/
COPY internal/cli/ internal/cli/
RUN CGO_ENABLED=1 go build -o producer -tags musl ./cmd/producer/main.go

FROM alpine

COPY --from=builder /usr/src/app/producer /producer

ENTRYPOINT ["/producer"]
