FROM golang:1.18-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /kswp

FROM scratch
COPY --from=builder /kswp /kswp
ENTRYPOINT ["/kswp"]
