-- Migration: Create users table for Aadhaar User Service
-- Version: 001
-- Description: Initial schema for storing Aadhaar application user details

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aadhaar_application_id VARCHAR(14) NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone VARCHAR(10) NOT NULL,
    address VARCHAR(500) NOT NULL,
    date_of_birth VARCHAR(10) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create unique indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_aadhaar_application_id ON users(aadhaar_application_id);

-- Create indexes for sorting and searching
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Add constraint for gender values
ALTER TABLE users ADD CONSTRAINT chk_gender CHECK (gender IN ('male', 'female', 'other'));

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Comments for documentation
COMMENT ON TABLE users IS 'Stores Aadhaar application user details';
COMMENT ON COLUMN users.id IS 'Unique identifier (UUID)';
COMMENT ON COLUMN users.aadhaar_application_id IS '14-character Aadhaar application ID';
COMMENT ON COLUMN users.name IS 'Full name of the applicant';
COMMENT ON COLUMN users.email IS 'Email address (unique)';
COMMENT ON COLUMN users.phone IS '10-digit phone number';
COMMENT ON COLUMN users.address IS 'Residential address';
COMMENT ON COLUMN users.date_of_birth IS 'Date of birth (YYYY-MM-DD format)';
COMMENT ON COLUMN users.gender IS 'Gender (male/female/other)';
