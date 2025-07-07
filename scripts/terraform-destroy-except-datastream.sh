#!/bin/bash

echo "🔥 Destroying all resources except Datastream configurations..."
echo ""
echo "📋 Resources that will be PRESERVED:"
echo "  - google_datastream_connection_profile.bq_destination"
echo "  - google_datastream_connection_profile.cloudsql_source[\"rainbow\"]"
echo "  - google_datastream_stream.stream[\"rainbow\"]"
echo "  - google_bigquery_dataset.datastream_dataset[\"rainbow\"]"
echo ""
echo "🗑️  Resources that will be DESTROYED:"
echo "  - GKE cluster and node pools"
echo "  - Networking (subnets, NAT, routers, etc.)"
echo "  - DNS records"
echo "  - Service accounts and IAM bindings (non-datastream)"
echo "  - Storage IAM bindings"
echo ""

read -p "Continue? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 1
fi

# Set working directory to the deployments/cloud/terraform directory
cd deployments/cloud/terraform

# Destroy specific non-datastream resources using -target
terraform destroy \
  -target="google_compute_global_address.dagster_ingress_ip" \
  -target="google_compute_global_address.private_ip_address" \
  -target="google_compute_global_address.superset_ingress_ip" \
  -target="google_compute_router.router" \
  -target="google_compute_router_nat.nat" \
  -target="google_compute_ssl_policy.rainbow_ssl_policy" \
  -target="google_compute_subnetwork.main" \
  -target="google_container_cluster.primary" \
  -target="google_container_node_pool.primary_nodes" \
  -target="google_dns_record_set.dagster_dns" \
  -target="google_dns_record_set.superset_dns" \
  -target="google_gke_hub_feature.config_management" \
  -target="google_gke_hub_feature_membership.configmanagement_feature_member" \
  -target="google_gke_hub_membership.membership" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/editor\"]" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/iam.serviceAccountAdmin\"]" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/iam.serviceAccountTokenCreator\"]" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/iam.serviceAccountUser\"]" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/iam.workloadIdentityPoolAdmin\"]" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/resourcemanager.projectIamAdmin\"]" \
  -target="google_project_iam_member.config_connector_bindings[\"roles/serviceusage.serviceUsageConsumer\"]" \
  -target="google_project_iam_member.gke_node_sa_bindings[\"roles/logging.logWriter\"]" \
  -target="google_project_iam_member.gke_node_sa_bindings[\"roles/monitoring.metricWriter\"]" \
  -target="google_project_iam_member.gke_node_sa_bindings[\"roles/monitoring.viewer\"]" \
  -target="google_project_iam_member.gke_node_sa_bindings[\"roles/storage.objectViewer\"]" \
  -target="google_service_account.config_connector" \
  -target="google_service_account.gke_node_sa" \
  -target="google_service_account_iam_member.config_connector_workload_identity" \
  -target="google_service_networking_connection.private_vpc_connection" \
  -target="google_storage_bucket_iam_member.dbt_docs_public_read"

# Return to the root of the repository
cd ../../..
