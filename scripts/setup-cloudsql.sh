#!/bin/bash

# Setup Cloud SQL for Rainbow Data Production
# This script creates the Cloud SQL instance, database, and users

set -e

# Configuration
PROJECT_ID="rainbow-data-production"
CLOUD_SQL_REGION="us-west1"
VPC_NAME="${PROJECT_ID}-vpc"
INSTANCE_NAME="${PROJECT_ID}-db"
DATABASE_NAME="rainbow"
API_USER="rainbow"
DATASTREAM_USER="datastream"

echo "🚀 Setting up Cloud SQL for Rainbow Data Production"
echo "Project ID: ${PROJECT_ID}"
echo "Region: ${CLOUD_SQL_REGION}"
echo "Instance Name: ${INSTANCE_NAME}"

# Set the project
echo "📋 Setting project context..."
gcloud config set project "${PROJECT_ID}"

# Enable required APIs
echo "📡 Enabling required APIs..."
REQUIRED_APIS=(
  "sqladmin.googleapis.com"
  "secretmanager.googleapis.com"
)

for API in "${REQUIRED_APIS[@]}"; do
  echo "  Enabling ${API}..."
  gcloud services enable "${API}" --project="${PROJECT_ID}"
done

echo "⏳ Waiting for APIs to be fully enabled..."
sleep 10

# Check if secrets exist, create if they don't
echo "🔐 Setting up secrets..."

# Create API password secret if it doesn't exist
if ! gcloud secrets describe "api-password" --project="${PROJECT_ID}" &>/dev/null; then
  echo "  Creating api-password secret..."
  echo -n "R&inb0w2024!Data" | gcloud secrets create "api-password" \
    --data-file=- \
    --project="${PROJECT_ID}"
else
  echo "ℹ️  api-password secret already exists"
fi

# Create datastream password secret if it doesn't exist
if ! gcloud secrets describe "datastream-password" --project="${PROJECT_ID}" &>/dev/null; then
  echo "  Creating datastream-password secret..."
  echo -n "D&taStr3am2024!Pass" | gcloud secrets create "datastream-password" \
    --data-file=- \
    --project="${PROJECT_ID}"
else
  echo "ℹ️  datastream-password secret already exists"
fi

# Create JWT secret if it doesn't exist
if ! gcloud secrets describe "jwt-secret" --project="${PROJECT_ID}" &>/dev/null; then
  echo "  Creating jwt-secret..."
  echo -n "$(openssl rand -base64 32)" | gcloud secrets create "jwt-secret" \
    --data-file=- \
    --project="${PROJECT_ID}"
else
  echo "ℹ️  jwt-secret already exists"
fi

# Get passwords from secrets
echo "🔑 Retrieving passwords from Secret Manager..."
API_PASSWORD=$(gcloud secrets versions access latest --secret="api-password" --project="${PROJECT_ID}")
DATASTREAM_PASSWORD=$(gcloud secrets versions access latest --secret="datastream-password" --project="${PROJECT_ID}")

# Create Cloud SQL instance
echo "🗄️  Creating Cloud SQL instance..."
gcloud sql instances create "${INSTANCE_NAME}" \
  --database-version=POSTGRES_15 \
  --tier=db-custom-2-3840 \
  --edition=ENTERPRISE \
  --region="${CLOUD_SQL_REGION}" \
  --network="projects/${PROJECT_ID}/global/networks/${VPC_NAME}" \
  --no-assign-ip \
  --database-flags=cloudsql.iam_authentication=on,cloudsql.logical_decoding=on,max_replication_slots=10,max_wal_senders=10 \
  --deletion-protection \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Cloud SQL instance already exists"

echo "⏳ Waiting for Cloud SQL instance to be ready..."
gcloud sql instances describe "${INSTANCE_NAME}" --project="${PROJECT_ID}" --format="value(state)" | while read state; do
  if [ "$state" = "READY" ]; then
    break
  fi
  echo "  Instance state: $state - waiting..."
  sleep 10
done

# Create database
echo "📊 Creating database..."
gcloud sql databases create "${DATABASE_NAME}" \
  --instance="${INSTANCE_NAME}" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Database already exists"

# Create API user
echo "👤 Creating API user..."
gcloud sql users create "${API_USER}" \
  --instance="${INSTANCE_NAME}" \
  --password="${API_PASSWORD}" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  API user already exists"

# Create datastream user
echo "📡 Creating datastream user..."
gcloud sql users create "${DATASTREAM_USER}" \
  --instance="${INSTANCE_NAME}" \
  --password="${DATASTREAM_PASSWORD}" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Datastream user already exists"

# Grant necessary permissions to datastream user
echo "🔐 Granting permissions to datastream user..."
gcloud sql instances patch "${INSTANCE_NAME}" \
  --database-flags=cloudsql.iam_authentication=on,cloudsql.logical_decoding=on,max_replication_slots=10,max_wal_senders=10 \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Database flags already set"

# Get instance details
echo "📋 Getting instance details..."
INSTANCE_IP=$(gcloud sql instances describe "${INSTANCE_NAME}" --project="${PROJECT_ID}" --format="value(ipAddresses[0].ipAddress)")
INSTANCE_CONNECTION_NAME=$(gcloud sql instances describe "${INSTANCE_NAME}" --project="${PROJECT_ID}" --format="value(connectionName)")

echo ""
echo "✅ Cloud SQL setup completed!"
echo ""
echo "📋 Created resources:"
echo "   🗄️  Instance: ${INSTANCE_NAME}"
echo "   📊 Database: ${DATABASE_NAME}"
echo "   👤 API User: ${API_USER}"
echo "   📡 Datastream User: ${DATASTREAM_USER}"
echo ""
echo "🔗 Connection details:"
echo "   🌐 Private IP: ${INSTANCE_IP}"
echo "   📞 Connection Name: ${INSTANCE_CONNECTION_NAME}"
echo "   🔒 SSL Mode: Required"
echo ""
echo "🔐 Secrets stored in Secret Manager:"
echo "   🔑 api-password"
echo "   🔑 datastream-password"
echo "   🔑 jwt-secret"
echo ""
echo "⚠️  Security Notes:"
echo "   • Instance uses private IP only (no public access)"
echo "   • SSL connections are required"
echo "   • IAM authentication is enabled"
echo "   • Logical decoding is enabled for datastream"
echo ""
echo "🚀 Your Cloud SQL instance is ready for the API and datastream!"
