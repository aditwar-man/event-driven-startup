CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    tier VARCHAR(50) DEFAULT 'free',
    ai_description_quota_used INTEGER DEFAULT 0,
    ai_description_quota_limit INTEGER DEFAULT 5,
    ai_video_quota_used INTEGER DEFAULT 0,
    ai_video_quota_limit INTEGER DEFAULT 0,
    auto_posting_quota_used INTEGER DEFAULT 0,
    auto_posting_quota_limit INTEGER DEFAULT 5,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_tier ON users(tier);
