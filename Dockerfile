# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Install git for go modules
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Move to working dir for main.go
WORKDIR /app/cmd/server

# Build Go binary
RUN go build -o server .

# Run stage
FROM debian:stable-slim

# Install certificates (for HTTPS support)
RUN apt-get update && apt-get install -y ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/cmd/server/server .

# (Optional) Copy configs if needed
COPY configs/ /root/configs/

# Expose port (เปลี่ยนตามที่ Go app ใช้)
EXPOSE 8080

# Run the binary
CMD ["./server"]