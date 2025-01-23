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

-- Grant privileges to the user
GRANT ALL PRIVILEGES ON DATABASE airflow TO airflow;

-- Grant privileges to the user on the public schema
\c airflow
GRANT ALL PRIVILEGES ON SCHEMA public TO airflow;
GRANT CREATE ON SCHEMA public TO airflow;
