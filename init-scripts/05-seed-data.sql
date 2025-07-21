-- 05-seed-data.sql
-- Insert initial test/seed data (optional, for development/testing)

-- Insert a test record for verification
INSERT INTO appointments (first_name, last_name, visit_date) 
VALUES ('Test', 'User', '2075-12-25')
ON CONFLICT (visit_date) DO NOTHING;

-- Add more seed data if needed for testing
INSERT INTO appointments (first_name, last_name, visit_date) 
VALUES ('Demo', 'Person', '2075-06-15')
ON CONFLICT (visit_date) DO NOTHING;
