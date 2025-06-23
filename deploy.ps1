# Deploy to Minikube script (PowerShell)
Write-Host "Deploying crawler services to Minikube..." -ForegroundColor Green

# Build Docker images with Minikube's Docker daemon
Write-Host "Setting up Minikube Docker environment..." -ForegroundColor Yellow
& minikube docker-env --shell powershell | Invoke-Expression

# Build service1
Write-Host "Building service1..." -ForegroundColor Yellow
docker build -t crawler/service1:latest ./service1/

# Build service2
Write-Host "Building service2..." -ForegroundColor Yellow
docker build -t crawler/service2:latest ./service2/

# Deploy using Helm
Write-Host "Deploying with Helm..." -ForegroundColor Yellow
helm upgrade --install crawler ./helm-chart/

Write-Host "Deployment complete!" -ForegroundColor Green
Write-Host ""
Write-Host "To access service2 externally:" -ForegroundColor Cyan
Write-Host "minikube service crawler-service2 --url" -ForegroundColor White
Write-Host ""
Write-Host "Or get the NodePort URL:" -ForegroundColor Cyan
$minikubeIP = "localhost"
Write-Host "http://${minikubeIP}:30080" -ForegroundColor White 