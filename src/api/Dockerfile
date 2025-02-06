# Build stage
FROM golang:1.23.2-alpine3.20 AS builder

# Set the current building directory inside the container
WORKDIR /build

# Copy files into the container
COPY . .

# Download dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Intermediate stage
FROM alpine:3.14 AS intermediate

# Install curl
RUN apk add --no-cache curl

# Final stage
FROM scratch

# Copy the binary from the build stage
COPY --from=builder /build/main /main

# Copy the curl binary from the intermediate stage
COPY --from=intermediate /usr/bin/curl /usr/bin/curl
COPY --from=intermediate /lib /lib
COPY --from=intermediate /usr/lib /usr/lib

# Expose the port
EXPOSE 5000

# Run the binary
CMD ["/main"]
