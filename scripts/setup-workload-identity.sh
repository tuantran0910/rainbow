#!/bin/bash

# Setup Workload Identity Federation for GitHub Actions
# This script needs to be run once to set up authentication between GitHub Actions and Google Cloud

set -e

# Configuration
PROJECT_ID="rainbow-data-production"
POOL_ID="github-actions-pool-v2"
PROVIDER_ID="github-actions-provider-v2"
SERVICE_ACCOUNT_NAME="github-actions-deployer"
GITHUB_REPO="tuantran0910/rainbow"

echo "🚀 Setting up Workload Identity Federation for GitHub Actions"
echo "Project ID: ${PROJECT_ID}"
echo "GitHub Repository: ${GITHUB_REPO}"

# Enable required APIs
echo "📡 Enabling required APIs..."
gcloud services enable iamcredentials.googleapis.com --project="${PROJECT_ID}"
gcloud services enable sts.googleapis.com --project="${PROJECT_ID}"

# Create service account for GitHub Actions
echo "👤 Creating service account for GitHub Actions..."
gcloud iam service-accounts create "${SERVICE_ACCOUNT_NAME}" \
  --display-name="GitHub Actions Deployer" \
  --description="Service account for GitHub Actions CI/CD" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Service account already exists"

# Grant necessary permissions to the service account
echo "🔐 Granting permissions to service account..."
ROLES=(
  "roles/run.developer"
  "roles/iam.serviceAccountUser"
  "roles/secretmanager.secretAccessor"
  "roles/cloudsql.client"
)

for ROLE in "${ROLES[@]}"; do
  echo "  Adding role: ${ROLE}"
  gcloud projects add-iam-policy-binding "${PROJECT_ID}" \
    --member="serviceAccount:${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
    --role="${ROLE}" >/dev/null
done

# Create Workload Identity Pool
echo "🌊 Creating Workload Identity Pool..."
gcloud iam workload-identity-pools create "${POOL_ID}" \
  --location="global" \
  --display-name="GitHub Actions Pool" \
  --description="Pool for GitHub Actions authentication" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Workload Identity Pool creation failed - may already exist"

# Create Workload Identity Provider
echo "🔗 Creating Workload Identity Provider..."
gcloud iam workload-identity-pools providers create-oidc "${PROVIDER_ID}" \
  --location="global" \
  --workload-identity-pool="${POOL_ID}" \
  --display-name="GitHub Actions Provider" \
  --description="OIDC provider for GitHub Actions" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository,attribute.repository_owner=assertion.repository_owner" \
  --attribute-condition="assertion.repository_owner == 'tuantran0910'" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --project="${PROJECT_ID}" 2>/dev/null || echo "ℹ️  Workload Identity Provider creation failed - may already exist"

# Wait a moment for resources to be fully created
echo "⏳ Waiting for resources to be fully provisioned..."
sleep 5

# Allow GitHub Actions to impersonate the service account
echo "🎭 Setting up service account impersonation..."
gcloud iam service-accounts add-iam-policy-binding \
  "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/$(gcloud projects describe ${PROJECT_ID} --format='value(projectNumber)')/locations/global/workloadIdentityPools/${POOL_ID}/attribute.repository/${GITHUB_REPO}" \
  --project="${PROJECT_ID}"

# Get the Workload Identity Provider resource name
echo "🔍 Getting Workload Identity Provider details..."
PROJECT_NUMBER=$(gcloud projects describe "${PROJECT_ID}" --format="value(projectNumber)")
WIF_PROVIDER="projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL_ID}/providers/${PROVIDER_ID}"
WIF_SERVICE_ACCOUNT="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

echo ""
echo "✅ Workload Identity Federation setup completed!"
echo ""
echo "🔧 Add these secrets to your GitHub repository:"
echo "   Go to: https://github.com/${GITHUB_REPO}/settings/secrets/actions"
echo ""
echo "   WIF_PROVIDER: ${WIF_PROVIDER}"
echo "   WIF_SERVICE_ACCOUNT: ${WIF_SERVICE_ACCOUNT}"
echo ""
echo "📋 You can also run these commands to add them via GitHub CLI:"
echo "   gh secret set WIF_PROVIDER --body \"${WIF_PROVIDER}\" --repo ${GITHUB_REPO}"
echo "   gh secret set WIF_SERVICE_ACCOUNT --body \"${WIF_SERVICE_ACCOUNT}\" --repo ${GITHUB_REPO}"
echo ""
echo "🚀 Your GitHub Actions workflow is now ready to deploy to Google Cloud!" 
