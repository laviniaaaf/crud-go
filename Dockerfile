FROM golang:1.25-alpine AS builder


WORKDIR /app

# Instala git e certificados
RUN apk add --no-cache git ca-certificates

# Copia go.mod e go.sum da raiz
COPY go.mod go.sum ./
RUN go mod download

# Copia todo o backend
COPY ./backend ./backend

# Compila o aplicativo
RUN go build -o main ./backend

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache ca-certificates

# Copia binário compilado
COPY --from=builder /app/main .

# Copia os arquivos estáticos do frontend para serem servidos pelo backend
COPY ./frontend ./frontend

EXPOSE 8080
CMD ["./main"]

#
#docker compose ps
# para o conteiner = docker compose down