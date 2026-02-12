-- Add default currency setting
INSERT INTO system_settings (setting_key, setting_value, description, created_at, updated_at)
VALUES ('currency_symbol', 'NL', 'Currency symbol displayed for all prices', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (setting_key) DO NOTHING;
