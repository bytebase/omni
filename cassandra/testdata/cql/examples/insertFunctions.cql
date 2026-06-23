-- INSERT statements with CQL functions

-- UUID generation function
INSERT INTO events (id, created_at) VALUES (uuid(), now());

-- Current time UUID and timestamp functions
INSERT INTO events (id, created_at) VALUES (currentTimeuuid(), currentTimestamp());

-- Date/time component functions  
INSERT INTO logs (id, log_date, log_time) VALUES (uuid(), currentDate(), currentTime());

-- fromJson function for complex data
INSERT INTO users (id, profile) VALUES (uuid(), fromJson('{"location": "NYC", "interests": ["music", "sports"]}'));