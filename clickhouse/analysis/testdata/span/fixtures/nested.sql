CREATE DATABASE lab
CREATE TABLE lab.events (id UInt32, attrs Nested(k String, v UInt32)) ENGINE = MergeTree ORDER BY id
