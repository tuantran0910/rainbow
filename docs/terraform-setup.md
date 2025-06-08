# Terraform CI/CD Setup Guide

## 🔐 Authentication Setup

### 1. Create Google Cloud Service Account

```bash
# Set your project ID
export PROJECT_ID="rainbow-data-production"

# Create service account
gcloud iam service-accounts create terraform-ci \
  --display-name="Terraform CI/CD Service Account" \
  --project=$PROJECT_ID

# Grant necessary permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:terraform-ci@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/editor"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:terraform-ci@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

# Create and download service account key
gcloud iam service-accounts keys create terraform-ci-key.json \
  --iam-account=terraform-ci@$PROJECT_ID.iam.gserviceaccount.com
```

### 2. Store Credentials in GitHub Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Add the following secrets:

| Secret Name | Value | Description |
|------------|-------|-------------|
| `GCP_SA_KEY` | Content of `terraform-ci-key.json` | Service account credentials |

⚠️ **Important**: Delete the local `terraform-ci-key.json` file after uploading to GitHub secrets.

### 3. Setup GitHub Environments (Optional but Recommended)

1. Go to **Settings** → **Environments**
2. Create two environments:
   - `production` - for apply operations
   - `production-destroy` - for destroy operations
3. Add protection rules:
   - Required reviewers
   - Wait timer
   - Deployment branches (main only)

## 🚀 CI/CD Workflow Features

### Automatic Triggers
- **Pull Requests**: Runs `terraform plan` and posts results as PR comment
- **Push to main**: Runs `terraform plan` then `terraform apply` if changes detected
- **Manual trigger**: Can run `terraform destroy` (requires manual approval)

### Security Features
- ✅ **Plan before apply**: Always shows what will change
- ✅ **Environment protection**: Production deployments require approval
- ✅ **Credential isolation**: Service account keys stored as secrets
- ✅ **State locking**: GCS backend prevents concurrent modifications
- ✅ **Path filtering**: Only triggers on Terraform file changes

### Workflow Steps
1. **Format Check**: Ensures code is properly formatted
2. **Initialize**: Sets up Terraform and authenticates to GCP
3. **Validate**: Checks syntax and configuration
4. **Plan**: Shows what changes will be made
5. **Apply**: Applies changes (only on main branch)

## 🛠️ Local Development

### Prerequisites
```bash
# Install Terraform
brew install terraform

# Install Google Cloud CLI
brew install google-cloud-sdk

# Authenticate locally
gcloud auth login
gcloud config set project rainbow-data-production
```

### Local Commands
```bash
# Navigate to terraform directory
cd deployments/cloud/terraform

# Initialize Terraform
terraform init

# Plan changes
terraform plan

# Apply changes
terraform apply

# Format code
terraform fmt -recursive

# Validate configuration
terraform validate
```

## 🔒 Security Best Practices

### Service Account Permissions
- Use **least privilege principle**
- Regularly rotate service account keys
- Monitor service account usage

### State File Security
- State stored in GCS with encryption
- Access controlled via IAM
- Version-controlled in GCS

### Secrets Management
- Never commit credentials to Git
- Use GitHub secrets for CI/CD
- Rotate secrets regularly

## 🚨 Troubleshooting

### Common Issues

#### Authentication Errors
```bash
# Verify service account permissions
gcloud iam service-accounts get-iam-policy terraform-ci@$PROJECT_ID.iam.gserviceaccount.com
```

#### State Lock Issues
```bash
# Force unlock if needed (use carefully)
terraform force-unlock LOCK_ID
```

#### Quota Issues
- Check GCP quotas in Console
- Request quota increases if needed
- Optimize resource sizes

## 📋 Manual Deployment Checklist

Before pushing to main:
- [ ] Run `terraform fmt`
- [ ] Run `terraform validate`
- [ ] Run `terraform plan` locally
- [ ] Review plan output carefully
- [ ] Ensure no sensitive data in plan
- [ ] Test in development environment first

## 🔄 Rollback Strategy

If deployment fails:
1. Check GitHub Actions logs
2. Fix issues in new PR
3. If urgent, manually run:
   ```bash
   terraform apply -target=resource.name
   ```
4. For complete rollback, revert commit and re-run pipeline
