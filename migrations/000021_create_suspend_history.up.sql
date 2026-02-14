-- Create suspend_history table
CREATE TABLE IF NOT EXISTS suspend_history (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_suspend_history_domain_id ON suspend_history(domain_id);
CREATE INDEX idx_suspend_history_created_at ON suspend_history(created_at DESC);
