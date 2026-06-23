-- SELECT statements with JSON output

-- Select all columns as JSON
SELECT JSON * FROM users WHERE id = 123e4567-e89b-12d3-a456-426614174000;

-- Select specific columns as JSON
SELECT JSON id, name, age FROM users LIMIT 10;

-- toJson function for specific column
SELECT id, toJson(profile) FROM users WHERE age > 25;