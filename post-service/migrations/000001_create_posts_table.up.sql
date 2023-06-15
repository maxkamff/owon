CREATE TABLE IF NOT EXISTS posts
    (
        id SERIAL PRIMARY KEY,
        title TEXT,
        description TEXT,
        user_id INT,
        created_at TIMESTAMP DEFAULT NOW(), 
        updated_at TIMESTAMP DEFAULT NOW(), 
        deleted_at TIMESTAMP
    )