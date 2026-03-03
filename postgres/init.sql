-- Create schemas
CREATE SCHEMA IF NOT EXISTS item_schema;
CREATE SCHEMA IF NOT EXISTS order_schema;
CREATE SCHEMA IF NOT EXISTS payment_schema;

-- item_schema.items
CREATE TABLE IF NOT EXISTS item_schema.items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- order_schema.orders
CREATE TABLE IF NOT EXISTS order_schema.orders (
    id SERIAL PRIMARY KEY,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    customer_id VARCHAR(50),
    status VARCHAR(50) DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT NOW()
);

-- payment_schema.payments
CREATE TABLE IF NOT EXISTS payment_schema.payments (
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    method VARCHAR(50),
    status VARCHAR(50) DEFAULT 'SUCCESS',
    created_at TIMESTAMP DEFAULT NOW()
);

-- Basic grants (optional if connected as postgres user, but good practice)
GRANT ALL PRIVILEGES ON SCHEMA item_schema TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA item_schema TO postgres;

GRANT ALL PRIVILEGES ON SCHEMA order_schema TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA order_schema TO postgres;

GRANT ALL PRIVILEGES ON SCHEMA payment_schema TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA payment_schema TO postgres;
