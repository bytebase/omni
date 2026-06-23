-- SELECT statements with CQL functions

-- toJson function for JSON output
SELECT id, toJson(profile) FROM users WHERE age > 25;

-- JSON format selection (all columns)
SELECT JSON * FROM users WHERE id = 123e4567-e89b-12d3-a456-426614174000;

-- JSON format selection (specific columns)
SELECT JSON id, name, age FROM users LIMIT 10;

-- dateOf function to extract timestamp from timeuuid
SELECT id, dateOf(created_id) FROM events WHERE id = 550e8400-e29b-41d4-a716-446655440000;

-- unixTimestampOf function to get UNIX timestamp
SELECT id, unixTimestampOf(created_id) FROM events LIMIT 5;