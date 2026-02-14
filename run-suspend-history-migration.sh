#!/bin/bash
# Run the suspend_history migration

# Load environment variables
source .env 2>/dev/null || true

# Set defaults if not set
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-opendomain}
DB_PASSWORD=${DB_PASSWORD:-opendomain123}
DB_NAME=${DB_NAME:-opendomain_dev}

echo "Creating suspend_history table..."

# Check if running in Docker
if docker ps | grep -q opendomain-postgres; then
    docker exec -i opendomain-postgres psql -U $DB_USER -d $DB_NAME <<EOF
-- Create suspend_history table
CREATE TABLE IF NOT EXISTS suspend_history (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_suspend_history_domain_id ON suspend_history(domain_id);
CREATE INDEX IF NOT EXISTS idx_suspend_history_created_at ON suspend_history(created_at DESC);
EOF
elif docker ps | grep -q opendomain-postgres-dev; then
    docker exec -i opendomain-postgres-dev psql -U $DB_USER -d $DB_NAME <<EOF
-- Create suspend_history table
CREATE TABLE IF NOT EXISTS suspend_history (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_suspend_history_domain_id ON suspend_history(domain_id);
CREATE INDEX IF NOT EXISTS idx_suspend_history_created_at ON suspend_history(created_at DESC);
EOF
else
    # Try direct connection
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
-- Create suspend_history table
CREATE TABLE IF NOT EXISTS suspend_history (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES domains(id) ON DELETE CASCADE,
    reason TEXT NOT NULL,
    details TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_suspend_history_domain_id ON suspend_history(domain_id);
CREATE INDEX IF NOT EXISTS idx_suspend_history_created_at ON suspend_history(created_at DESC);
EOF
fi

if [ $? -eq 0 ]; then
    echo "✓ Successfully created suspend_history table and indexes"
else
    echo "✗ Failed to create table. Please check your database connection."
    exit 1
fi
