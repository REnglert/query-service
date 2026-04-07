#!/bin/bash
set -e

IMAGE_NAME="query-service"
IMAGE_TAG="v1"
CLUSTER_NAME="query-service"

echo "==> Switching kubectl context to kind-${CLUSTER_NAME}"
kubectl config use-context kind-${CLUSTER_NAME}

echo "==> Building Docker image ${IMAGE_NAME}:${IMAGE_TAG}"
docker build -t ${IMAGE_NAME}:${IMAGE_TAG} .

echo "==> Loading image into kind cluster"
kind load docker-image ${IMAGE_NAME}:${IMAGE_TAG} --name ${CLUSTER_NAME}

echo "==> Applying Terraform"
cd terraform
terraform init -input=false
terraform apply -auto-approve
cd ..

echo "==> Waiting for rollout"
kubectl rollout status deployment/${IMAGE_NAME}

echo ""
echo "✓ Deploy complete. Run ./scripts/start.sh to begin port-forwarding."