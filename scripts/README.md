# Infrastructure Setup Scripts

This directory contains scripts to manually set up the infrastructure components for the Rainbow Data Production project.

## Overview

These scripts allow you to provision Google Cloud infrastructure manually instead of using Terraform resources. This is useful when you want to manage certain infrastructure components outside of Terraform.

## Prerequisites

Before running these scripts, ensure you have:

1. **Google Cloud CLI installed and configured**
   ```bash
   gcloud auth login
   gcloud config set project rainbow-data-production
   ```

2. **Required permissions** in the Google Cloud project:
   - Compute Admin
   - SQL Admin
   - Secret Manager Admin
   - Service Networking Admin

3. **OpenSSL** (for generating JWT secrets)

## Scripts

### 1. `setup-vpc.sh` - VPC Network Setup

Creates the VPC network and related networking components.

**What it creates:**
- Custom VPC network (`rainbow-data-production-vpc`)
- Subnet with secondary ranges for GKE (`rainbow-data-production-subnet`)
- Private IP address range for Cloud SQL
- Service networking connection for Cloud SQL
- Cloud Router and NAT Gateway
- Firewall rules (internal, SSH, HTTP, HTTPS)

**Usage:**
```bash
./scripts/setup-vpc.sh
```

**Resources created:**
- VPC Network: `rainbow-data-production-vpc`
- Subnet: `rainbow-data-production-subnet` (10.0.0.0/16)
- Secondary ranges: pods (10.1.0.0/16), services (10.2.0.0/16)
- Private IP: `rainbow-data-production-private-ip`
- Router: `rainbow-data-production-router`
- NAT: `rainbow-data-production-nat`

### 2. `setup-cloudsql.sh` - Cloud SQL Setup

Creates the Cloud SQL instance, database, users, and secrets.

**What it creates:**
- Cloud SQL PostgreSQL 15 instance (Enterprise edition)
- Database and users
- Secrets in Secret Manager

**Usage:**
```bash
./scripts/setup-cloudsql.sh
```

**Important:** Run the VPC setup script first, as Cloud SQL requires the VPC network.

**Resources created:**
- Cloud SQL Instance: `rainbow-data-production-db`
- Database: `rainbow`
- Users: `rainbow` (API), `datastream` (for data streaming)
- Secrets: `api-password`, `datastream-password`, `jwt-secret`

**Configuration:**
- Region: `us-west1`
- Tier: `db-custom-2-3840` (2 vCPU, 3.75GB RAM)
- Private IP only (no public access)
- SSL required
- IAM authentication enabled
- Logical decoding enabled for datastream

### 3. `setup-workload-identity.sh` - GitHub Actions Authentication

Sets up Workload Identity Federation for GitHub Actions deployment.

**Usage:**
```bash
./scripts/setup-workload-identity.sh
```

### 4. `cleanup-workload-identity.sh` - Cleanup Workload Identity

Removes the Workload Identity Federation setup.

**Usage:**
```bash
./scripts/cleanup-workload-identity.sh
```

## Execution Order

For a fresh setup, run the scripts in this order:

1. **VPC Network** (required first):
   ```bash
   ./scripts/setup-vpc.sh
   ```

2. **Cloud SQL** (requires VPC):
   ```bash
   ./scripts/setup-cloudsql.sh
   ```

3. **Workload Identity** (for GitHub Actions):
   ```bash
   ./scripts/setup-workload-identity.sh
   ```

## Terraform Integration

After running these scripts, the Terraform configuration will reference these manually created resources using data sources:

```hcl
# VPC Network
data "google_compute_network" "main" {
  name = "rainbow-data-production-vpc"
}

# Cloud SQL Instance
data "google_sql_database_instance" "main" {
  name = "rainbow-data-production-db"
}

# Database
data "google_sql_database" "rainbow" {
  name     = "rainbow"
  instance = data.google_sql_database_instance.main.name
}
```

## Security Features

### VPC Network
- Custom subnets with defined IP ranges
- Private Google Access enabled
- Cloud NAT for outbound internet access
- Restrictive firewall rules

### Cloud SQL
- Private IP only (no public access)
- SSL/TLS encryption required
- IAM authentication enabled
- Regular backups enabled
- Deletion protection enabled

### Secrets Management
- All passwords stored in Google Secret Manager
- JWT secret auto-generated with secure random data
- IAM-based access control

## Monitoring and Logs

The scripts include:
- Error handling with `set -e`
- Progress indicators with emojis
- Detailed output of created resources
- Graceful handling of existing resources

## Troubleshooting

### Common Issues

1. **API not enabled**: The scripts automatically enable required APIs and wait for them to be ready.

2. **Permissions**: Ensure your account has the required IAM roles listed in prerequisites.

3. **Existing resources**: Scripts handle existing resources gracefully with informational messages.

4. **Network connectivity**: Cloud SQL requires the VPC network to be created first.

### Verification

After running the scripts, verify the setup:

```bash
# Check VPC network
gcloud compute networks describe rainbow-data-production-vpc

# Check Cloud SQL instance
gcloud sql instances describe rainbow-data-production-db

# Check secrets
gcloud secrets list --filter="name~rainbow"
```

## Cleanup

To remove manually created resources:

1. Delete Cloud SQL instance:
   ```bash
   gcloud sql instances delete rainbow-data-production-db
   ```

2. Delete VPC network (after removing all dependent resources):
   ```bash
   gcloud compute networks delete rainbow-data-production-vpc
   ```

3. Delete secrets:
   ```bash
   gcloud secrets delete api-password
   gcloud secrets delete datastream-password
   gcloud secrets delete jwt-secret
   ```

**⚠️ Warning**: Deleting these resources will affect your production environment. Always backup data before cleanup.
