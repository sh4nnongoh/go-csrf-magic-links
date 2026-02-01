# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
# -ldflags="-s -w" strips debug symbols, reducing binary size
# CGO_ENABLED=0 creates static binary, no C dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -a -o main .

# Stage 2: Create minimal runtime image
FROM scratch

# Copy only the built binary, and public assets from the builder stage
COPY --from=builder /app/main /main
COPY --from=builder /app/static /static
COPY --from=builder /app/templates /templates

# Copy CA certificates for HTTPS if needed
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Expose ports
EXPOSE 8080 6060

# Run as non-root user for security
USER 1000:1000

# Set minimal environment
ENV GOMAXPROCS=1
ENV GOMEMLIMIT=58MiB

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/main", "health"] || exit 1

ENTRYPOINT ["/main"]
