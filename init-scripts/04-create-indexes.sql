-- 04-create-indexes.sql
-- Create indexes for better query performance

-- Create an index for better performance on appointment date queries
CREATE INDEX IF NOT EXISTS idx_appointments_visit_date ON appointments(visit_date);

-- Create an index for name searches (optional, for future use)
CREATE INDEX IF NOT EXISTS idx_appointments_names ON appointments(first_name, last_name);
