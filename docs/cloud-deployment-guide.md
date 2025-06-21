# Cloud Run Deployment Guide

This guide explains the changes made to support deploying the Rainbow API to Google Cloud Run with Cloud SQL.

## Overview

The API codebase has been modified to support both local development (Docker) and cloud deployment (Cloud Run + Cloud SQL) environments.

## Key Changes Made

### 1. Database Configuration (`src/api/config/config.go`)

**Problem**: The original code hardcoded `sslmode=disable` which is incompatible with Cloud SQL security requirements.

**Solution**: Added dynamic SSL mode configuration based on environment:

```go
// Determine SSL mode based on environment and USE_CLOUD_SQL flag
useCloudSQL := getEnv("USE_CLOUD_SQL", "false")
var sslMode string
if useCloudSQL == "true" || env == "production" {
    sslMode = "require"
} else {
    sslMode = "disable"
}
```

**Environment Variable Standardization**:
- Normalized to use `DB_*` environment variables for both local and cloud deployments
- Simplified configuration by removing dual naming convention
- Consistent variable names across all environments

### 2. Dockerfile Updates (`src/api/Dockerfile`)

**Problem**: Migrations directory was not included in the Docker image.

**Solution**: Added migrations directory copy:

```dockerfile
# Copy the migrations directory
COPY --from=builder /build/migrations /migrations
```

### 3. Cloud Run Configuration (`deployments/cloud/terraform/cloudrun.tf`)

**Added Environment Variables**:

- `ENV=production` - Sets production environment mode
- `USE_CLOUD_SQL=true` - Enables Cloud SQL specific configurations
- `MIGRATIONS_DIR=/migrations` - Points to migrations in the container
- `LOG_LEVEL=info` - Appropriate log level for production
- `ENABLE_CONSOLE=false` - Disables console logging in favor of structured logging
- `JWT_SECRET` - From Google Secret Manager
- `SERVER_PORT=5000` - Application port

**Secret Management**: Added JWT secret data source:

```hcl
data "google_secret_manager_secret_version" "jwt_secret" {
  secret = "jwt-secret"
}
```

## Environment Variables

### Local Development (Docker Compose)
```env
# Database
DB_HOST=postgres
DB_USER=rainbow
DB_PASSWORD=R&inb0w2024!Data
DB_NAME=rainbow
DB_PORT=5432

# Server
SERVER_PORT=5000
ENV=development
MIGRATIONS_DIR=migrations

# Security
JWT_SECRET=jwt_secret
```

### Cloud Run Production
```env
# Database (Cloud SQL)
DB_HOST=<private-ip-of-cloud-sql>
DB_USER=rainbow
DB_PASSWORD=<from-secret-manager>
DB_NAME=rainbow
DB_PORT=5432
USE_CLOUD_SQL=true

# Server
SERVER_PORT=5000
ENV=production
MIGRATIONS_DIR=/migrations

# Logging
LOG_LEVEL=info
ENABLE_CONSOLE=false

# Security
JWT_SECRET=<from-secret-manager>
```

## SSL/TLS Configuration

### Local Development
- Uses `sslmode=disable` for simplicity
- Connects to local PostgreSQL container

### Cloud SQL Production
- Uses `sslmode=require` for security
- Connects via VPC connector to private Cloud SQL instance
- SSL certificates handled automatically by Cloud SQL

## Prerequisites for Cloud Deployment

### 1. Google Secret Manager Secrets
Create the following secrets in Google Secret Manager:

```bash
# API Database Password
gcloud secrets create api-password --data-file=<password-file>

# JWT Secret
gcloud secrets create jwt-secret --data-file=<jwt-secret-file>
```

### 2. Cloud SQL Instance
The terraform configuration creates:
- Private Cloud SQL PostgreSQL 16 instance
- VPC peering for secure connectivity
- Database user and application database

### 3. Container Image
Build and push the API container image:

```bash
# Build the image
cd src/api
docker build -t gcr.io/<PROJECT_ID>/rainbow-api:latest .

# Push to Google Container Registry
docker push gcr.io/<PROJECT_ID>/rainbow-api:latest

# Update terraform to use the correct image
# Edit cloudrun.tf and change the image from "hello" to your image
```

## Database Migrations

Migrations are handled automatically on application startup:
- Local: Reads from `./migrations` directory
- Cloud Run: Reads from `/migrations` directory (copied during build)

The `goose` migration tool is used to apply migrations automatically when the application starts.

## Monitoring and Logging

### Local Development
- Console logging enabled with debug level
- Structured JSON logs to stdout

### Cloud Production
- Console logging disabled
- Info level logging
- Logs sent to Google Cloud Logging automatically
- Health check endpoint: `/health`

## Security Considerations

1. **Database Access**: Cloud SQL uses private IP with VPC connector
2. **Secret Management**: Sensitive data stored in Google Secret Manager
3. **SSL/TLS**: Required for all Cloud SQL connections
4. **Service Account**: Dedicated service account with minimal required permissions

## Testing the Deployment

### Health Check
```bash
curl https://<cloud-run-url>/health
```

### API Endpoints
```bash
# Check API is running
curl https://<cloud-run-url>/

# Test authentication
curl -X POST https://<cloud-run-url>/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Admin123!"}'
```

## Troubleshooting

### Common Issues

1. **Migration Failures**: Check that `/migrations` directory exists in container
2. **Database Connection**: Verify VPC connector and Cloud SQL private IP
3. **Secret Access**: Ensure service account has Secret Manager access
4. **SSL Errors**: Verify `USE_CLOUD_SQL=true` is set for Cloud Run

### Debugging Commands

```bash
# Check Cloud Run logs
gcloud run services logs read rainbow-api --region=<region>

# Test database connectivity from Cloud Run
gcloud run jobs create db-test \
  --image=postgres:16-alpine \
  --env-vars=PGHOST=<cloud-sql-private-ip>,PGUSER=rainbow,PGDATABASE=rainbow
```

## Migration from Local to Cloud

1. **Data Migration**: Use `pg_dump` and `pg_restore` to migrate data
2. **DNS Setup**: Configure custom domain in Cloud Run if needed
3. **Load Testing**: Test with expected traffic patterns
4. **Monitoring**: Set up alerts for errors and performance metrics

This configuration ensures a smooth transition from local development to cloud production while maintaining security and performance best practices.
