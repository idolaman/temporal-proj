# Deploy to Docker Desktop Kubernetes script (PowerShell)
Write-Host "Deploying crawler services to Docker Desktop Kubernetes..." -ForegroundColor Green

# Ensure we are using the docker-desktop context
Write-Host "Switching kubectl context to 'docker-desktop'..." -ForegroundColor Yellow
kubectl config use-context docker-desktop | Out-Null

# Build Docker images (Docker Desktop shares daemon with its Kubernetes)
Write-Host "Building service1 image..." -ForegroundColor Yellow
docker build -t crawler/service1:latest -f service1/Dockerfile .

Write-Host "Building service2 image..." -ForegroundColor Yellow
docker build -t crawler/service2:latest -f service2/Dockerfile .

# Deploy using Helm
Write-Host "Deploying with Helm..." -ForegroundColor Yellow
helm upgrade --install crawler ./helm-chart/

Write-Host "Deployment complete!" -ForegroundColor Green
Write-Host "" 
Write-Host "Service #2 (NodePort) should be reachable at:" -ForegroundColor Cyan
$nodePort = 30080
Write-Host "http://localhost:$nodePort" -ForegroundColor White 