# Build stage
FROM golang:1.24 AS builder
WORKDIR /app
COPY . .
COPY data.csv ./data.csv
RUN go build -o server .

# Runtime stage (minimal)
FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/server /server
COPY --from=builder /app/data.csv /data.csv 
CMD ["/server"]