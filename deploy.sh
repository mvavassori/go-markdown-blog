#!/bin/bash
set -e

# Configuration
PROD_SERVER="tutorial-server"
LOCAL_TEMP_DIR="/tmp/docker-images"
# PROD_COMPOSE_DIR="/home/batman/go-markdown-blog"
# BUILD_COMPOSE_FILE="docker-compose.build.yaml"
# ENV_FILE=".env.prod"

# Create temp directory if it doesn't exist
mkdir -p ${LOCAL_TEMP_DIR}

# Function to build, save, and deploy a service
deploy_service() {
    local service=$1
    echo "Deploying $service..."
    
    # Build the image
    echo "Building $service image..."
    
    # # #
    # docker compose -f ${BUILD_COMPOSE_FILE} --env-file ${ENV_FILE} build ${service}
    # # #
    
    # For x86 cpus (assuming you're using a x86 cpu on your machine)
    # docker build -t ${service} .

    # For ARM cpus
    docker buildx build --platform linux/arm64 -t ${service} .
    
    
    # Get the source and target image names
    if [ "$service" == "go-markdown-blog" ]; then
        SOURCE_IMAGE="go-markdown-blog:latest"  # change if necessary
        TARGET_IMAGE="go-markdown-blog:latest"
    # elif [ "$service" == "frontend" ]; then
    #     SOURCE_IMAGE="frontend-prod:latest"
    #     TARGET_IMAGE="frontend:latest"
    # elif [ "$service" == "backend" ]; then
    #     SOURCE_IMAGE="backend-prod:latest"
    #     TARGET_IMAGE="backend:latest" 
    else
        echo "Unknown service: $service"
        exit 1
    fi
    
    # Tag the image to match production expectations
    echo "Tagging $SOURCE_IMAGE as $TARGET_IMAGE..."
    docker tag ${SOURCE_IMAGE} ${TARGET_IMAGE}
    
    # Save the image to a tar file
    echo "Saving $service image to tar file..."
    TAR_FILE="${LOCAL_TEMP_DIR}/${service}.tar"
    docker save ${TARGET_IMAGE} > ${TAR_FILE}
    
    # Transfer the image to production server
    echo "Transferring $service image to production server..."
    scp ${TAR_FILE} ${PROD_SERVER}:/tmp/
    
    # Load the image on the production server
    echo "Loading $service image on production server..."
    ssh ${PROD_SERVER} "docker load < /tmp/$(basename ${TAR_FILE}) && rm /tmp/$(basename ${TAR_FILE})"
    
    # Clean up local tar file
    rm ${TAR_FILE}
    
    echo "$service deployment preparation completed!"
}

# Show usage if no arguments provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <service1> [<service2> ...]"
    echo "Example: $0 frontend backend db-tasks"
    exit 1
fi

# Deploy each service specified as an argument
for service in "$@"; do
    deploy_service ${service}
done

# Restart services on production server
echo "Restarting services on production server..."

# # #
# ssh ${PROD_SERVER} "cd ${PROD_COMPOSE_DIR} && docker compose down && docker compose up -d"
# # #


# Just for this example
ssh ${PROD_SERVER} << 'EOF'
  # Check if the container exists and remove it if it does
  if docker ps -a --format '{{.Names}}' | grep -q "^temp-blog-test$"; then
    echo "Container temp-blog-test exists. Stopping it..."
    docker stop temp-blog-test
    
    echo "Removing container temp-blog-test..."
    docker rm temp-blog-test
    
    echo "Container removed successfully."
  else
    echo "Container temp-blog-test does not exist."
  fi
  
  # Create new container regardless of whether it previously existed
  echo "Creating new container temp-blog-test..."
  docker run -d --name temp-blog-test -p 80:8080 go-markdown-blog:latest
  echo "Container created successfully."
EOF

echo "Deployment completed successfully!"