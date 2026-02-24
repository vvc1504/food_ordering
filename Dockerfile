# Frontend build stage
FROM node:20-alpine AS node-builder

WORKDIR /app/web
COPY web/package*.json ./
RUN npm install
COPY web/ ./
RUN npm run build

# Backend build stage
FROM golang:1.24-alpine AS go-builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Remove the web folder to avoid copying unnecessary files into the go build context
# although we already copied everything above.
RUN go build -o food_ordering ./cmd/server/main.go

# Production stage
FROM alpine:latest

# Add ca-certificates for any HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the backend binary
COPY --from=go-builder /app/food_ordering .

# Copy the frontend built assets
COPY --from=node-builder /app/web/dist ./web/dist

# Copy the coupon files
COPY requirements ./requirements

# Default environment variables
ENV PORT=8080
ENV COUPON_FILES="requirements/couponbase1.gz,requirements/couponbase2.gz,requirements/couponbase3.gz"

# Expose the API port
EXPOSE 8080

# Clean entrypoint
ENTRYPOINT ["./food_ordering"]

# Support both environment variables and direct flags
CMD ["-port", "8080"]
