#!/bin/bash
set -e

# Requires DOCKERHUB_USERNAME to be set in the environment
# e.g. export DOCKERHUB_USERNAME=credemol
if [[ -z "${DOCKERHUB_USERNAME}" ]]; then
  echo "ERROR: DOCKERHUB_USERNAME is not set. Run: export DOCKERHUB_USERNAME=<your-dockerhub-username>"
  exit 1
fi

# Define the image tag (use semantic versioning for production)
IMAGE_REPOSITORY=appconfig-operator
IMAGE_TAG=${IMAGE_TAG:-0.2.0}

echo "Image Tag: $IMAGE_TAG"

# Build and push linux/amd64 (for cloud Kubernetes clusters)
docker buildx create --use

# in appconfig-operator directory.
#docker buildx build --push --platform linux/amd64,linux/arm64 \
#  -t $DOCKERHUB_USERNAME/$IMAGE_REPOSITORY:$IMAGE_TAG \
#  -f ./Dockerfile .

docker buildx build --platform linux/amd64,linux/arm64 \
  -t $DOCKERHUB_USERNAME/$IMAGE_REPOSITORY:$IMAGE_TAG \
  ./appconfig-operator --push


echo "Done. Pushed:"
