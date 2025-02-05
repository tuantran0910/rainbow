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
CREATE DATABASE rainbow_test;

-- Create the user with a password
CREATE ROLE rainbow WITH LOGIN PASSWORD 'R&inb0w2024!Data';

-- Grant necessary privileges to the user
GRANT CONNECT ON DATABASE rainbow TO rainbow;
GRANT TEMPORARY ON DATABASE rainbow TO rainbow;

-- Switch to the rainbow database and grant schema permissions
\c rainbow
GRANT USAGE, CREATE ON SCHEMA public TO rainbow;
