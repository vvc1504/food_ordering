# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod ./
# Skip go sum check for internal/custom modules if any
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
# Copy requirements folder (containing coupons)
COPY requirements ./requirements

# Expose port
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]
CMD ["-port", "8080", "-coupon-dir", "requirements"]
