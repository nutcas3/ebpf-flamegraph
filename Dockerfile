FROM golang:1.26-bookworm AS builder

# Install build dependencies for eBPF
RUN apt-get update && apt-get install -y \
    clang \
    llvm \
    libbpf-dev \
    linux-headers-generic \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Install bpf2go tool
RUN go install github.com/cilium/ebpf/cmd/bpf2go@latest

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate eBPF bindings
RUN go generate ./...

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ebpf-demo .

# Production stage
FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -m -u 1000 ebpf && \
    mkdir -p /app/ebpf && \
    chown -R ebpf:ebpf /app

# Set working directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/ebpf-demo .

# Copy eBPF bytecode files
COPY --from=builder /app/bpf_*_*.o ./ebpf/

# Change ownership to ebpf user
RUN chown -R ebpf:ebpf /app

# Switch to non-root user
USER ebpf

# Expose metrics port (if needed)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep ebpf-demo > /dev/null || exit 1

# Default command - requires privileged access for eBPF
ENTRYPOINT ["./ebpf-demo"]
