#!/bin/bash
set -e

# Wait for MinIO to be ready
until (/usr/bin/mc config host add minio ${MINIO_ENDPOINT} ${MINIO_ROOT_USER} ${MINIO_ROOT_PASSWORD}); do
    echo 'Waiting for MinIO to be ready...'
    sleep 1
done

# Create bucket for Dagster compute logs (ignore if exists)
/usr/bin/mc mb minio/${DAGSTER_COMPUTE_LOGS_BUCKET} || true

# Set bucket to public access
/usr/bin/mc anonymous set public minio/${DAGSTER_COMPUTE_LOGS_BUCKET}
