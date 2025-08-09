# Step 1: build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o app ./main.go

# Step 2: run stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 3000
CMD ["./app"]
