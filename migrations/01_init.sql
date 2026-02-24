CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    role VARCHAR(20) NOT NULL, -- 'sales', 'supervisor', 'farmer', 'distributor'
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS farmers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location GEOGRAPHY(Point, 4326), -- Koordinat GPS (PostGIS)
    crop_type VARCHAR(100),
    land_area DECIMAL(10,2),
    fertilizer_usage DECIMAL(10,2),
    competitor_product VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS visits (
    id SERIAL PRIMARY KEY,
    sales_id INT REFERENCES users(id),
    farmer_id INT REFERENCES farmers(id),
    photo_path VARCHAR(255), -- Path foto terkompresi
    visit_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    notes TEXT
);

CREATE TABLE IF NOT EXISTS distributors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location GEOGRAPHY(Point, 4326),
    credit_limit INT DEFAULT 0, -- 0 = COD, 30 = tempo 30 hari
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    distributor_id INT REFERENCES distributors(id),
    supervisor_id INT REFERENCES users(id),
    quantity INT,
    order_date DATE DEFAULT CURRENT_DATE,
    payment_status VARCHAR(20) DEFAULT 'unpaid' -- 'paid', 'unpaid'
);
