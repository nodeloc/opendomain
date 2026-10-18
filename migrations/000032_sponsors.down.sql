DROP TABLE IF EXISTS sponsors;

DELETE FROM system_settings WHERE setting_key IN ('sponsor_usdt_trc20_address', 'sponsor_github_url');
