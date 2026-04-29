#!/bin/bash

# eBPF Demo Setup Script
# This script sets up the environment for running eBPF demos

set -e

echo "🚀 Setting up eBPF Demo Environment..."

# Check if running on Linux
if [[ "$OSTYPE" != "linux-gnu"* ]]; then
    echo "❌ This demo requires Linux. Current OS: $OSTYPE"
    echo "💡 Use Docker with Linux container or run on a Linux system"
    exit 1
fi

# Check kernel version
KERNEL_VERSION=$(uname -r | cut -d. -f1-2)
REQUIRED_VERSION="5.8"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$KERNEL_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "⚠️  Warning: Kernel version $KERNEL_VERSION is older than recommended $REQUIRED_VERSION"
    echo "   Some eBPF features may not work properly"
fi

# Check BTF support
if [ -f /sys/kernel/btf/vmlinux ]; then
    echo "✅ BTF (BPF Type Format) is enabled"
else
    echo "❌ BTF is not enabled. Some CO-RE features may not work"
    echo "   Consider enabling CONFIG_DEBUG_INFO_BTF=y in kernel config"
fi

# Install dependencies
echo "📦 Installing dependencies..."

# Update package list
sudo apt-get update

# Install build tools
sudo apt-get install -y \
    build-essential \
    clang \
    llvm \
    libbpf-dev \
    linux-headers-$(uname -r) \
    pkg-config \
    git

# Install Go if not present
if ! command -v go &> /dev/null; then
    echo "📦 Installing Go..."
    GO_VERSION="1.22.0"
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
    rm go${GO_VERSION}.linux-amd64.tar.gz
    
    # Add Go to PATH
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    export PATH=$PATH:/usr/local/go/bin
fi

# Install eBPF tools
echo "🔧 Installing eBPF tools..."
go install github.com/cilium/ebpf/cmd/bpf2go@latest
go install github.com/cilium/ebpf/cmd/bpf2go@latest

# Install additional tools
sudo apt-get install -y \
    bpfcc-tools \
    bpftrace \
    linux-tools-$(uname -r) \
    linux-cloud-tools-$(uname -r)

# Set up permissions
echo "🔒 Setting up permissions..."
sudo setcap cap_sys_admin+ep $(which go) 2>/dev/null || true

# Build the demo
echo "🏗️  Building the eBPF demo..."
cd "$(dirname "$0")"

# Initialize Go modules
if [ ! -f go.mod ]; then
    go mod init ebpf-demo
fi

# Download dependencies
go mod tidy

# Generate eBPF bindings
echo "⚡ Generating eBPF bindings..."
go generate ./...

# Build the application
echo "🔨 Building application..."
go build -o ebpf-demo .

# Build advanced demo
echo "🔨 Building advanced demo..."
cd cmd/advanced
go build -o advanced-ebpf .
cd ../..

# Create Docker image
echo "🐳 Building Docker image..."
docker build -t ebpf-demo .

# Test installation
echo "🧪 Testing installation..."

# Test basic eBPF functionality
if command -v bpftrace &> /dev/null; then
    echo "✅ bpftrace is available"
    sudo bpftrace -l 'tracepoint:syscalls:*' | head -5 > /dev/null
    echo "✅ bpftrace can access tracepoints"
fi

# Test kernel module loading
if lsmod | grep -q "bpf"; then
    echo "✅ eBPF kernel module is loaded"
else
    echo "⚠️  eBPF kernel module may not be loaded"
fi

# Create demo scripts
echo "📝 Creating demo scripts..."

# Basic monitoring script
cat > run-monitor.sh << 'EOF'
#!/bin/bash
echo "🚀 Starting Basic eBPF Monitor..."
echo "📊 This will show network statistics in real-time"
echo "Press Ctrl+C to stop"
echo ""
sudo ./ebpf-demo monitor
EOF

# Security monitoring script
cat > run-security.sh << 'EOF'
#!/bin/bash
echo "🔒 Starting eBPF Security Monitor..."
echo "🛡️  This will monitor for security threats and block attacks"
echo "Press Ctrl+C to stop"
echo ""
sudo ./ebpf-demo security
EOF

# Advanced monitoring script
cat > run-advanced.sh << 'EOF'
#!/bin/bash
echo "🚀 Starting Advanced eBPF Monitor..."
echo "📊 This will start the advanced monitor with HTTP API"
echo "🌐 Metrics will be available at http://localhost:8080/metrics"
echo "Press Ctrl+C to stop"
echo ""
sudo ./cmd/advanced/advanced-ebpf
EOF

# Docker demo script
cat > run-docker.sh << 'EOF'
#!/bin/bash
echo "🐳 Starting eBPF Demo in Docker..."
echo "📊 This will run the monitor in a Docker container"
echo ""
docker-compose up ebpf-monitor
EOF

# Make scripts executable
chmod +x run-*.sh

# Create traffic generation script
cat > generate-traffic.sh << 'EOF'
#!/bin/bash
echo "🌊 Generating traffic for demo..."
echo "This will create various types of network traffic"
echo ""

# HTTP traffic
echo "📡 Generating HTTP traffic..."
for i in {1..10}; do
    curl -s http://httpbin.org/get > /dev/null 2>&1 || true
    sleep 0.1
done

# DNS queries
echo "🔍 Generating DNS queries..."
for domain in google.com facebook.com github.com; do
    nslookup $domain > /dev/null 2>&1 || true
    sleep 0.1
done

# Pings
echo "🏓 Generating ping traffic..."
for target in 8.8.8.8 1.1.1.1; do
    ping -c 1 $target > /dev/null 2>&1 || true
    sleep 0.1
done

echo "✅ Traffic generation complete!"
EOF

chmod +x generate-traffic.sh

echo ""
echo "🎉 Setup complete!"
echo ""
echo "📋 Available commands:"
echo "  ./run-monitor.sh      - Start basic network monitoring"
echo "  ./run-security.sh     - Start security monitoring"
echo "  ./run-advanced.sh     - Start advanced monitoring with API"
echo "  ./run-docker.sh       - Run in Docker container"
echo "  ./generate-traffic.sh - Generate test traffic"
echo ""
echo "🌐 Advanced monitor API endpoints:"
echo "  http://localhost:8080/metrics        - Full metrics"
echo "  http://localhost:8080/top-connections - Top connections"
echo "  http://localhost:8080/security-events - Security events"
echo "  http://localhost:8080/health          - Health check"
echo ""
echo "⚠️  Note: All eBPF programs require sudo privileges"
echo "💡 Tip: Run ./generate-traffic.sh in another terminal to see activity"
echo ""
echo "🚀 Ready to start the eBPF demo!"
