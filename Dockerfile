# ====================================================================
# Stage 1: Build the Go binary
# ====================================================================

FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and compile statically
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.version={{.Version}}" -o /app/server .

# ====================================================================
# Stage 2: Create lightweight production image
# ====================================================================

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/server /server

# Use a non-root user for security
USER nonroot:nonroot

#EXPOSE 8443

ENTRYPOINT ["/server"]
