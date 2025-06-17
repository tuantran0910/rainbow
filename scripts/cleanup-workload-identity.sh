#!/bin/bash

# Cleanup Workload Identity Federation resources
# Run this script to clean up existing misconfigured resources

set -e

# Configuration
PROJECT_ID="rainbow-data-production"
POOL_ID="github-actions-pool-v2"
PROVIDER_ID="github-actions-provider-v2"
SERVICE_ACCOUNT_NAME="github-actions-deployer"

echo "🧹 Cleaning up Workload Identity Federation resources"
echo "Project ID: ${PROJECT_ID}"

# Remove Workload Identity Provider (if exists)
echo "🔗 Removing Workload Identity Provider..."
if gcloud iam workload-identity-pools providers describe "${PROVIDER_ID}" \
  --location="global" \
  --workload-identity-pool="${POOL_ID}" \
  --project="${PROJECT_ID}" \
  --quiet >/dev/null 2>&1; then
  
  gcloud iam workload-identity-pools providers delete "${PROVIDER_ID}" \
    --location="global" \
    --workload-identity-pool="${POOL_ID}" \
    --project="${PROJECT_ID}" \
    --quiet
  echo "✅ Workload Identity Provider removed"
else
  echo "ℹ️  Workload Identity Provider doesn't exist"
fi

# Handle Workload Identity Pool (including DELETED state)
echo "🌊 Removing Workload Identity Pool..."
POOL_STATE=$(gcloud iam workload-identity-pools describe "${POOL_ID}" \
  --location="global" \
  --project="${PROJECT_ID}" \
  --format="value(state)" \
  --quiet 2>/dev/null || echo "NOT_FOUND")

if [ "$POOL_STATE" = "DELETED" ]; then
  echo "⚠️  Pool is in DELETED state, undeleting first..."
  gcloud iam workload-identity-pools undelete "${POOL_ID}" \
    --location="global" \
    --project="${PROJECT_ID}" \
    --quiet
  echo "✅ Pool undeleted, now deleting permanently..."
  gcloud iam workload-identity-pools delete "${POOL_ID}" \
    --location="global" \
    --project="${PROJECT_ID}" \
    --quiet
  echo "✅ Workload Identity Pool removed"
elif [ "$POOL_STATE" = "ACTIVE" ]; then
  gcloud iam workload-identity-pools delete "${POOL_ID}" \
    --location="global" \
    --project="${PROJECT_ID}" \
    --quiet
  echo "✅ Workload Identity Pool removed"
else
  echo "ℹ️  Workload Identity Pool doesn't exist"
fi

# Remove service account (if exists)
echo "👤 Removing service account..."
if gcloud iam service-accounts describe "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --quiet >/dev/null 2>&1; then
  
  gcloud iam service-accounts delete "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --project="${PROJECT_ID}" \
    --quiet
  echo "✅ Service account removed"
else
  echo "ℹ️  Service account doesn't exist"
fi

echo ""
echo "✅ Cleanup completed!"
