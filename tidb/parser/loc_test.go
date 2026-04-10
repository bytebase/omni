package parser

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// LocViolation records a single location validation failure.
type LocViolation struct {
	Path    string // e.g. "Items[0](SelectStmt).Columns[0](ColumnRef)"
	NodeTag string // node type name
	Start   int
	End     int
	Reason  string
}

func (v LocViolation) String() string {
	return fmt.Sprintf("%s [%s]: Start=%d End=%d — %s", v.Path, v.NodeTag, v.Start, v.End, v.Reason)
}

// CheckLocations parses sql via Parse(), recursively walks the AST using
// reflection, and returns all Loc violations where Start >= 0 but End <= Start.
func CheckLocations(t *testing.T, sql string) []LocViolation {
	t.Helper()

	result, err := Parse(sql)
	if err != nil {
		t.Fatalf("CheckLocations Parse(%q): %v", sql, err)
	}

	var violations []LocViolation
	if result != nil {
		for i, item := range result.Items {
			path := fmt.Sprintf("Items[%d]", i)
			walkNodeLocs(reflect.ValueOf(item), path, &violations)
		}
	}
	return violations
}

// walkNodeLocs recursively walks a reflected AST value, checking every Loc field.
func walkNodeLocs(v reflect.Value, path string, violations *[]LocViolation) {
	// Dereference pointers and interfaces.
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		typeName := v.Type().Name()

		// Check if this struct has a Loc field of type ast.Loc.
		locField := v.FieldByName("Loc")
		if locField.IsValid() && locField.Type() == reflect.TypeOf(ast.Loc{}) {
			loc := locField.Interface().(ast.Loc)
			if loc.Start >= 0 && loc.End <= loc.Start {
				*violations = append(*violations, LocViolation{
					Path:    path,
					NodeTag: typeName,
					Start:   loc.Start,
					End:     loc.End,
					Reason:  "Start >= 0 but End <= Start",
				})
			}
		}

		// Recurse into all fields.
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			if field.Name == "Loc" {
				continue // already checked
			}
			childPath := path + "." + field.Name
			walkNodeLocs(v.Field(i), childPath, violations)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			elemPath := fmt.Sprintf("%s[%d]", path, i)
			// Add type name for interface elements.
			actual := elem
			for actual.Kind() == reflect.Ptr || actual.Kind() == reflect.Interface {
				if actual.IsNil() {
					break
				}
				actual = actual.Elem()
			}
			if actual.IsValid() && actual.Kind() == reflect.Struct {
				elemPath = fmt.Sprintf("%s[%d](%s)", path, i, actual.Type().Name())
			}
			walkNodeLocs(elem, elemPath, violations)
		}
	}
}

// TestLocAudit runs CheckLocations against a large sample of SQL statements
// and reports all Loc.End violations grouped by AST node type.
func TestLocAudit(t *testing.T) {
	// Comprehensive set of MySQL SQL test cases covering all major statement types.
	testCases := []string{
		// SELECT basics
		"SELECT 1",
		"SELECT 1, 2, 3",
		"SELECT a, b FROM t",
		"SELECT a AS x, b AS y FROM t",
		"SELECT * FROM t",
		"SELECT DISTINCT a FROM t",
		"SELECT ALL a FROM t",
		"SELECT a FROM t WHERE a > 1",
		"SELECT a FROM t WHERE a = 1 AND b = 2",
		"SELECT a FROM t WHERE a = 1 OR b = 2",
		"SELECT a FROM t WHERE NOT a = 1",

		// SELECT with GROUP BY, HAVING, ORDER BY, LIMIT
		"SELECT a, COUNT(*) FROM t GROUP BY a",
		"SELECT a, COUNT(*) FROM t GROUP BY a HAVING COUNT(*) > 1",
		"SELECT a FROM t ORDER BY a ASC",
		"SELECT a FROM t ORDER BY a DESC",
		"SELECT a FROM t ORDER BY a ASC, b DESC",
		"SELECT a FROM t LIMIT 10",
		"SELECT a FROM t LIMIT 10 OFFSET 5",
		"SELECT a FROM t LIMIT 5, 10",

		// SELECT with JOINs
		"SELECT a FROM t1 JOIN t2 ON t1.id = t2.id",
		"SELECT a FROM t1 LEFT JOIN t2 ON t1.id = t2.id",
		"SELECT a FROM t1 RIGHT JOIN t2 ON t1.id = t2.id",
		"SELECT a FROM t1 INNER JOIN t2 ON t1.id = t2.id",
		"SELECT a FROM t1 CROSS JOIN t2",
		"SELECT a FROM t1 NATURAL JOIN t2",
		"SELECT a FROM t1 JOIN t2 USING (id)",
		"SELECT a FROM t1, t2 WHERE t1.id = t2.id",

		// SELECT with subqueries
		"SELECT * FROM (SELECT 1 AS a) AS sub",
		"SELECT a FROM t WHERE a IN (SELECT b FROM t2)",
		"SELECT EXISTS (SELECT 1 FROM t)",
		"SELECT (SELECT MAX(a) FROM t) AS m",

		// UNION
		"SELECT 1 UNION SELECT 2",
		"SELECT 1 UNION ALL SELECT 2",
		"SELECT 1 UNION SELECT 2 UNION SELECT 3",

		// SELECT FOR UPDATE
		"SELECT a FROM t FOR UPDATE",
		"SELECT a FROM t LOCK IN SHARE MODE",

		// Expressions
		"SELECT 1 + 2",
		"SELECT a - b FROM t",
		"SELECT a * b FROM t",
		"SELECT a / b FROM t",
		"SELECT a % b FROM t",
		"SELECT -a FROM t",
		"SELECT a BETWEEN 1 AND 10 FROM t",
		"SELECT a NOT BETWEEN 1 AND 10 FROM t",
		"SELECT a IN (1, 2, 3) FROM t",
		"SELECT a NOT IN (1, 2, 3) FROM t",
		"SELECT a LIKE 'foo%' FROM t",
		"SELECT a NOT LIKE 'foo%' FROM t",
		"SELECT a IS NULL FROM t",
		"SELECT a IS NOT NULL FROM t",
		"SELECT CASE WHEN a > 0 THEN 'pos' ELSE 'neg' END FROM t",
		"SELECT CASE a WHEN 1 THEN 'one' WHEN 2 THEN 'two' END FROM t",
		"SELECT CAST(a AS CHAR) FROM t",
		"SELECT CONVERT(a, CHAR) FROM t",
		"SELECT COALESCE(a, b) FROM t",
		"SELECT IF(a > 0, 'yes', 'no') FROM t",
		"SELECT IFNULL(a, 0) FROM t",
		"SELECT NULLIF(a, 0) FROM t",

		// Function calls
		"SELECT COUNT(*) FROM t",
		"SELECT COUNT(DISTINCT a) FROM t",
		"SELECT SUM(a) FROM t",
		"SELECT AVG(a) FROM t",
		"SELECT MIN(a) FROM t",
		"SELECT MAX(a) FROM t",
		"SELECT CONCAT(a, b) FROM t",
		"SELECT UPPER(a) FROM t",
		"SELECT NOW()",
		"SELECT CURDATE()",

		// INSERT
		"INSERT INTO t VALUES (1)",
		"INSERT INTO t (a, b) VALUES (1, 2)",
		"INSERT INTO t VALUES (1, 2), (3, 4)",
		"INSERT INTO t (a) VALUES (1) ON DUPLICATE KEY UPDATE a = 1",
		"INSERT INTO t SET a = 1, b = 2",
		"INSERT INTO t SELECT * FROM t2",
		"INSERT IGNORE INTO t VALUES (1)",
		"REPLACE INTO t VALUES (1, 2)",

		// UPDATE
		"UPDATE t SET a = 1",
		"UPDATE t SET a = 1, b = 2 WHERE c = 3",
		"UPDATE t1 JOIN t2 ON t1.id = t2.id SET t1.a = t2.b",
		"UPDATE t SET a = a + 1 ORDER BY b LIMIT 10",

		// DELETE
		"DELETE FROM t",
		"DELETE FROM t WHERE a = 1",
		"DELETE FROM t ORDER BY a LIMIT 10",
		"DELETE t1 FROM t1 JOIN t2 ON t1.id = t2.id",

		// CREATE TABLE
		"CREATE TABLE t (id INT)",
		"CREATE TABLE t (id INT PRIMARY KEY)",
		"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY)",
		"CREATE TABLE t (id INT, name VARCHAR(100))",
		"CREATE TABLE t (id INT, name VARCHAR(100) NOT NULL DEFAULT '')",
		"CREATE TABLE t (id INT, name TEXT, INDEX idx_name (name(10)))",
		"CREATE TABLE t (id INT, UNIQUE KEY uk_id (id))",
		"CREATE TABLE t (id INT, a INT, FOREIGN KEY (a) REFERENCES t2 (id))",
		"CREATE TABLE IF NOT EXISTS t (id INT)",
		"CREATE TEMPORARY TABLE t (id INT)",
		"CREATE TABLE t (id INT) ENGINE=InnoDB",
		"CREATE TABLE t (id INT) DEFAULT CHARSET=utf8mb4",
		"CREATE TABLE t (id INT, name VARCHAR(100), PRIMARY KEY (id))",
		"CREATE TABLE t (id INT) COMMENT 'test table'",
		"CREATE TABLE t (id INT, CHECK (id > 0))",
		"CREATE TABLE t LIKE t2",
		"CREATE TABLE t AS SELECT * FROM t2",
		"CREATE TABLE t (id INT, ts TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)",

		// ALTER TABLE
		"ALTER TABLE t ADD COLUMN c INT",
		"ALTER TABLE t ADD COLUMN c INT AFTER id",
		"ALTER TABLE t ADD COLUMN c INT FIRST",
		"ALTER TABLE t DROP COLUMN c",
		"ALTER TABLE t MODIFY COLUMN c VARCHAR(200)",
		"ALTER TABLE t CHANGE COLUMN c d INT",
		"ALTER TABLE t RENAME TO t2",
		"ALTER TABLE t ADD INDEX idx_c (c)",
		"ALTER TABLE t ADD UNIQUE INDEX uk_c (c)",
		"ALTER TABLE t DROP INDEX idx_c",
		"ALTER TABLE t ADD PRIMARY KEY (id)",
		"ALTER TABLE t DROP PRIMARY KEY",
		"ALTER TABLE t ADD FOREIGN KEY (a) REFERENCES t2 (id)",
		"ALTER TABLE t ADD CONSTRAINT chk CHECK (a > 0)",
		"ALTER TABLE t DROP FOREIGN KEY fk_name",
		"ALTER TABLE t ENGINE = InnoDB",
		"ALTER TABLE t DEFAULT CHARACTER SET utf8mb4",
		"ALTER TABLE t CONVERT TO CHARACTER SET utf8mb4",
		"ALTER TABLE t ADD COLUMN a INT, ADD COLUMN b INT",
		"ALTER TABLE t RENAME INDEX idx_old TO idx_new",

		// CREATE INDEX
		"CREATE INDEX idx ON t (a)",
		"CREATE UNIQUE INDEX idx ON t (a)",
		"CREATE INDEX idx ON t (a, b)",
		"CREATE INDEX idx ON t (a(10))",
		"CREATE FULLTEXT INDEX idx ON t (a)",
		"CREATE SPATIAL INDEX idx ON t (a)",

		// DROP statements
		"DROP TABLE t",
		"DROP TABLE IF EXISTS t",
		"DROP TABLE t1, t2",
		"DROP TABLE t CASCADE",
		"DROP INDEX idx ON t",
		"DROP DATABASE db1",
		"DROP DATABASE IF EXISTS db1",

		// CREATE/DROP VIEW
		"CREATE VIEW v AS SELECT * FROM t",
		"CREATE OR REPLACE VIEW v AS SELECT a, b FROM t",
		"CREATE VIEW v (x, y) AS SELECT a, b FROM t",
		"DROP VIEW v",
		"DROP VIEW IF EXISTS v",

		// CREATE DATABASE
		"CREATE DATABASE db1",
		"CREATE DATABASE IF NOT EXISTS db1",
		"CREATE DATABASE db1 DEFAULT CHARACTER SET utf8mb4",

		// TRUNCATE
		"TRUNCATE TABLE t",
		"TRUNCATE t",

		// RENAME TABLE
		"RENAME TABLE t1 TO t2",
		"RENAME TABLE t1 TO t2, t3 TO t4",

		// SET / SHOW
		"SET @a = 1",
		"SET NAMES utf8mb4",
		"SET CHARACTER SET utf8mb4",
		"SET GLOBAL max_connections = 100",
		"SET SESSION sql_mode = 'STRICT_TRANS_TABLES'",
		"SHOW DATABASES",
		"SHOW TABLES",
		"SHOW COLUMNS FROM t",
		"SHOW INDEX FROM t",
		"SHOW CREATE TABLE t",
		"SHOW TABLE STATUS",
		"SHOW VARIABLES LIKE 'max%'",
		"SHOW PROCESSLIST",
		"SHOW WARNINGS",
		"SHOW ERRORS",

		// USE
		"USE db1",

		// Transaction
		"BEGIN",
		"START TRANSACTION",
		"COMMIT",
		"ROLLBACK",
		"SAVEPOINT sp1",
		"ROLLBACK TO SAVEPOINT sp1",
		"RELEASE SAVEPOINT sp1",

		// EXPLAIN / DESCRIBE
		"EXPLAIN SELECT 1",
		"DESCRIBE t",

		// PREPARE / EXECUTE / DEALLOCATE
		"PREPARE stmt FROM 'SELECT ?'",
		"EXECUTE stmt USING @a",
		"DEALLOCATE PREPARE stmt",

		// GRANT / REVOKE
		"GRANT SELECT ON db1.* TO 'user'@'localhost'",
		"GRANT ALL PRIVILEGES ON *.* TO 'user'@'localhost'",
		"REVOKE SELECT ON db1.* FROM 'user'@'localhost'",

		// LOCK / UNLOCK
		"LOCK TABLES t WRITE",
		"LOCK TABLES t1 READ, t2 WRITE",
		"UNLOCK TABLES",

		// FLUSH
		"FLUSH TABLES",
		"FLUSH PRIVILEGES",

		// CREATE TRIGGER
		"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET NEW.a = 1",
		"CREATE TRIGGER tr AFTER UPDATE ON t FOR EACH ROW SET NEW.a = OLD.a + 1",

		// CREATE PROCEDURE / FUNCTION (simple cases)
		"DROP PROCEDURE IF EXISTS p",
		"DROP FUNCTION IF EXISTS f",

		// CALL
		"CALL my_proc()",
		"CALL my_proc(1, 2)",

		// DO
		"DO 1",
		"DO SLEEP(1)",

		// HANDLER
		"HANDLER t OPEN",
		"HANDLER t CLOSE",

		// LOAD DATA
		"LOAD DATA INFILE '/tmp/data.txt' INTO TABLE t",

		// Multi-statement
		"SELECT 1; SELECT 2",
		"INSERT INTO t VALUES (1); UPDATE t SET a = 2",
	}

	// Track violations grouped by NodeTag
	violationsByTag := make(map[string][]LocViolation)
	totalViolations := 0
	totalStatements := 0
	statementsWithViolations := 0

	for _, sql := range testCases {
		violations := func() []LocViolation {
			// Some statements may fail to parse (e.g., if our parser doesn't support them).
			// Use a recovery mechanism.
			result, err := Parse(sql)
			if err != nil {
				return nil
			}
			var vs []LocViolation
			if result != nil {
				for i, item := range result.Items {
					path := fmt.Sprintf("Items[%d]", i)
					walkNodeLocs(reflect.ValueOf(item), path, &vs)
				}
			}
			return vs
		}()

		totalStatements++
		if len(violations) > 0 {
			statementsWithViolations++
			totalViolations += len(violations)
			for _, v := range violations {
				violationsByTag[v.NodeTag] = append(violationsByTag[v.NodeTag], v)
			}
		}
	}

	// Sort tags for stable output
	var tags []string
	for tag := range violationsByTag {
		tags = append(tags, tag)
	}
	sort.Strings(tags)

	t.Logf("=== Loc.End Audit Summary ===")
	t.Logf("Total SQL statements tested: %d", totalStatements)
	t.Logf("Statements with violations: %d", statementsWithViolations)
	t.Logf("Total violations: %d", totalViolations)
	t.Logf("")
	t.Logf("=== Violations by Node Type ===")

	for _, tag := range tags {
		vs := violationsByTag[tag]
		t.Logf("")
		t.Logf("--- %s (%d violations) ---", tag, len(vs))
		// Show up to 5 example paths per tag
		shown := 0
		for _, v := range vs {
			if shown >= 5 {
				t.Logf("  ... and %d more", len(vs)-5)
				break
			}
			t.Logf("  %s", v)
			shown++
		}
	}

	// Also log a concise table
	t.Logf("")
	t.Logf("=== Concise Violation Counts ===")
	t.Logf("%-30s %s", "NodeTag", "Count")
	t.Logf("%-30s %s", strings.Repeat("-", 30), strings.Repeat("-", 5))
	for _, tag := range tags {
		t.Logf("%-30s %d", tag, len(violationsByTag[tag]))
	}
}
