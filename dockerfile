# Build stage
FROM golang:1.27-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o lavendera-API .

# Run stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/lavendera-API .

EXPOSE 8000

CMD ["./lavendera-API"]