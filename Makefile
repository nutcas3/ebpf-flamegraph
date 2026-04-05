.PHONY: help build run clean demo profile docker-build docker-run docker-stop docker-clean

# Default target
help:
	@echo "🚀 eBPF Demo Makefile"
	@echo "===================="
	@echo ""
	@echo "Build Targets:"
	@echo "  build          - Build eBPF applications"
	@echo "  docker-build   - Build Docker image with eBPF"
	@echo ""
	@echo "Run Targets:"
	@echo "  run            - Run eBPF monitor locally"
	@echo "  docker-run     - Run eBPF demo in Docker"
	@echo "  demo           - Start full demo with traffic"
	@echo ""
	@echo "Profiling Targets:"
	@echo "  profile        - Run performance profiling"
	@echo "  profile-cpu    - CPU profiling only"
	@echo "  profile-offcpu - Off-CPU profiling only"
	@echo "  profile-all    - All profiling types"
	@echo ""
	@echo "Management:"
	@echo "  clean          - Clean build artifacts"
	@echo "  docker-clean   - Clean Docker resources"
	@echo "  setup          - Setup development environment"
	@echo ""
	@echo "Examples:"
	@echo "  make docker-build && make docker-run"
	@echo "  make profile-all"
	@echo "  make demo"

# Variables
DOCKER_IMAGE = ebpf-demo
DOCKER_TAG = latest
CONTAINER_NAME = ebpf-demo-container
FLAMEGRAPH_DIR = $(HOME)/FlameGraph
OUTPUT_DIR = profiling-output

# Setup development environment
setup:
	@echo "🔧 Setting up eBPF development environment..."
	@mkdir -p $(OUTPUT_DIR)
	@if [ ! -d "$(FLAMEGRAPH_DIR)" ]; then \
		echo "📥 Cloning FlameGraph tools..."; \
		git clone https://github.com/brendangregg/FlameGraph.git $(FLAMEGRAPH_DIR); \
	fi
	@echo "✅ Environment setup complete"

# Build eBPF applications locally
build:
	@echo "🔨 Building eBPF applications..."
	@go mod tidy
	@go build -o ebpf-demo .
	@cd cmd/advanced && go build -o advanced-ebpf .
	@echo "✅ Build complete"

# Run eBPF monitor locally (requires Linux)
run:
	@echo "🚀 Starting eBPF monitor..."
	@if [ "$$(uname)" != "Linux" ]; then \
		echo "❌ eBPF requires Linux. Use 'make docker-run' for other platforms."; \
		exit 1; \
	fi
	@sudo ./ebpf-demo monitor

# Start advanced monitor with HTTP API
run-advanced:
	@echo "🚀 Starting advanced eBPF monitor..."
	@if [ "$$(uname)" != "Linux" ]; then \
		echo "❌ eBPF requires Linux. Use 'make docker-run' for other platforms."; \
		exit 1; \
	fi
	@sudo ./cmd/advanced/advanced-ebpf

# Build Docker image with eBPF support
docker-build:
	@echo "🐳 Building Docker image with eBPF support..."
	@docker build -f Dockerfile.ebpf -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✅ Docker image built"

# Run eBPF demo in Docker (works on all platforms)
docker-run:
	@echo "🚀 Starting eBPF demo in Docker..."
	@if ! docker image inspect $(DOCKER_IMAGE):$(DOCKER_TAG) >/dev/null 2>&1; then \
		echo "📦 Building Docker image first..."; \
		$(MAKE) docker-build; \
	fi
	@docker run --rm -it --name $(CONTAINER_NAME) \
		--privileged \
		--net=host \
		-v /sys/kernel/debug:/sys/kernel/debug:ro \
		-v /lib/modules:/lib/modules:ro \
		-v /usr/src:/usr/src:ro \
		-p 8080:8080 \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		./cmd/advanced/advanced-ebpf

# Run full demo with traffic generation
demo:
	@echo "🎭 Starting full eBPF demo with traffic..."
	@if ! docker image inspect $(DOCKER_IMAGE):$(DOCKER_TAG) >/dev/null 2>&1; then \
		echo "📦 Building Docker image first..."; \
		$(MAKE) docker-build; \
	fi
	@docker-compose -f docker-compose.ebpf.yml --profile demo up

# Stop Docker services
docker-stop:
	@echo "🛑 Stopping Docker services..."
	@docker-compose -f docker-compose.ebpf.yml down
	@docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@docker rm $(CONTAINER_NAME) 2>/dev/null || true

# Performance profiling targets
profile: profile-all

profile-cpu:
	@echo "🔥 Starting CPU profiling..."
	@$(MAKE) setup
	@./advanced-profiling-suite.sh cpu 30
	@$(MAKE) show-results

profile-offcpu:
	@echo "⏱️  Starting Off-CPU profiling..."
	@$(MAKE) setup
	@./advanced-profiling-suite.sh offcpu 30
	@$(MAKE) show-results

profile-memory:
	@echo "💾 Starting memory profiling..."
	@$(MAKE) setup
	@./advanced-profiling-suite.sh memory 30
	@$(MAKE) show-results

profile-cache:
	@echo "🗄️  Starting cache profiling..."
	@$(MAKE) setup
	@./advanced-profiling-suite.sh cache 30
	@$(MAKE) show-results

profile-all:
	@echo "🔥 Starting comprehensive profiling..."
	@$(MAKE) setup
	@./advanced-profiling-suite.sh all 60
	@$(MAKE) show-results
	@$(MAKE) analyze-results

# Show profiling results
show-results:
	@echo "📊 Opening profiling results..."
	@if command -v xdg-open >/dev/null 2>&1; then \
		xdg-open $(OUTPUT_DIR)/cpu-flamegraph.svg 2>/dev/null || true; \
	elif command -v open >/dev/null 2>&1; then \
		open $(OUTPUT_DIR)/cpu-flamegraph.svg 2>/dev/null || true; \
	else \
		echo "📂 Results available in: $(OUTPUT_DIR)"; \
		ls -la $(OUTPUT_DIR)/*.svg 2>/dev/null || echo "No SVG files found"; \
	fi

# Analyze profiling results
analyze-results:
	@echo "🤖 Analyzing profiling results..."
	@./advanced-profiling-suite.sh analyze

# Generate differential analysis
compare:
	@echo "📈 Generating differential analysis..."
	@if [ -z "$(BASELINE)" ] || [ -z "$(OPTIMIZED)" ]; then \
		echo "Usage: make compare BASELINE=file.folded OPTIMIZED=file.folded"; \
		exit 1; \
	fi
	@./advanced-profiling-suite.sh compare $(BASELINE) $(OPTIMIZED)
	@$(MAKE) show-results

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f ebpf-demo cmd/advanced/advanced-ebpf
	@rm -f *.o *.elf
	@rm -f bpf_*.go cmd/advanced/bpf_*.go
	@rm -rf $(OUTPUT_DIR)
	@echo "✅ Clean complete"

# Clean Docker resources
docker-clean:
	@echo "🧹 Cleaning Docker resources..."
	@docker-compose -f docker-compose.ebpf.yml down -v
	@docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@docker system prune -f
	@echo "✅ Docker clean complete"

# Full clean
clean-all: clean docker-clean

# Development targets
dev-setup: setup docker-build
	@echo "🚀 Development environment ready"

dev-run: docker-run
	@echo "🚀 Development server running"

dev-profile: docker-build
	@echo "🔥 Starting profiling in Docker..."
	@docker run --rm -it --name profiling-$(shell date +%s) \
		--privileged \
		--net=host \
		-v $(PWD):/app \
		-w /app \
		$(DOCKER_IMAGE):$(DOCKER_TAG) \
		make profile-all

# Test targets
test-build:
	@echo "🧪 Testing build..."
	@$(MAKE) clean
	@$(MAKE) build
	@echo "✅ Build test passed"

test-docker:
	@echo "🧪 Testing Docker build..."
	@$(MAKE) docker-clean
	@$(MAKE) docker-build
	@echo "✅ Docker test passed"

test-all: test-build test-docker
	@echo "✅ All tests passed"

# Status and information
status:
	@echo "📊 eBPF Demo Status"
	@echo "==================="
	@echo "Docker Image: $$(docker images -q $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || echo 'Not built')"
	@echo "Running Containers: $$(docker ps --filter 'name=$(CONTAINER_NAME)' --format 'table {{.Names}}\t{{.Status}}' 2>/dev/null | tail -n +2 || echo 'None')"
	@echo "Build Artifacts: $$(ls -la ebpf-demo cmd/advanced/advanced-ebpf 2>/dev/null | wc -l || echo '0')"
	@echo "Profiling Output: $$(ls -la $(OUTPUT_DIR) 2>/dev/null | wc -l || echo '0')"

# Quick commands for common tasks
quick-demo: docker-build docker-run
	@echo "🚀 Quick demo started"

quick-profile: docker-build dev-profile
	@echo "🔥 Quick profiling started"

# Help for technical talk
talk-help:
	@echo "🎤 Technical Talk Commands"
	@echo "========================="
	@echo "Pre-talk setup:"
	@echo "  make dev-setup      - Setup everything"
	@echo "  make test-all       - Test everything"
	@echo ""
	@echo "During talk:"
	@echo "  make quick-demo     - Start eBPF demo"
	@echo "  make quick-profile  - Start profiling"
	@echo "  make show-results   - Show flame graphs"
	@echo ""
	@echo "Post-talk cleanup:"
	@echo "  make clean-all      - Clean everything"

# Default target
.DEFAULT_GOAL := help
