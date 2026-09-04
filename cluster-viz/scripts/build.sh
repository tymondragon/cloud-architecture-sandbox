#!/bin/bash
set -e

echo "Building cluster-viz Docker image..."

# Build the Docker image
docker build -t cluster-viz:latest .

echo "Loading image into Kind cluster..."

# Load the image into Kind
kind load docker-image cluster-viz:latest

echo "Build complete! Image loaded into Kind cluster."
