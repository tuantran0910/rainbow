#!/bin/bash

# Setup VPC Network for Rainbow Data Production
# This script creates the VPC network, subnet, and related networking components

set -e

# Configuration
PROJECT_ID="rainbow-data-production"
REGION="us-central1"
ZONE="us-central1-a"
VPC_NAME="${PROJECT_ID}-vpc"
SUBNET_NAME="${PROJECT_ID}-subnet"
ROUTER_NAME="${PROJECT_ID}-router"
NAT_NAME="${PROJECT_ID}-nat"
PRIVATE_IP_NAME="${PROJECT_ID}-private-ip"

echo "🚀 Setting up VPC Network for Rainbow Data Production"
echo "Project ID: ${PROJECT_ID}"
echo "Region: ${REGION}"
echo "VPC Name: ${VPC_NAME}"

# Set the project
echo "📋 Setting project context..."
gcloud config set project "${PROJECT_ID}"

# Enable required APIs
echo "📡 Enabling required APIs..."
REQUIRED_APIS=(
  "compute.googleapis.com"
  "servicenetworking.googleapis.com"
  "vpcaccess.googleapis.com"
)

for API in "${REQUIRED_APIS[@]}"; do
  echo "  Enabling ${API}..."
  gcloud services enable "${API}" --project="${PROJECT_ID}"
done

echo "⏳ Waiting for APIs to be fully enabled..."
sleep 10

# Create VPC Network
echo "🌐 Creating VPC network..."
gcloud compute networks create "${VPC_NAME}" \
  --subnet-mode=custom \
  --bgp-routing-mode=regional \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  VPC network already exists"

# Create Subnet
echo "🔗 Creating subnet..."
gcloud compute networks subnets create "${SUBNET_NAME}" \
  --network="${VPC_NAME}" \
  --range="10.0.0.0/16" \
  --region="${REGION}" \
  --secondary-range="pods=10.1.0.0/16,services=10.2.0.0/16" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Subnet already exists"

# Create private IP address for Cloud SQL
echo "🔒 Creating private IP address range for Cloud SQL..."
gcloud compute addresses create "${PRIVATE_IP_NAME}" \
  --global \
  --purpose=VPC_PEERING \
  --prefix-length=16 \
  --network="${VPC_NAME}" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Private IP address already exists"

# Create service networking connection for Cloud SQL
echo "🔗 Creating service networking connection for Cloud SQL..."
gcloud services vpc-peerings connect \
  --service=servicenetworking.googleapis.com \
  --ranges="${PRIVATE_IP_NAME}" \
  --network="${VPC_NAME}" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Service networking connection already exists"

# Create Cloud Router for NAT
echo "🚦 Creating Cloud Router..."
gcloud compute routers create "${ROUTER_NAME}" \
  --network="${VPC_NAME}" \
  --region="${REGION}" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Cloud Router already exists"

# Create Cloud NAT
echo "🌍 Creating Cloud NAT..."
gcloud compute routers nats create "${NAT_NAME}" \
  --router="${ROUTER_NAME}" \
  --region="${REGION}" \
  --nat-all-subnet-ip-ranges \
  --auto-allocate-nat-external-ips \
  --enable-logging \
  --log-filter=ERRORS_ONLY \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Cloud NAT already exists"

# Create firewall rules
echo "🔥 Creating firewall rules..."

# Allow internal communication
gcloud compute firewall-rules create "${VPC_NAME}-allow-internal" \
  --network="${VPC_NAME}" \
  --allow=tcp,udp,icmp \
  --source-ranges="10.0.0.0/8" \
  --description="Allow internal communication within VPC" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Internal firewall rule already exists"

# Allow SSH
gcloud compute firewall-rules create "${VPC_NAME}-allow-ssh" \
  --network="${VPC_NAME}" \
  --allow=tcp:22 \
  --source-ranges="0.0.0.0/0" \
  --description="Allow SSH access" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  SSH firewall rule already exists"

# Allow HTTPS
gcloud compute firewall-rules create "${VPC_NAME}-allow-https" \
  --network="${VPC_NAME}" \
  --allow=tcp:443 \
  --source-ranges="0.0.0.0/0" \
  --description="Allow HTTPS access" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  HTTPS firewall rule already exists"

# Allow HTTP (for health checks)
gcloud compute firewall-rules create "${VPC_NAME}-allow-http" \
  --network="${VPC_NAME}" \
  --allow=tcp:80 \
  --source-ranges="0.0.0.0/0" \
  --description="Allow HTTP access for health checks" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  HTTP firewall rule already exists"

echo ""
echo "✅ VPC Network setup completed!"
echo ""
echo "📋 Created resources:"
echo "   🌐 VPC Network: ${VPC_NAME}"
echo "   🔗 Subnet: ${SUBNET_NAME} (10.0.0.0/16)"
echo "   🚦 Router: ${ROUTER_NAME}"
echo "   🌍 NAT Gateway: ${NAT_NAME}"
echo "   🔒 Private IP Range: ${PRIVATE_IP_NAME}"
echo ""
echo "🔗 Secondary IP ranges for GKE:"
echo "   📦 Pods: 10.1.0.0/16"
echo "   🔧 Services: 10.2.0.0/16"
echo ""
echo "🚀 Your VPC is now ready for Cloud SQL and other services!"
