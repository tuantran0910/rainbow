---------------------------------------------------------------------
--                      Apache Airflow                             --
---------------------------------------------------------------------

-- Initialize the necessary resources for the Apache Airflow.
CREATE DATABASE airflow;

-- Drop the user if it already exists
DO $$
BEGIN
    IF EXISTS (
        SELECT FROM pg_catalog.pg_roles WHERE rolname = 'airflow'
    ) THEN
        DROP ROLE airflow;
    END IF;
END $$;

-- Create the user with a password
CREATE ROLE airflow WITH LOGIN PASSWORD 'airflow';

-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE airflow TO airflow;
GRANT TEMPORARY ON DATABASE airflow TO airflow;

-- Switch to the airflow database and grant schema permissions
\c airflow
GRANT USAGE, CREATE ON SCHEMA public TO airflow;


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
