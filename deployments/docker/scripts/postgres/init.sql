---------------------------------------------------------------------
--                      Dagster                                    --
---------------------------------------------------------------------
\echo 'Setting up Dagster database and user...'

-- Initialize the necessary resources for the Dagster.
CREATE DATABASE dagster;

-- Drop the user if it already exists (global operation)
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'dagster'
    ) THEN
        DROP ROLE dagster;
    END IF;
END $$;

-- Create the user with a password (global operation)
CREATE ROLE dagster WITH LOGIN PASSWORD 'Dagster2024!Data';

-- Grant necessary connection privileges
GRANT CONNECT ON DATABASE dagster TO dagster;
GRANT TEMPORARY ON DATABASE dagster TO dagster;

-- Switch to the dagster database to grant schema permissions
\echo 'Connecting to dagster database...'
\c dagster
\echo 'Granting schema permissions in dagster database...'
GRANT USAGE, CREATE ON SCHEMA public TO dagster;
\echo 'Dagster setup complete.'

---------------------------------------------------------------------
--                      API (Rainbow)                              --
---------------------------------------------------------------------
\echo 'Setting up API (Rainbow) database and user...'

-- Initialize the necessary resources for the API
CREATE DATABASE rainbow;

-- Drop the user if it already exists (global operation)
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'rainbow'
    ) THEN
        DROP ROLE rainbow;
    END IF;
END $$;

-- Create the user with a password (global operation)
CREATE ROLE rainbow WITH LOGIN PASSWORD 'R&inb0w2024!Data';

-- Grant necessary connection privileges
GRANT CONNECT ON DATABASE rainbow TO rainbow;
GRANT TEMPORARY ON DATABASE rainbow TO rainbow;

-- Switch to the rainbow database to grant schema permissions and create extensions
\echo 'Connecting to rainbow database...'
\c rainbow
\echo 'Granting schema permissions and creating extensions in rainbow database...'
GRANT USAGE, CREATE ON SCHEMA public TO rainbow;

-- Create the extension for UUID generation (needs to be done within the rainbow database)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
\echo 'API (Rainbow) setup complete.'

---------------------------------------------------------------------
--                      CDC                                        --
---------------------------------------------------------------------
\echo 'Setting up CDC user and permissions...'

-- NOTE: We are currently connected to the 'rainbow' database from the previous section.

-- Drop the user if it already exists (global operation)
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'cdc'
    ) THEN
        DROP ROLE cdc;
    END IF;
END $$;

-- Create the user with a password (global operation)
-- REPLICATION privilege allows logical replication connections
CREATE ROLE cdc WITH REPLICATION LOGIN PASSWORD 'Cdc2024!Data';

-- Grant necessary connection privileges
GRANT CONNECT ON DATABASE rainbow TO cdc;
GRANT TEMPORARY ON DATABASE rainbow TO cdc; -- Usually not needed for CDC, but keeping for consistency

-- Grant schema and table privileges
\echo 'Granting CDC permissions within rainbow database...'
GRANT USAGE ON SCHEMA public TO cdc;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO cdc;

-- Set the default privileges for future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT ON TABLES TO cdc;

-- Create the publication for logical replication
\echo 'Creating CDC publication in rainbow database...'
CREATE PUBLICATION cdc_publication FOR ALL TABLES;
\echo 'CDC setup complete.'

---------------------------------------------------------------------
--                      Metabase                                   --
---------------------------------------------------------------------
\echo 'Setting up Metabase database and user...'

-- Initialize the necessary resources for Metabase
CREATE DATABASE metabase;

-- Drop the user if it already exists (global operation)
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'metabase'
    ) THEN
        DROP ROLE metabase;
    END IF;
END $$;

-- Create the user with a password (global operation)
CREATE ROLE metabase WITH LOGIN PASSWORD 'Metabase2024!Data';

-- Grant necessary connection privileges
GRANT CONNECT ON DATABASE metabase TO metabase;
GRANT TEMPORARY ON DATABASE metabase TO metabase;
-- GRANT CREATE ON DATABASE metabase TO metabase; -- Grants ability to create schemas within metabase

-- Switch to the metabase database to grant schema permissions
\echo 'Connecting to metabase database...'
\c metabase
\echo 'Granting schema permissions in metabase database...'
GRANT USAGE, CREATE ON SCHEMA public TO metabase;
GRANT CREATE ON DATABASE metabase TO metabase; -- Redundant if granted before \c, but confirms intent within the DB context.
\echo 'Metabase setup complete.'

---------------------------------------------------------------------
--                      EL                                         --
---------------------------------------------------------------------
\echo 'Setting up EL user and permissions...'

-- Drop the user if it already exists (global operation)
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'el'
    ) THEN
        DROP ROLE el;
    END IF;
END $$;

-- Create the user with a password (global operation, connected to metabase now)
CREATE ROLE el WITH LOGIN PASSWORD 'El2024!Data';

-- Grant necessary connection privilege to rainbow database
GRANT CONNECT ON DATABASE rainbow TO el;

-- Switch to the 'rainbow' database
\echo 'Connecting to rainbow database for EL permissions...'
\c rainbow

-- Grant schema and table privileges
\echo 'Granting EL permissions within rainbow database...'
GRANT USAGE ON SCHEMA public TO el;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO el;

-- Ensure future tables are also accessible
ALTER DEFAULT PRIVILEGES FOR ROLE rainbow IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO el;
\echo 'EL setup complete.'

\echo 'Database initialization script finished.'
