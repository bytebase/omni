CREATE DATABASE lab
CREATE TABLE lab.events (id UInt32, attrs Nested(k String, v UInt32)) ENGINE = MergeTree ORDER BY id
CREATE TABLE lab.calc (id UInt32, secret String, a String ALIAS concat(secret, '!'), b String ALIAS upper(a), m UInt32 MATERIALIZED id * 2) ENGINE = MergeTree ORDER BY id
CREATE TABLE lab.lookup (id UInt32, email String) ENGINE = Join(ANY, LEFT, id)
