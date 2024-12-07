FROM golang:1.22.0-alpine AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates gcc libc-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o hacko-app ./cmd/bin/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/hacko-app .
EXPOSE 3000
CMD ["./hacko-app"]
