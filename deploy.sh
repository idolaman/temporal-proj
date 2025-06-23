#!/bin/bash

# Deploy to Minikube script
echo "Deploying crawler services to Minikube..."

# Build Docker images with Minikube's Docker daemon
echo "Building Docker images..."
eval $(minikube docker-env)

# Build service1
echo "Building service1..."
docker build -t crawler/service1:latest ./service1/

# Build service2
echo "Building service2..."
docker build -t crawler/service2:latest ./service2/

# Deploy using Helm
echo "Deploying with Helm..."
helm upgrade --install crawler ./helm-chart/

echo "Deployment complete!"
echo ""
echo "To access service2 externally:"
echo "minikube service crawler-service2 --url"
echo ""
echo "Or get the NodePort URL:"
echo "echo http://$(minikube ip):30080" 