# Build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server .

# Runtime stage (minimal)
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/server /server
CMD ["/server"]