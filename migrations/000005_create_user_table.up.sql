CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role_id INT NOT NULL
        REFERENCES roles(id)
        ON DELETE RESTRICT
);