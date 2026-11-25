-- The query span corpus of cosmosdb/analysis, judged by the Cosmos DB
-- emulator. One query per line, on the fixture and in the container the last
-- "-- fixture:" and "-- session:" lines name. testdata/README.md says how
-- each answer is read.
-- fixture: document
-- session: document

-- Property paths.
SELECT c.name FROM c
SELECT c.name, c.address.city FROM c
SELECT c.address FROM c
SELECT c.address.geo.lat FROM c
SELECT c.tags FROM c
SELECT c.id, c.pk FROM c
SELECT c._ts FROM c
SELECT c["name"] FROM c
SELECT c.address["city"] FROM c
SELECT c["address"]["geo"] FROM c
SELECT c["tags"][1]["name"] FROM c
SELECT c.tags[0] FROM c
SELECT c.tags[0].name FROM c
SELECT c.tags[0]["name"] FROM c
SELECT c.tags[0].skills FROM c
SELECT (c.name) FROM c

-- The whole document.
SELECT * FROM c
SELECT c FROM c
SELECT VALUE c FROM c
SELECT root FROM root
SELECT TOP 1 * FROM c

-- Aliases.
SELECT c.name AS n FROM c
SELECT c.name n FROM c
SELECT c.address.city AS name FROM c
SELECT c.name AS address, c.address AS name FROM c
SELECT c.address AS a FROM c
SELECT r.name FROM root r
SELECT root.name FROM root
SELECT x.name FROM mycontainer x
SELECT c FROM container c
SELECT c.name FROM container c

-- Items Cosmos DB numbers.
SELECT UPPER(c.name) FROM c
SELECT c.name, UPPER(c.name) FROM c
SELECT UPPER(c.name), c.name, LOWER(c.nick) FROM c
SELECT c.name, c.nick, c.n + 1, UPPER(c.name) AS u, LOWER(c.name) FROM c
SELECT c.name AS n, c.tags[0], [c.name] FROM c
SELECT -c.n FROM c
SELECT {"a": c.name} FROM c
SELECT IS_DEFINED(c.nick) FROM c
SELECT {"k": c.name}.k FROM c

-- Expressions and functions.
SELECT CONCAT(c.name, "-", c.address.city) AS label FROM c
SELECT c.name ?? c.nick AS who FROM c
SELECT IIF(c.n > 0, c.name, c.address.city) AS pick FROM c
SELECT c.n > 0 ? c.name : c.address.zip AS pick FROM c
SELECT {"a": c.name, "b": c.address.city} AS o FROM c
SELECT [c.name, c.address.zip] AS arr FROM c
SELECT [c.name, c.address.city][0] AS first FROM c
SELECT {"k": c.name}.k AS kk FROM c
SELECT (c.address ?? c.nick).city AS city FROM c
SELECT (c.n > 0 ? c.address : c.tags).zip AS zip FROM c
SELECT IIF(c.n > 0, c.address, c.tags).geo.lat FROM c
SELECT LOWER(c.name) AS l, UPPER(c.nick) AS u FROM c
SELECT REPLACE(c.name, "x", "y") AS r FROM c
SELECT ToString(c.n) AS s, c.n * 2 AS d FROM c
SELECT ARRAY_SLICE(c.phones, 0, 1) AS ph FROM c
SELECT ARRAY_CONCAT(c.phones, [c.name]) AS all FROM c
SELECT ARRAY_LENGTH(c.tags) AS n FROM c
SELECT 1 AS one, "k" AS k FROM c

-- SELECT VALUE.
SELECT VALUE c.name FROM c
SELECT VALUE c.address FROM c
SELECT VALUE {"n": c.name, "z": c.address.zip} FROM c
SELECT VALUE UPPER(c.name) FROM c
SELECT VALUE [c.name, c.n] FROM c
SELECT DISTINCT VALUE c.address.city FROM c

-- JOIN over arrays and FROM over a property.
SELECT t.name FROM c JOIN t IN c.tags
SELECT t FROM c JOIN t IN c.tags
SELECT c.name, t.name AS tag FROM c JOIN t IN c.tags
SELECT VALUE t FROM c JOIN t IN c.tags
SELECT t.name AS tag, s AS skill FROM c JOIN t IN c.tags JOIN s IN t.skills
SELECT t.skills[0] FROM c JOIN t IN c.tags
SELECT p FROM c JOIN p IN c.phones
SELECT c.name, t FROM c JOIN t IN c.phones
SELECT VALUE t.name FROM t IN c.tags
SELECT t.name FROM t IN c.tags WHERE t.name != ""
SELECT * FROM t IN c.tags
SELECT a.city FROM c.address a
SELECT a FROM c.address a
SELECT * FROM c.address
SELECT address.city FROM c.address
SELECT t.name FROM c JOIN t IN c.tags WHERE ARRAY_CONTAINS(t.skills, "«$.tags[0].skills[0]»")

-- Subqueries.
SELECT ARRAY(SELECT VALUE t.name FROM t IN c.tags) AS names FROM c
SELECT ARRAY(SELECT t.name FROM t IN c.tags) AS names FROM c
SELECT ARRAY(SELECT VALUE p FROM p IN c.phones) FROM c
SELECT ARRAY(SELECT VALUE s FROM s IN t.skills) AS sk FROM c JOIN t IN c.tags
SELECT VALUE ARRAY(SELECT VALUE t.name FROM t IN c.tags) FROM c
SELECT (SELECT VALUE COUNT(1) FROM t IN c.tags) AS cnt FROM c
SELECT (SELECT VALUE c.name) AS n FROM c
SELECT FIRST(SELECT VALUE t.name FROM t IN c.tags) AS f FROM c
SELECT EXISTS(SELECT VALUE t FROM t IN c.tags) AS has FROM c
SELECT c.name FROM c WHERE EXISTS(SELECT VALUE t FROM t IN c.tags WHERE t.name = "«$.tags[1].name»")
SELECT t.name FROM c JOIN (SELECT VALUE x FROM x IN c.tags) t
SELECT s.name FROM (SELECT c.name FROM c) s
SELECT s.city FROM (SELECT c.name, c.address.city FROM c) s
SELECT s FROM (SELECT c.name, c.address.city FROM c) s
SELECT * FROM (SELECT c.name, c.address.city FROM c) s
SELECT * FROM (SELECT VALUE c.address FROM c) s
SELECT s.zip FROM (SELECT VALUE c.address FROM c) s
SELECT VALUE s FROM (SELECT VALUE c.address FROM c) s
SELECT s.n FROM (SELECT c.name AS n FROM c) s
SELECT VALUE s.name FROM (SELECT c.name FROM c) s
SELECT s.k FROM (SELECT VALUE {"k": c.name} FROM c) s
SELECT s.o.k FROM (SELECT {"k": c.name} AS o FROM c) s
SELECT * FROM (SELECT VALUE {"a": c.name, "b": c.nick} FROM c) s

-- Aggregates, grouping and ordering.
SELECT COUNT(1) AS n FROM c
SELECT COUNT(1) FROM c
SELECT VALUE COUNT(1) FROM c
SELECT MAX(c.name) AS m FROM c
SELECT MIN(c.n) FROM c
SELECT SUM(c.n) AS total, AVG(c.n) AS mean FROM c
SELECT c.address.city, COUNT(1) AS n FROM c GROUP BY c.address.city
SELECT c.address.city, COUNT(1) FROM c GROUP BY c.address.city
SELECT t.name, COUNT(1) AS n FROM c JOIN t IN c.tags GROUP BY t.name
SELECT DISTINCT c.address.city FROM c
SELECT DISTINCT VALUE t.name FROM t IN c.tags
SELECT TOP 1 c.name FROM c ORDER BY c.name
SELECT c.name FROM c ORDER BY c.address.city DESC
SELECT c.address.city FROM c WHERE c.n > 0 OFFSET 0 LIMIT 1
