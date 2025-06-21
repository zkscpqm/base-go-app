FROM golang:1.23 as builder

WORKDIR /app

COPY pkg/ pkg/
COPY cmd/ cmd/
COPY go.mod .
COPY go.sum .

RUN GOOS=linux GOARCH=amd64 go build -o unnamed ./cmd/unnamed/main.go

# Step 2: Use a minimal image for running the application
FROM debian:stable-slim

WORKDIR /app

RUN apt-get update && apt-get install -y curl

COPY deploy/dev/config.json .
COPY --from=builder /app/unnamed .

EXPOSE 8080
ENTRYPOINT ["./unnamed"]
