# Stage 1: Build the Go binary
FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o kswp main.go

# Stage 2: Create the final image
FROM alpine:3.15
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
WORKDIR /app
COPY --from=builder /app/kswp .
ENTRYPOINT ["./kswp"]
