FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /repomind ./cmd/repomind

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /repomind /repomind

EXPOSE 3000

ENTRYPOINT ["/repomind"]