CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);


CREATE TYPE link_type_enum AS ENUM ('general', 'payment', 'kyc', 'onboarding');

-- URL records
CREATE TABLE IF NOT EXISTS url_records (
    id VARCHAR(36) PRIMARY KEY,
    slug VARCHAR(255) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
    is_active BOOLEAN NOT NULL DEFAULT TRUE
    max_clicks INT,
    total_clicks INT NOT NULL DEFAULT 0,
    link_type link_type_enum NOT NULL DEFAULT 'general'
);

--click events (audit trail)
CREATE TABLE IF NOT EXISTS click_events (
    id VARCHAR(36) PRIMARY KEY,
    slug VARCHAR(20) NOT NULL,
    clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address VARCHAR(45),
    user_agent TEXT,
    was_valid BOOLEAN NOT NULL,
    rejection_reason VARCHAR(50)
);

--Indexes for performance
CREATE INDEX IF NOT EXISTS idx_url_records_slug ON url_records(slug);
CREATE INDEX IF NOT EXISTS idx_url_records_user_id ON url_records(user_id);
CREATE INDEX IF NOT EXISTS idx_click_events_slug ON click_events(slug);