#!/bin/bash
set -e

CLUSTER_NAME="query-service"

echo "==> Switching kubectl context to kind-${CLUSTER_NAME}"
kubectl config use-context kind-${CLUSTER_NAME}

echo "==> Checking pods are running"
kubectl get pods -l app=query-service

echo ""
echo "==> Port-forwarding service/query-service to localhost:8080"
echo "    Test with:"
echo "    curl -X POST http://localhost:8080/query \\"
echo "      -H 'Content-Type: application/json' \\"
echo "      -d '{\"query\": \"What is the capital of France?\"}'"
echo ""

kubectl port-forward service/query-service 8080:80