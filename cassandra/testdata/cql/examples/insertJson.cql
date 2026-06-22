-- INSERT statements with JSON support

-- Basic JSON insert
INSERT INTO users JSON '{"id": "123e4567-e89b-12d3-a456-426614174000", "name": "John Doe", "age": 30}';

-- JSON insert with DEFAULT UNSET clause
INSERT INTO users JSON '{"id": "123e4567-e89b-12d3-a456-426614174000", "name": "Jane Doe"}' DEFAULT UNSET;

-- fromJson function in VALUES clause
INSERT INTO users (id, profile) VALUES (uuid(), fromJson('{"location": "NYC", "interests": ["music", "sports"]}'));