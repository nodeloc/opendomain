-- Add pricing fields to root_domains table
ALTER TABLE root_domains ADD COLUMN price_per_year DECIMAL(10,2) DEFAULT 0.00;
ALTER TABLE root_domains ADD COLUMN lifetime_price DECIMAL(10,2);
ALTER TABLE root_domains ADD COLUMN is_free BOOLEAN DEFAULT TRUE;

CREATE INDEX idx_root_domains_is_free ON root_domains(is_free);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
  order_number VARCHAR(32) UNIQUE NOT NULL,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,

  -- Domain information
  subdomain VARCHAR(63) NOT NULL,
  root_domain_id INT NOT NULL REFERENCES root_domains(id) ON DELETE RESTRICT,
  full_domain VARCHAR(255) NOT NULL,

  -- Pricing details
  years INT NOT NULL CHECK (years >= 1 AND years <= 10),
  is_lifetime BOOLEAN DEFAULT FALSE,
  base_price DECIMAL(10,2) NOT NULL,
  discount_amount DECIMAL(10,2) DEFAULT 0.00,
  final_price DECIMAL(10,2) NOT NULL,

  -- Coupon information
  coupon_id BIGINT REFERENCES coupons(id) ON DELETE SET NULL,
  coupon_code VARCHAR(50),

  -- Order status
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'cancelled', 'refunded', 'expired')),

  -- Domain activation
  domain_id BIGINT REFERENCES domains(id) ON DELETE SET NULL,

  -- Timestamps
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  paid_at TIMESTAMP,
  expires_at TIMESTAMP NOT NULL,

  -- Notes
  notes TEXT
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_created_at ON orders(created_at);
CREATE INDEX idx_orders_expires_at ON orders(expires_at);

-- Payments table
CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,

  -- NodeLoc payment details
  transaction_id VARCHAR(100) UNIQUE,
  nodeloc_payment_id VARCHAR(50) NOT NULL,

  -- Payment information
  amount DECIMAL(10,2) NOT NULL,
  currency VARCHAR(3) DEFAULT 'CNY',

  -- Payment status
  status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'refunded')),

  -- Gateway response
  gateway_response TEXT,
  signature VARCHAR(255),

  -- Callback information
  callback_received_at TIMESTAMP,
  callback_ip VARCHAR(45),

  -- Timestamps
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  completed_at TIMESTAMP
);

CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_transaction_id ON payments(transaction_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at);

-- Payment configuration table
CREATE TABLE IF NOT EXISTS payment_configs (
  id SERIAL PRIMARY KEY,
  provider VARCHAR(20) NOT NULL DEFAULT 'nodeloc',
  payment_id VARCHAR(50) NOT NULL,
  secret_key VARCHAR(255) NOT NULL,
  callback_url VARCHAR(255),
  is_active BOOLEAN DEFAULT TRUE,
  is_test_mode BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Add domain_id to coupon_usage index for tracking
CREATE INDEX idx_coupon_usage_domain_id ON coupon_usage(domain_id);
