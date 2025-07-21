-- 00-create-database.sql
-- Create the main application database
-- Note: This database is also created by POSTGRES_DB environment variable,
-- but having it explicit here makes the setup clearer and more maintainable

-- Create database if it doesn't exist
SELECT 'CREATE DATABASE citynext_appointments'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'citynext_appointments')\gexec
