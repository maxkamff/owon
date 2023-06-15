CREATE TABLE IF NOT EXISTS comments
    (
        id SERIAL PRIMARY KEY,
        user_id INT,
        post_id INT,
        description TEXT,
        liked BOOLEAN,
        created_at TIMESTAMP DEFAULT NOW(), 
        updated_at TIMESTAMP DEFAULT NOW(), 
        deleted_at TIMESTAMP
    )