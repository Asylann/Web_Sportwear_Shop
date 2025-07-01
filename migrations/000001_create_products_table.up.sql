CREATE TABLE IF NOT EXISTS products (
       id SERIAL PRIMARY KEY,
       name TEXT NOT NULL UNIQUE ,
       description TEXT,
       category_id INT NOT NULL,
       size INT DEFAULT 0,
       price NUMERIC(10,2) NOT NULL DEFAULT 0,
       imageURL TEXT
);