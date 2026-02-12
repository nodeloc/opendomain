CREATE TABLE IF NOT EXISTS pages (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    content TEXT NOT NULL,
    category VARCHAR(50) NOT NULL, -- company, resources
    is_published BOOLEAN DEFAULT true,
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert default pages
INSERT INTO pages (title, slug, content, category, is_published, display_order) VALUES
('About Us', 'about-us', '<h2>About OpenDomain</h2><p>OpenDomain is a free subdomain service that provides users with easy-to-use domain management tools.</p>', 'company', true, 1),
('Contact', 'contact', '<h2>Contact Us</h2><p>Email: support@opendomain.com</p><p>We are here to help you!</p>', 'company', true, 2),
('Terms of Service', 'terms-of-service', '<h2>Terms of Service</h2><p>By using OpenDomain, you agree to our terms of service.</p>', 'company', true, 3),
('Privacy Policy', 'privacy-policy', '<h2>Privacy Policy</h2><p>We respect your privacy and protect your personal information.</p>', 'company', true, 4),
('Documentation', 'documentation', '<h2>Documentation</h2><p>Learn how to use OpenDomain with our comprehensive documentation.</p>', 'resources', true, 1),
('API Reference', 'api-reference', '<h2>API Reference</h2><p>Integrate OpenDomain into your applications with our API.</p>', 'resources', true, 2),
('FAQ', 'faq', '<h2>Frequently Asked Questions</h2><p>Find answers to common questions about OpenDomain.</p>', 'resources', true, 3),
('Support', 'support', '<h2>Support</h2><p>Need help? Our support team is ready to assist you.</p>', 'resources', true, 4);

CREATE INDEX idx_pages_category ON pages(category);
CREATE INDEX idx_pages_slug ON pages(slug);
