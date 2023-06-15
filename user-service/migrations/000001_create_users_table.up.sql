CREATE TABLE IF NOT EXISTS users
    (
    id UUID, 
    name TEXT, 
    last_name TEXT,
    email TEXT,
    username TEXT,
    password TEXT,
    refresh_token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(), 
    updated_at TIMESTAMP DEFAULT NOW(), 
    deleted_at TIMESTAMP);