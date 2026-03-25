# Docker Deployment Guide

## Quick Start: Local Testing

### Build the Image
```bash
cd /home/emma/Downloads/quickbooks-online-mcp
docker build -t qbo-mcp-server:latest .
```

### Run Locally (Stdio Mode)
```bash
docker run -it \
  -e QUICKBOOKS_CLIENT_ID=test_client_id \
  -e QUICKBOOKS_CLIENT_SECRET=test_client_secret \
  -e QUICKBOOKS_REALM_ID=test_realm_id \
  -e QUICKBOOKS_REFRESH_TOKEN=test_refresh_token \
  -e QUICKBOOKS_ENVIRONMENT=sandbox \
  qbo-mcp-server:latest
```

### Run Locally (HTTP Mode)
```bash
docker run -d \
  -p 3000:3000 \
  -e PORT=3000 \
  -e MCP_TRANSPORT=http \
  -e QUICKBOOKS_CLIENT_ID=your_client_id \
  -e QUICKBOOKS_CLIENT_SECRET=your_client_secret \
  -e QUICKBOOKS_REALM_ID=your_realm_id \
  -e QUICKBOOKS_REFRESH_TOKEN=your_refresh_token \
  -e QUICKBOOKS_ENVIRONMENT=sandbox \
  qbo-mcp-server:latest
```

Test it:
```bash
curl http://localhost:3000/mcp -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'
```

---

## Using Docker Compose (Recommended)

### 1. Create `.env` file (local credentials)
```env
QUICKBOOKS_CLIENT_ID=your_client_id
QUICKBOOKS_CLIENT_SECRET=your_client_secret
QUICKBOOKS_REALM_ID=your_realm_id
QUICKBOOKS_REFRESH_TOKEN=your_refresh_token
QUICKBOOKS_ENVIRONMENT=sandbox
PORT=3000
MCP_TRANSPORT=http
```

### 2. Update `compose.yaml`
Located at: `/home/emma/Downloads/quickbooks-online-mcp/compose.yaml`

```yaml
services:
  qbo-mcp-server:
    build: .
    image: qbo-mcp-server:latest
    container_name: qbo-mcp
    env_file:
      - .env
    ports:
      - "3000:3000"
    environment:
      - PORT=3000
      - MCP_TRANSPORT=http
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/mcp"]
      interval: 30s
      timeout: 10s
      retries: 3
```

### 3. Run with Compose
```bash
# Start
docker compose up -d

# View logs
docker compose logs -f

# Stop
docker compose down

# Rebuild
docker compose up --build
```

---

## Production Deployment

### Environment Setup for Production

Create `.env.production`:
```env
QUICKBOOKS_CLIENT_ID=prod_client_id
QUICKBOOKS_CLIENT_SECRET=prod_client_secret
QUICKBOOKS_REALM_ID=prod_realm_id
QUICKBOOKS_REFRESH_TOKEN=prod_refresh_token
QUICKBOOKS_ENVIRONMENT=production
PORT=3000
MCP_TRANSPORT=http

# Optional: Organization mode
PLATFORM_INT_URL=https://api.yourcompany.com
ENFORCE_AUTH=true

# Optional: Monitoring
USAGE_REPORT_ENDPOINT=https://monitoring.yourcompany.com/usage
```

### Production Docker Compose
```yaml
version: '3.8'

services:
  qbo-mcp-server:
    build:
      context: .
      dockerfile: Dockerfile
    image: qbo-mcp-server:1.0.0
    container_name: qbo-mcp-prod
    env_file:
      - .env.production
    ports:
      - "3000:3000"
    restart: always
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/mcp"]
      interval: 60s
      timeout: 10s
      retries: 5
      start_period: 20s
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
    # Resource limits
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

---

## Docker Hub / Registry Deployment

### Push to Docker Hub
```bash
# Login
docker login

# Tag image
docker tag qbo-mcp-server:latest yourusername/qbo-mcp-server:latest
docker tag qbo-mcp-server:latest yourusername/qbo-mcp-server:1.0.0

# Push
docker push yourusername/qbo-mcp-server:latest
docker push yourusername/qbo-mcp-server:1.0.0
```

### Pull and Run from Docker Hub
```bash
docker run -d \
  -p 3000:3000 \
  --env-file .env.production \
  yourusername/qbo-mcp-server:latest
```

---

## Azure Container Registry (ACR) Deployment

### 1. Create ACR (if not exists)
```bash
az acr create \
  --resource-group myResourceGroup \
  --name myacrname \
  --sku Basic
```

### 2. Push to ACR
```bash
# Login to ACR
az acr login --name myacrname

# Build and push
az acr build \
  --registry myacrname \
  --image qbo-mcp-server:latest .

# Or manually:
docker tag qbo-mcp-server:latest myacrname.azurecr.io/qbo-mcp-server:latest
docker push myacrname.azurecr.io/qbo-mcp-server:latest
```

### 3. Deploy to Azure Container Instances
```bash
az container create \
  --resource-group myResourceGroup \
  --name qbo-mcp-container \
  --image myacrname.azurecr.io/qbo-mcp-server:latest \
  --registry-login-server myacrname.azurecr.io \
  --registry-username <username> \
  --registry-password <password> \
  --environment-variables \
    QUICKBOOKS_CLIENT_ID=your_id \
    QUICKBOOKS_CLIENT_SECRET=your_secret \
    QUICKBOOKS_REALM_ID=your_realm \
    QUICKBOOKS_REFRESH_TOKEN=your_token \
    PORT=3000 \
    MCP_TRANSPORT=http \
  --ports 3000 \
  --dns-name-label qbo-mcp
```

---

## Kubernetes Deployment (Advanced)

### Create ConfigMap for configuration
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: qbo-mcp-config
data:
  PORT: "3000"
  MCP_TRANSPORT: "http"
  QUICKBOOKS_ENVIRONMENT: "sandbox"
```

### Create Secret for credentials
```bash
kubectl create secret generic qbo-mcp-secrets \
  --from-literal=QUICKBOOKS_CLIENT_ID=your_id \
  --from-literal=QUICKBOOKS_CLIENT_SECRET=your_secret \
  --from-literal=QUICKBOOKS_REALM_ID=your_realm \
  --from-literal=QUICKBOOKS_REFRESH_TOKEN=your_token
```

### Deploy
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: qbo-mcp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: qbo-mcp
  template:
    metadata:
      labels:
        app: qbo-mcp
    spec:
      containers:
      - name: qbo-mcp
        image: myregistry.azurecr.io/qbo-mcp-server:latest
        ports:
        - containerPort: 3000
        envFrom:
        - configMapRef:
            name: qbo-mcp-config
        - secretRef:
            name: qbo-mcp-secrets
        livenessProbe:
          httpGet:
            path: /mcp
            port: 3000
          initialDelaySeconds: 10
          periodSeconds: 30
        readinessProbe:
          httpGet:
            path: /mcp
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 10
        resources:
          requests:
            memory: "256Mi"
            cpu: "500m"
          limits:
            memory: "512Mi"
            cpu: "1000m"
---
apiVersion: v1
kind: Service
metadata:
  name: qbo-mcp-service
spec:
  selector:
    app: qbo-mcp
  type: LoadBalancer
  ports:
  - protocol: TCP
    port: 3000
    targetPort: 3000
```

Deploy:
```bash
kubectl apply -f k8s-configmap.yaml
kubectl apply -f k8s-secret.yaml
kubectl apply -f k8s-deployment.yaml
```

---

## Monitoring & Logging

### View Logs
```bash
# Docker container
docker logs -f qbo-mcp

# Docker compose
docker compose logs -f qbo-mcp-server

# Kubernetes
kubectl logs -f deployment/qbo-mcp
```

### Health Check Endpoint
The MCP server automatically responds to health checks on:
```
GET/POST http://localhost:3000/mcp
```

Any request to `/mcp` with proper MCP headers will succeed if the server is healthy.

---

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Container exits immediately | Check logs: `docker logs <container>` |
| Port already in use | Change PORT env var or use `docker run -p 3001:3000` |
| Authentication fails | Verify credentials in `.env` |
| Container fails health check | Ensure QB credentials are valid |
| Out of memory | Increase memory limit in docker-compose.yaml |

---

## Security Best Practices

1. **Don't hardcode credentials** in Dockerfile
2. **Use .env files** for local testing only
3. **Use secrets management** in production (Docker Secrets, K8s Secrets, Vault)
4. **Use private registries** for container images
5. **Enable HTTPS** with reverse proxy (nginx, Traefik)
6. **Limit network access** with firewalls
7. **Regular updates** - rebuild images monthly
8. **Image scanning** - scan for vulnerabilities before deployment

## Next Steps

✅ Build locally: `docker build -t qbo-mcp-server .`
✅ Test locally: `docker compose up`
✅ Get QB credentials: See QUICKBOOKS_CREDENTIALS.md
✅ Deploy to production: Choose your platform above
