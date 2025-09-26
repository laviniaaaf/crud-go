FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source files
COPY . .

# Compile the application
RUN go build -o /main ./backend

FROM alpine:latest
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /main /app/main

# Copy frontend files
COPY ./frontend ./frontend

# Expose port
EXPOSE 8080
CMD ["./main"]
# docker compose up --build
# docker compose ps
#  stop conteiner = docker compose down
# docker compose up -d --remove-orphans --> up new projct and to clean the containers that nnot use more