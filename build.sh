# !/bin/bash

set -e

REGION="us-west-1"
CLUSTER_NAME="rtqa7"
PROD=true
SECRET_NAME=www-api

SERVICE_NAME="www-api"
ECR_REPO_NAME=$CLUSTER_NAME-$SERVICE_NAME
ECR_REPO_BASE_URI="280563394466.dkr.ecr."$REGION".amazonaws.com"
NEW_VERSION="latest"
OLD_VERSION="last"

echo "Building $SERVICE_NAME version $NEW_VERSION..."

aws ecr get-login-password --region $REGION | docker login --username AWS --password-stdin $ECR_REPO_BASE_URI

aws ecr batch-delete-image --region $REGION --repository-name $ECR_REPO_NAME --image-ids imageTag=$OLD_VERSION

MY_MANIFEST=$(aws ecr batch-get-image --repository-name $ECR_REPO_NAME --image-ids imageTag=latest --region $REGION --query images[].imageManifest --output text)

if [ -n "$MY_MANIFEST" ]; then
aws ecr put-image --repository-name $ECR_REPO_NAME --image-tag $OLD_VERSION --image-manifest "$MY_MANIFEST" --region $REGION > /dev/null
fi

if $PROD; then
docker build -f prod.dockerfile --build-arg region=$REGION --build-arg secret_name=$SECRET_NAME -t $SERVICE_NAME:$NEW_VERSION .
else
docker build -f dev.dockerfile --build-arg region=$REGION --build-arg secret_name=$SECRET_NAME -t $SERVICE_NAME:$NEW_VERSION .
fi

docker tag $SERVICE_NAME:$NEW_VERSION $ECR_REPO_BASE_URI/$ECR_REPO_NAME:$NEW_VERSION

docker push $ECR_REPO_BASE_URI/$ECR_REPO_NAME:$NEW_VERSION

aws ecs update-service --region $REGION --cluster $CLUSTER_NAME --service $ECR_REPO_NAME --force-new-deployment > /dev/null

echo "Build & pushed $SERVICE_NAME version $NEW_VERSION, deploying the same on ecs $ECR_REPO_NAME service of $CLUSTER_NAME cluster"

