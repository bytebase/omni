-- @name: DELETE missing FROM and table — WHERE only
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE WHERE id = 1

-- @name: DELETE missing FROM and table — LIMIT only
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE LIMIT 10

-- @name: DELETE missing FROM and table — ORDER BY LIMIT
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE ORDER BY id LIMIT 1

-- @name: DELETE FROM t missing BY after ORDER
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE FROM t ORDER id

-- @name: DELETE FROM t WHERE — single-table delete accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE FROM t WHERE id = 1

-- @name: DELETE multi-table syntax 1 accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id

-- @name: DELETE multi-table USING syntax accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.3
DELETE FROM t1 USING t1 JOIN t2

-- @name: UPDATE missing SET — WHERE only
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.5
UPDATE t WHERE id = 1

-- @name: UPDATE missing SET — LIMIT only
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.5
UPDATE t LIMIT 10

-- @name: UPDATE missing SET — ORDER BY only
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.5
UPDATE t ORDER BY id

-- @name: UPDATE SET x ORDER missing BY
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.5
UPDATE t SET x = 1 ORDER id

-- @name: UPDATE SET WHERE — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.5
UPDATE t SET x = 1 WHERE id = 1

-- @name: UPDATE SET without WHERE — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.5
UPDATE t SET x = 1

-- @name: INSERT INTO with columns but no value source
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t (id, name)

-- @name: INSERT INTO with no value source at all
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t

-- @name: INSERT INTO with WHERE — not a valid value source
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t (id) WHERE 1=1

-- @name: REPLACE INTO with columns but no value source
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
REPLACE INTO t (id, name)

-- @name: REPLACE INTO with no value source at all
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
REPLACE INTO t

-- @name: INSERT ON DUPLICATE incomplete — missing KEY UPDATE
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t VALUES (1) ON DUPLICATE

-- @name: INSERT ON DUPLICATE KEY incomplete — missing UPDATE
-- @valid: false
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t VALUES (1) ON DUPLICATE KEY

-- @name: INSERT INTO VALUES — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t VALUES (1, 'a')

-- @name: INSERT INTO SET — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t SET id = 1

-- @name: INSERT INTO SELECT — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t SELECT * FROM t2

-- @name: INSERT INTO columns VALUES — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t (id) VALUES (1)

-- @name: REPLACE INTO VALUES — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.4
REPLACE INTO t VALUES (1, 'a')

-- @name: INSERT ON DUPLICATE KEY UPDATE — accepted
-- @valid: true
-- @source: SCENARIOS-mysql-strict.md section 2.4
INSERT INTO t VALUES (1) ON DUPLICATE KEY UPDATE x = 1
