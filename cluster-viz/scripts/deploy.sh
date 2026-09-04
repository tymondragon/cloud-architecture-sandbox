#!/bin/bash
set -e

echo "Deploying cluster-viz to Kubernetes..."

# Apply all manifests
kubectl apply -f manifests/

echo "Waiting for deployment to be ready..."

# Wait for deployment
kubectl wait --for=condition=available --timeout=60s deployment/cluster-viz -n cluster-viz

echo "Deployment complete!"
echo ""
echo "Access the UI:"
echo "  Port-forward: kubectl port-forward -n cluster-viz svc/cluster-viz 9090:80"
echo "  Then open: http://localhost:9090"
echo ""
echo "Or add to /etc/hosts:"
echo "  127.0.0.1 cluster-viz.local"
echo "  Then open: http://cluster-viz.local (requires port 80 forwarding)"
