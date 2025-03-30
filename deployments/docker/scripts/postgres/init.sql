---------------------------------------------------------------------
--                      Dagster                                    --
---------------------------------------------------------------------

-- Initialize the necessary resources for the Dagster.
CREATE DATABASE dagster;

-- Drop the user if it already exists
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'dagster'
    ) THEN
        DROP ROLE dagster;
    END IF;
END $$;

-- Create the user with a password
CREATE ROLE dagster WITH LOGIN PASSWORD 'Dagster2024!Data';

-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE dagster TO dagster;
GRANT TEMPORARY ON DATABASE dagster TO dagster;

-- Switch to the dagster database and grant schema permissions
\c dagster
GRANT USAGE, CREATE ON SCHEMA public TO dagster;

---------------------------------------------------------------------
--                      API                                        --
---------------------------------------------------------------------

-- Initialize the necessary resources for the API
CREATE DATABASE rainbow;

-- Create the user with a password
CREATE ROLE rainbow WITH LOGIN PASSWORD 'R&inb0w2024!Data';

-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE rainbow TO rainbow;
GRANT TEMPORARY ON DATABASE rainbow TO rainbow;

-- Switch to the rainbow database and grant schema permissions
\c rainbow
GRANT USAGE, CREATE ON SCHEMA public TO rainbow;

-- Create the extension for UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

---------------------------------------------------------------------
--                      CDC                                        --
---------------------------------------------------------------------

-- Create the user with a password
CREATE ROLE cdc WITH REPLICATION LOGIN PASSWORD 'Cdc2024!Data';

-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE rainbow TO cdc;
GRANT TEMPORARY ON DATABASE rainbow TO cdc;

-- Create the extension for logical replication
CREATE PUBLICATION cdc_publication FOR ALL TABLES;

-- Grant schema and table privileges
\c rainbow
GRANT USAGE ON SCHEMA public TO cdc;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO cdc;

-- Set the default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT ON TABLES TO cdc;

---------------------------------------------------------------------
--                      Metabase                                   --
---------------------------------------------------------------------

-- Initialize the necessary resources for the Metabase
CREATE DATABASE metabase;

-- Create the user with a password
CREATE ROLE metabase WITH LOGIN PASSWORD 'Metabase2024!Data';

-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE metabase TO metabase;
GRANT TEMPORARY ON DATABASE metabase TO metabase;

-- Switch to the metabase database and grant schema permissions
\c metabase
GRANT CREATE ON DATABASE metabase TO metabase;
GRANT USAGE, CREATE ON SCHEMA public TO metabase;

---------------------------------------------------------------------
--                      EL                                         --
---------------------------------------------------------------------

-- Create the user with a password
CREATE ROLE el WITH LOGIN PASSWORD 'El2024!Data';

-- TODO: switch to database rainbow
-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE dagster TO el;

-- Grant schema and table privileges
\c dagster
GRANT USAGE ON SCHEMA public TO el;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO jerry;

-- Ensure future tables are also accessible
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO el;
