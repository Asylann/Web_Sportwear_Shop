CREATE TABLE IF NOT EXISTS roles(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE ,
    description TEXT NOT NULL
);

INSERT INTO roles (name, description) VALUES
('customer', 'End user who can browse and purchase products'),
('seller',   'User who can manage (create/update/delete) their own products'),
('admin',    'User with full access to all data and operations');