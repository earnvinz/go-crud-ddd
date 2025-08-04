CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    name_th TEXT NOT NULL,
    name_en TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW(),
    updated_by TEXT NOT NULL
);
