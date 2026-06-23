-- Create table with VECTOR data type (Cassandra 5.0+)
CREATE TABLE products (
    id uuid PRIMARY KEY,
    name text,
    description text,
    embedding VECTOR<float, 1536>,
    price decimal
);

-- Create table with DURATION type
CREATE TABLE events (
    id uuid PRIMARY KEY,
    name text,
    event_duration DURATION,
    created_at timestamp
);