FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy shared module first
COPY ../../shared ./shared

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy shared package if needed for runtime
COPY --from=builder /app/shared ./shared

EXPOSE 8082

CMD ["./main"]
