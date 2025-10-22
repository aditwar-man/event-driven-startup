FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy entire project
COPY . .

# Build using workspace
RUN go work sync
RUN cd services/auth && go build -o /app/main ./cmd

EXPOSE 8081

CMD ["/app/main"]