-- Announcements table
CREATE TABLE IF NOT EXISTS announcements (
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  type VARCHAR(20) NOT NULL CHECK (type IN ('general', 'maintenance', 'update', 'important')),
  priority INT DEFAULT 0, -- Higher number = higher priority
  is_published BOOLEAN DEFAULT FALSE,
  published_at TIMESTAMP,
  author_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  views INT DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX idx_announcements_type ON announcements(type);
CREATE INDEX idx_announcements_published ON announcements(is_published, published_at DESC);
CREATE INDEX idx_announcements_priority ON announcements(priority DESC);
