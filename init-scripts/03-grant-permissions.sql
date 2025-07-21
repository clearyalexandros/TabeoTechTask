-- 03-grant-permissions.sql
-- Grant necessary permissions to application user

-- Grant full permissions to the citynext_user
GRANT ALL PRIVILEGES ON DATABASE citynext_appointments TO citynext_user;
GRANT ALL PRIVILEGES ON TABLE appointments TO citynext_user;
GRANT USAGE, SELECT ON SEQUENCE appointments_id_seq TO citynext_user;

-- Ensure future tables created by citynext_user have proper permissions
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO citynext_user;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO citynext_user;
