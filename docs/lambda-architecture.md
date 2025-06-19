# Lambda Architecture Implementation

This document describes the lambda architecture implementation for the Rainbow dbt project, which combines batch processing (T-1) with real-time streaming data (T) to provide unified, real-time views.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐
│   Raw Dataset   │    │Raw__Datastream  │
│   (Batch T-1)   │    │   Dataset (T)   │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│  Staging/Batch  │    │ Staging/Stream  │
│     Models      │    │     Models      │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          ▼                      ▼
┌─────────────────┐    ┌─────────────────┐
│Intermediate/    │    │                 │
│Batch_Layer      │    │                 │
│   Models        │    │                 │
└─────────┬───────┘    │                 │
          │            │                 │
          └────────────┼─────────────────┘
                       ▼
               ┌─────────────────┐
               │  Unified Lambda │
               │      Views      │
               │ (Intermediate)  │
               └─────────────────┘
```

## Directory Structure

```
models/
├── sources/
│   └── sources.yml (includes both rainbow and rainbow_datastream)
├── staging/
│   ├── batch/
│   │   ├── stg_*.sql (existing batch staging models)
│   │   └── schema.yml
│   └── stream/
│       ├── stg_stream_*.sql (new streaming staging models)
│       └── schema.yml
├── intermediate/
│   ├── batch_layer/
│   │   ├── dim_*_batch.sql (batch-processed dimensions)
│   │   ├── fct_*_batch.sql (batch-processed facts)
│   │   └── schema.yml
│   ├── dim_*.sql (unified lambda views)
│   └── fct_*.sql (unified lambda views)
└── marts/ (future)
```

## Key Components

### 1. Macros

#### `lambda_view` Macro
- **Purpose**: Combines batch layer data (T-1) with streaming data (T)
- **Logic**: 
  - Takes data older than cutoff from batch layer
  - Takes recent data from streaming layer
  - Stream data takes precedence for overlapping time periods
- **Usage**:
```sql
{{ lambda_view(
    batch_relation=ref('dim_users_batch'),
    stream_stg_relation=ref('stg_stream_users'),
    unique_key='id',
    updated_at_field='updated_at_tz_hcm',
    cutoff_hours=24
) }}
```

#### `lambda_fact_sales` Macro
- **Purpose**: Specialized macro for sales fact table lambda processing
- **Features**: 
  - Preserves complex business logic for promotions and discounts
  - Handles multiple source tables (orders, order_items)
  - Maintains same calculation logic as batch processing
- **Usage**:
```sql
{{ lambda_fact_sales(
    batch_orders_relation=ref('stg_orders'),
    batch_order_items_relation=ref('stg_order_items'),
    stream_orders_relation=ref('stg_stream_orders'),
    stream_order_items_relation=ref('stg_stream_order_items'),
    cutoff_hours=24
) }}
```

### 2. Data Sources

#### Batch Sources (`rainbow` dataset)
- Traditional batch processing data
- Updated via scheduled ETL processes
- Contains historical data (T-1)

#### Streaming Sources (`rainbow_datastream` dataset)
- Real-time CDC data from DataStream
- Contains metadata fields:
  - `_metadata_timestamp`: DataStream ingestion timestamp
  - `_metadata_deleted`: Deletion indicator
- Provides near real-time data (T)

### 3. Model Types

#### Staging Models
- **Batch**: Simple views over raw batch data
- **Stream**: Views over streaming data with DataStream metadata handling

#### Intermediate Models
- **Batch Layer**: Traditional SCD Type 2 and fact processing for batch data
- **Lambda Views**: Unified views combining batch + stream using lambda macros

## Implementation Guide

### Step 1: Verify Source Connections
Ensure both `rainbow` and `rainbow_datastream` datasets are accessible:
```bash
dbt test --select source:rainbow
dbt test --select source:rainbow_datastream
```

### Step 2: Build Batch Layer
Build the batch processing models first:
```bash
dbt build --select models/staging/batch
dbt build --select models/intermediate/batch_layer
```

### Step 3: Build Stream Staging
Build the streaming staging models:
```bash
dbt build --select models/staging/stream
```

### Step 4: Build Lambda Views
Build the unified lambda views:
```bash
dbt build --select models/intermediate --exclude models/intermediate/batch_layer
```

### Step 5: Validation
Compare data between batch and lambda views to ensure consistency:
```sql
-- Example validation query
SELECT 
    'batch' as source,
    COUNT(*) as record_count,
    MAX(updated_at_tz_hcm) as latest_update
FROM {{ ref('dim_users_batch') }}
WHERE is_current = TRUE

UNION ALL

SELECT 
    'lambda' as source,
    COUNT(*) as record_count,
    MAX(updated_at_tz_hcm) as latest_update
FROM {{ ref('dim_users') }}
```

## Configuration

### dbt_project.yml
```yaml
models:
  rainbow_transformations:
    +schema: warehouse
    staging:
      batch:
        +schema: warehouse_staging_batch
      stream:
        +schema: warehouse_staging_stream
    intermediate:
      batch_layer:
        +schema: warehouse_batch
      +schema: warehouse  # Unified lambda views
```

## Best Practices

### 1. Cutoff Time Management
- Default cutoff is 24 hours
- Adjust based on batch processing frequency and data freshness requirements
- Consider business hours and processing windows

### 2. Data Quality
- Monitor DataStream metadata fields for data quality issues
- Implement alerts for high deletion rates or ingestion delays
- Validate transformation logic consistency between batch and stream

### 3. Performance Optimization
- Lambda views are materialized as views for real-time access
- Consider incremental materialization for large fact tables
- Monitor query performance and optimize as needed

### 4. Testing Strategy
- Test both batch and stream paths independently
- Validate business logic consistency
- Implement data quality tests for lambda views

## Troubleshooting

### Common Issues

1. **Reference Errors**: Ensure model names match between batch and lambda layers
2. **Schema Mismatches**: Verify DataStream and batch schemas are compatible
3. **Performance Issues**: Consider materialization strategies for large datasets
4. **Data Freshness**: Monitor cutoff times and adjust based on processing schedules

### Monitoring Queries

```sql
-- Check lambda view freshness
SELECT 
    'users' as table_name,
    COUNT(*) as total_records,
    COUNT(*) FILTER (WHERE updated_at_tz_hcm >= CURRENT_TIMESTAMP - INTERVAL '1 hour') as recent_records
FROM {{ ref('dim_users') }}

UNION ALL

SELECT 
    'sales' as table_name,
    COUNT(*) as total_records,
    COUNT(*) FILTER (WHERE created_at_tz_hcm >= CURRENT_TIMESTAMP - INTERVAL '1 hour') as recent_records
FROM {{ ref('fct_sales') }}
```

## Migration Path

To implement lambda architecture on existing models:

1. **Backup**: Ensure current models are backed up
2. **Reorganize**: Move existing models to batch_layer directory
3. **Add Stream Sources**: Create DataStream source definitions
4. **Create Stream Staging**: Build stream staging models
5. **Implement Lambda Views**: Create unified views using lambda macros
6. **Test**: Validate data consistency and performance
7. **Switch References**: Update downstream models to use lambda views

This implementation ensures that you maintain the same transformation logic while gaining real-time capabilities through the lambda architecture pattern. 