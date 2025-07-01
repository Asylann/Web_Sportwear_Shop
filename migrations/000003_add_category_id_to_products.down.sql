ALTER TABLE products
    DROP CONSTRAINT IF EXISTS products_category_fkey;

ALTER TABLE products
    DROP COLUMN IF EXISTS category_id;