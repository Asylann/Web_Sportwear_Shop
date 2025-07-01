ALTER TABLE products
    ADD COLUMN category_id INT;

ALTER TABLE products
    ADD CONSTRAINT products_category_fkey
        FOREIGN KEY (category_id)
            REFERENCES categories(id)
            ON DELETE SET NULL;