# Build Stage
FROM golang:1.23 AS build

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application (static linking)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o xm-microservice ./cmd/server

# Final Stage
FROM alpine:latest

WORKDIR /root/

# Copy the built binary
COPY --from=build /app/xm-microservice .
COPY --from=build /app/internal ./internal

# Ensure binary has execute permissions
RUN chmod +x xm-microservice

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./xm-microservice"]
