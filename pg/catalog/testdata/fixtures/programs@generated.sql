CREATE TABLE gen (a int, b text GENERATED ALWAYS AS (a::text) STORED, d text DEFAULT 'x', e int DEFAULT 1);
