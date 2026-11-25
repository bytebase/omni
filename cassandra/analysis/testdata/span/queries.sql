-- The span corpus: one CQL query per line, run in the keyspace the last
-- "-- session:" line names, on the fixture the last "-- fixture:" line names.
-- fixture: base
-- session: shop

-- Plain columns, aliases and star.
SELECT id, name FROM users
SELECT name, id FROM users
SELECT id, name, name FROM users
SELECT NAME, Email FROM Users
SELECT * FROM users
SELECT * FROM orders
SELECT * FROM events
SELECT * FROM counters
SELECT * FROM pairs
SELECT id AS user_id, name AS n FROM users
SELECT id AS "ID", name AS "Name" FROM users
SELECT name AS email FROM users
SELECT email FROM users WHERE id = 1
SELECT tags, prefs, history, home FROM users
SELECT hits FROM counters
SELECT pair FROM pairs
SELECT customer_id, order_id, amount FROM orders WHERE customer_id = 1 ORDER BY order_id ASC
SELECT * FROM orders WHERE customer_id = 1 LIMIT 10
SELECT * FROM orders PER PARTITION LIMIT 1
SELECT status FROM orders WHERE status = 'open' ALLOW FILTERING
SELECT region FROM events WHERE tenant = 't' AND day = '2024-01-01'
SELECT name FROM users WHERE token(id) > 0
SELECT name FROM shop.users

-- DISTINCT.
SELECT DISTINCT id FROM users
SELECT DISTINCT customer_id FROM orders
SELECT DISTINCT tenant, day FROM events
SELECT DISTINCT tenant, day, region FROM events

-- JSON.
SELECT JSON id, name FROM users
SELECT JSON * FROM users
SELECT JSON * FROM orders
SELECT JSON DISTINCT id FROM users
SELECT JSON name AS n FROM users
SELECT JSON home.city FROM users
SELECT JSON count(*) FROM users
SELECT JSON "AccountId", "eMail" FROM "Accounts"

-- WRITETIME and TTL.
SELECT WRITETIME(name) FROM users
SELECT TTL(email) FROM users
SELECT writetime(name) AS w, ttl(name) AS t FROM users
SELECT id, WRITETIME(email), email FROM users
SELECT MAXWRITETIME(history) FROM users
SELECT writetime(region) FROM events

-- Aggregates.
SELECT count(*) FROM users
SELECT COUNT(1) FROM users
SELECT count(email) FROM users
SELECT max(score), min(score), avg(score), sum(score) FROM users
SELECT sum(score) AS total FROM users WHERE id IN (1, 2)
SELECT customer_id, sum(amount) FROM orders GROUP BY customer_id
SELECT customer_id, count(*) AS n FROM orders GROUP BY customer_id
SELECT tenant, day, max(seq) FROM events GROUP BY tenant, day

-- Functions, casts and arithmetic.
SELECT toJson(home) FROM users
SELECT toJson(email) AS j FROM users
SELECT toJson(prefs), toJson(tags) FROM users
SELECT token(id) FROM users
SELECT token(id), id FROM users
SELECT blobAsText(payload) FROM events
SELECT textAsBlob(email) FROM users
SELECT CAST(score AS text) FROM users
SELECT CAST(id AS bigint) AS big FROM users
SELECT toDate(created), toUnixTimestamp(created) FROM users
SELECT toTimestamp(order_id), toUnixTimestamp(order_id) FROM orders
SELECT now() FROM users
SELECT uuid(), id FROM users
SELECT id + 1 FROM users
SELECT score * 2 AS doubled FROM users
SELECT -score FROM users
SELECT id, (int) 1 FROM users
SELECT mask_default(email) FROM users
SELECT mask_inner(phone, 1, 2) FROM users
SELECT mask_hash(ssn) AS h FROM users
SELECT collection_count(tags) FROM users
SELECT map_keys(prefs) FROM users

-- Subselection of a user-defined type and a collection.
SELECT home.city FROM users
SELECT home.city AS c, home.zip FROM users
SELECT prefs['lang'] FROM users
SELECT prefs['a'..'m'] AS some FROM users
SELECT tags['red'] FROM users

-- Materialized views.
SELECT * FROM users_by_email
SELECT email, name FROM users_by_email WHERE email = 'a@b.c'
SELECT * FROM orders_by_status
SELECT status, amount FROM orders_by_status
SELECT count(*) FROM users_by_email
SELECT JSON * FROM users_by_email

-- Quoted and mixed-case names.
SELECT * FROM "Accounts"
SELECT "AccountId", "DisplayName" FROM "Accounts"
SELECT "eMail" AS "Mail" FROM "Accounts"
SELECT plain, PLAIN FROM "Accounts"
SELECT "DisplayName" AS display FROM shop."Accounts"
SELECT * FROM "Audit"."Log"
SELECT "Who", what FROM "Audit"."Log"

-- Other keyspaces.
SELECT * FROM crm.contacts
SELECT full_name, email FROM crm.contacts
SELECT id, owner FROM crm.contacts WHERE id = 5132b130-ae79-11e4-ab27-0800200c9a66
SELECT similarity_cosine(vec, [0.1, 0.2, 0.3]) FROM crm.contacts

-- session: crm
SELECT * FROM contacts
SELECT email FROM contacts
SELECT name, email FROM shop.users
SELECT * FROM shop.users_by_email
