-- 01-create-user.sql
-- Create application database user with proper permissions
-- Depends on: 00-create-database.sql

-- Create the citynext_user if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'citynext_user') THEN
        CREATE USER citynext_user WITH PASSWORD 'citynext_password' SUPERUSER;
        GRANT CONNECT ON DATABASE citynext_appointments TO citynext_user;
        GRANT ALL PRIVILEGES ON DATABASE citynext_appointments TO citynext_user;
    END IF;
END
$$;
