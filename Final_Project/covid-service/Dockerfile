# Stage 1: Build
FROM golang:1.24.4 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o covid-service .

# Stage 2: Run
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/covid-service .
CMD ["./covid-service"]
