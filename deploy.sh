#!/usr/bin/env bash
# Deploy to Docker Desktop Kubernetes (macOS/Linux)
set -euo pipefail

GREEN="$(tput setaf 2)"
YELLOW="$(tput setaf 3)"
CYAN="$(tput setaf 6)"
RESET="$(tput sgr0)"

info()  { printf "%s\n" "${YELLOW}$*${RESET}"; }
success() { printf "%s\n" "${GREEN}$*${RESET}"; }

success "Deploying crawler services to Docker Desktop Kubernetes..."

# Switch kubectl context to docker-desktop
info "Switching kubectl context to 'docker-desktop'..."
kubectl config use-context docker-desktop >/dev/null

# Change to the Go module directory
cd temporal-proj

# Build Docker images (Docker Desktop shares daemon with its Kubernetes cluster)
info "Building broker image..."
docker build -t crawler/broker:latest -f cmd/broker/Dockerfile .

info "Building broker-manager image..."
docker build -t crawler/broker-manager:latest -f cmd/broker-manager/Dockerfile .

# Change back to root directory for Helm deployment
cd ..

# Deploy using Helm
info "Deploying with Helm..."
helm upgrade --install crawler ./helm-chart/

success "Deployment complete!"

NODE_PORT=30080
printf "\n${CYAN}Broker Manager (NodePort) should be reachable at:${RESET}\n"
printf "http://localhost:%s\n" "$NODE_PORT" 