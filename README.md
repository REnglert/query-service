# query-service

A Go-based microservice that accepts natural-language queries, routes them through a local LLM (TinyLlama via llama.cpp), and returns structured responses. Built to demonstrate production-grade backend practices including clean API design, Kubernetes deployment, infrastructure-as-code with Terraform, graceful shutdown, and LLM integration.

## Architecture
```
curl /query → Go service (port 8080) → llama-server (port 8081) → TinyLlama 1.1B
                     ↓
              Kubernetes (kind)
              Managed by Terraform
```

## Prerequisites

- [Go 1.25+](https://golang.org/dl/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [kind](https://kind.sigs.k8s.io/) — `brew install kind`
- [kubectl](https://kubernetes.io/docs/tasks/tools/) — `brew install kubectl`
- [Terraform](https://www.terraform.io/) — `brew install hashicorp/tap/terraform`
- [llama.cpp](https://github.com/ggerganov/llama.cpp) built locally with `llama-server`
- TinyLlama 1.1B model (`tinyllama-1.1b-chat-v1.0.Q4_0.gguf`)

## Quick Start

### 1. Start the LLM server
```bash
./llama.cpp/build/bin/llama-server \
  -m ~/Documents/Projects/llama/tinyllama-1.1b-chat-v1.0.Q4_0.gguf \
  -c 2048 \
  --port 8081 \
  --host 0.0.0.0
```

### 2. Deploy everything
```bash
./scripts/deploy.sh
```

### 3. Start port-forwarding
```bash
./scripts/start.sh
```

### 4. Test it
```bash
curl -X POST http://localhost:8080/query \
  -H "Content-Type: application/json" \
  -d '{"query": "What is the capital of France?"}'
```

Expected response:
```json
{"result":"The capital of France is Paris."}
```

## API

| Endpoint  | Method | Purpose |
| --------- | ------ | ------- |
| `/health` | GET    | Liveness check. Returns `ok`. |
| `/ready`  | GET    | Readiness check. Returns `ready`. |
| `/query`  | POST   | Accepts `{"query": "..."}`, returns `{"result": "..."}`. |

## Project Structure
```
query-service/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── api/                    # HTTP handlers and routes
│   ├── config/                 # Environment-based config
│   ├── llm/                    # LLM client wrapper
│   ├── server/                 # HTTP server setup
│   └── shutdown/               # Graceful shutdown
├── k8s/                        # Raw Kubernetes manifests (reference)
├── terraform/                  # Terraform-managed infra
├── scripts/                    # Automation scripts
└── Dockerfile                  # Multi-stage build
```

## Useful Commands

### Local development
```bash
# Run the service locally (no Kubernetes)
go run ./cmd/server/main.go

# Run tests
go test ./...

# Build binary
go build -o bin/server ./cmd/server/main.go
```

### Docker
```bash
# Build image
docker build -t query-service:v1 .

# Run container locally
docker run -p 8080:8080 -e LLM_BASE_URL=http://host.docker.internal:8081 query-service:v1
```

### Kubernetes
```bash
# Check pod status
kubectl get pods

# Check logs
kubectl logs -l app=query-service

# Exec into a pod
kubectl exec -it $(kubectl get pods -o name | head -1) -- sh

# Port-forward the service
kubectl port-forward service/query-service 8080:80

# Restart deployment (e.g. after new image load)
kubectl rollout restart deployment/query-service

# Watch rollout status
kubectl rollout status deployment/query-service
```

### Terraform
```bash
# Preview changes
terraform -chdir=terraform plan

# Apply changes
terraform -chdir=terraform apply

# Tear down
terraform -chdir=terraform destroy
```

### kind
```bash
# List clusters
kind get clusters

# Load a new image into the cluster
kind load docker-image query-service:v1 --name query-service

# Delete and recreate the cluster
kind delete cluster --name query-service
kind create cluster --name query-service
```

## Configuration

| Env Var | Default | Description |
| ------- | ------- | ----------- |
| `PORT` | `8080` | Port the Go service listens on |
| `LLM_BASE_URL` | `http://localhost:8081` | Base URL for the llama-server |

## Design Decisions

- **Two-stage Docker build** — keeps the final image small (~26MB), only contains the binary
- **kind over Docker Desktop K8s** — more reliable local image loading via `kind load`
- **Terraform over raw kubectl** — infra expressed as reviewable, versioned code
- **Configurable LLM URL** — allows the same image to run locally or in Kubernetes without rebuilding
- **Graceful shutdown** — handles SIGTERM cleanly, drains in-flight requests before exit