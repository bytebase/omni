package catalog

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// What OceanBase takes and refuses of the statements the migration generator chooses
// between. The generator's rules cite these tests; each statement below is held to the
// pinned server, with the error number it answers, 0 for a statement it takes.

// taken is the error number of a statement OceanBase takes.
const taken = 0

type engineStep struct {
	sql   string
	errno uint16
}

// runEngineSteps runs each case's statements in a database of its own and compares the
// error number of every one.
func runEngineSteps(t *testing.T, o *oracleConn, prefix string, cases []struct {
	name  string
	steps []engineStep
}) {
	ctx := context.Background()
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := fmt.Sprintf("%s_%02d", prefix, i)
			for _, s := range []string{"DROP DATABASE IF EXISTS " + db, "CREATE DATABASE " + db, "USE " + db} {
				if _, err := o.db.ExecContext(ctx, s); err != nil {
					t.Fatalf("%s: %v", s, err)
				}
			}
			for _, step := range tc.steps {
				_, err := o.db.ExecContext(ctx, step.sql)
				var got uint16
				var myErr *gomysql.MySQLError
				switch {
				case errors.As(err, &myErr):
					got = myErr.Number
				case err != nil:
					t.Fatalf("%s: %v", step.sql, err)
				}
				if got != step.errno {
					t.Errorf("%s\n  answered %d (%v), recorded %d", step.sql, got, err, step.errno)
				}
			}
		})
	}
}

// TestOracle_OceanBaseAlterTableClauses holds to the engine what the column, index and
// partition generators rest on: which ALTER TABLE clauses OceanBase takes together, and
// which statement has to come first.
func TestOracle_OceanBaseAlterTableClauses(t *testing.T) {
	o := oceanBaseOracle(t)
	runEngineSteps(t, o, "obalter", []struct {
		name  string
		steps []engineStep
	}{
		{"an AUTO_INCREMENT column needs no key", []engineStep{
			{"CREATE TABLE t (id INT NOT NULL AUTO_INCREMENT, a INT)", taken},
			{"CREATE TABLE u (id INT NOT NULL PRIMARY KEY, a INT)", taken},
			{"ALTER TABLE u ADD COLUMN seq BIGINT UNSIGNED NOT NULL AUTO_INCREMENT", taken},
			{"ALTER TABLE u ADD UNIQUE KEY uk_seq (seq)", taken},
			{"ALTER TABLE u DROP INDEX uk_seq", taken},
		}},
		{"one AUTO_INCREMENT column at a time", []engineStep{
			{"CREATE TABLE t (a INT NOT NULL AUTO_INCREMENT, b INT NOT NULL, KEY ka (a), KEY kb (b))", taken},
			{"ALTER TABLE t MODIFY COLUMN b INT NOT NULL AUTO_INCREMENT", 1075},
			{"ALTER TABLE t MODIFY COLUMN a INT NOT NULL, MODIFY COLUMN b INT NOT NULL AUTO_INCREMENT", 1075},
			{"ALTER TABLE t MODIFY COLUMN a INT NOT NULL", taken},
			{"ALTER TABLE t MODIFY COLUMN b INT NOT NULL AUTO_INCREMENT", taken},
		}},
		{"clauses it does not take together", []engineStep{
			{"CREATE TABLE t (id INT NOT NULL, c INT NOT NULL, a INT, PRIMARY KEY (id), UNIQUE KEY uk_c (c), KEY ka (a))", taken},
			{"ALTER TABLE t ADD COLUMN seq BIGINT NOT NULL AUTO_INCREMENT, DROP PRIMARY KEY, ADD PRIMARY KEY (seq)", 1235},
			{"ALTER TABLE t DROP INDEX uk_c, ADD UNIQUE KEY uk_c (c, a)", 4179},
			{"ALTER TABLE t DROP INDEX ka, DROP COLUMN a", 1235},
			{"ALTER TABLE t MODIFY COLUMN c INT NOT NULL AUTO_INCREMENT, ADD UNIQUE KEY uk_c2 (c)", 1235},
			{"ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (c), DROP INDEX uk_c", 1235},
			{"ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (c)", taken},
			{"CREATE TABLE u (a INT)", taken},
			{"ALTER TABLE u ADD COLUMN seq BIGINT NOT NULL AUTO_INCREMENT, ADD PRIMARY KEY (seq)", 1176},
			{"ALTER TABLE u ADD COLUMN seq BIGINT NOT NULL AUTO_INCREMENT", taken},
			{"ALTER TABLE u ADD PRIMARY KEY (seq)", taken},
		}},
		{"a column the primary key holds", []engineStep{
			{"CREATE TABLE t (a INT NOT NULL, c INT NOT NULL, b INT NOT NULL, PRIMARY KEY (c, a, b))", taken},
			{"ALTER TABLE t DROP COLUMN a", 1235},
			{"ALTER TABLE t MODIFY COLUMN b INT DEFAULT NULL", 1171},
			{"ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (c)", taken},
			{"ALTER TABLE t DROP COLUMN a", taken},
			{"ALTER TABLE t MODIFY COLUMN b INT DEFAULT NULL", taken},
			{"ALTER TABLE t DROP COLUMN c", 1235},
			{"ALTER TABLE t DROP PRIMARY KEY", taken},
			{"ALTER TABLE t DROP COLUMN c", taken},
		}},
		{"a column an index holds goes with it", []engineStep{
			{"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, c INT, a INT, UNIQUE KEY uk_c (c), KEY k_ca (c, a))", taken},
			{"ALTER TABLE t DROP COLUMN c", taken},
			{"ALTER TABLE t DROP INDEX uk_c", 1091},
			{"ALTER TABLE t DROP INDEX k_ca", taken},
		}},
		{"partitioning is replaced, not removed", []engineStep{
			{"CREATE TABLE t (id INT NOT NULL, dt DATE NOT NULL, PRIMARY KEY (id, dt)) PARTITION BY RANGE (YEAR(dt)) (PARTITION p0 VALUES LESS THAN (2010), PARTITION pmax VALUES LESS THAN MAXVALUE)", taken},
			{"ALTER TABLE t REMOVE PARTITIONING", 1064},
			{"ALTER TABLE t DROP COLUMN dt", 3855},
			{"ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (id)", 1503},
			{"ALTER TABLE t PARTITION BY HASH (id) PARTITIONS 4", taken},
			{"ALTER TABLE t DROP PRIMARY KEY, ADD PRIMARY KEY (id)", taken},
			{"ALTER TABLE t DROP COLUMN dt", taken},
			{"ALTER TABLE t PARTITION BY HASH (id) PARTITIONS 1", taken},
			{"ALTER TABLE t DROP PARTITION p0", 1235},
		}},
		{"a spatial column", []engineStep{
			{"CREATE TABLE t (id INT NOT NULL PRIMARY KEY)", taken},
			{"ALTER TABLE t ADD COLUMN g GEOMETRY NOT NULL", 1138},
			{"ALTER TABLE t ADD COLUMN p POINT NOT NULL /*!80003 SRID 0 */", 1138},
			{"ALTER TABLE t ADD COLUMN j JSON NOT NULL", taken},
			{"ALTER TABLE t ADD COLUMN g GEOMETRY /*!80003 SRID 0 */", taken},
			{"ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 0 */", taken},
			{"ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 4326 */", 1235},
			{"ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL", 1235},
			// NOT NULL already: no MODIFY is taken, whatever it changes, but through nullable.
			{"ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 0 */", 1138},
			{"ALTER TABLE t MODIFY COLUMN g GEOMETRY NOT NULL /*!80003 SRID 0 */ COMMENT 'c'", 1138},
			{"ALTER TABLE t MODIFY COLUMN g POINT NOT NULL /*!80003 SRID 0 */", 1138},
			{"ALTER TABLE t MODIFY COLUMN g POINT /*!80003 SRID 0 */ COMMENT 'c'", taken},
			{"ALTER TABLE t MODIFY COLUMN g POINT NOT NULL /*!80003 SRID 0 */ COMMENT 'c'", taken},
			{"ALTER TABLE t MODIFY COLUMN j JSON NOT NULL COMMENT 'c'", taken},
		}},
		{"a generated column's expression", []engineStep{
			{"CREATE TABLE t (id INT NOT NULL PRIMARY KEY, a INT, g INT GENERATED ALWAYS AS (a + 1) STORED, v INT GENERATED ALWAYS AS (a + 1) VIRTUAL, KEY kg (g))", taken},
			{"ALTER TABLE t MODIFY COLUMN g INT GENERATED ALWAYS AS (a + 2) STORED", 1235},
			{"ALTER TABLE t MODIFY COLUMN v INT GENERATED ALWAYS AS (a + 2) VIRTUAL", 1235},
			{"ALTER TABLE t MODIFY COLUMN g INT GENERATED ALWAYS AS (a + 1) STORED COMMENT 'c'", taken},
			{"ALTER TABLE t DROP COLUMN a", 3108},
			{"ALTER TABLE t DROP COLUMN g", taken},
			{"ALTER TABLE t DROP INDEX kg", 1091},
			{"ALTER TABLE t ADD COLUMN g INT GENERATED ALWAYS AS (a + 2) STORED", taken},
			{"ALTER TABLE t ADD KEY kg (g)", taken},
		}},
	})
}

// TestOracle_OceanBaseForeignKeyIndex holds to the engine what the loader and the foreign
// key generator rest on: OceanBase makes no index for a foreign key, and while the foreign
// key lives it holds the last index over exactly its columns.
func TestOracle_OceanBaseForeignKeyIndex(t *testing.T) {
	o := oceanBaseOracle(t)
	const parent = "CREATE TABLE p (id INT NOT NULL, u INT NOT NULL, v INT NOT NULL, PRIMARY KEY (id), UNIQUE KEY uu (u), UNIQUE KEY uv (u, v))"
	runEngineSteps(t, o, "obfk", []struct {
		name  string
		steps []engineStep
	}{
		{"no index is made, so there is none to drop", []engineStep{
			{parent, taken},
			{"CREATE TABLE c (id INT PRIMARY KEY, pid INT, x INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id), FOREIGN KEY (pid) REFERENCES p (id))", taken},
			{"ALTER TABLE c ADD CONSTRAINT fkx FOREIGN KEY (x) REFERENCES p (id)", taken},
			{"ALTER TABLE c DROP INDEX fk", 1091},
			{"ALTER TABLE c DROP INDEX pid", 1091},
			{"ALTER TABLE c DROP INDEX fkx", 1091},
			{"ALTER TABLE c DROP FOREIGN KEY fk, DROP INDEX fk", 1091},
		}},
		{"the last index over its columns is held", []engineStep{
			{parent, taken},
			{"CREATE TABLE c (id INT PRIMARY KEY, pid INT, x INT, KEY a1 (pid), KEY a2 (pid), KEY wide (pid, x), CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))", taken},
			{"ALTER TABLE c DROP INDEX wide", taken},
			{"ALTER TABLE c DROP INDEX a2", taken},
			{"ALTER TABLE c DROP INDEX a1", 5320},
			{"ALTER TABLE c DROP FOREIGN KEY fk", taken},
			{"ALTER TABLE c DROP INDEX a1", taken},
			{"ALTER TABLE c ADD CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id)", taken},
		}},
		{"an index added after the foreign key is held too", []engineStep{
			{parent, taken},
			{"CREATE TABLE c (id INT PRIMARY KEY, pid INT, CONSTRAINT fk FOREIGN KEY (pid) REFERENCES p (id))", taken},
			{"ALTER TABLE c ADD UNIQUE KEY late (pid)", taken},
			{"ALTER TABLE c DROP INDEX late", 5320},
		}},
		{"in any column order, and not a prefix", []engineStep{
			{parent, taken},
			{"CREATE TABLE c (id INT PRIMARY KEY, a INT, b INT, s VARCHAR(20), KEY ba (b, a), KEY pre (s(5)), CONSTRAINT fab FOREIGN KEY (a, b) REFERENCES p (u, v))", taken},
			{"ALTER TABLE c DROP INDEX pre", taken},
			{"ALTER TABLE c DROP INDEX ba", 5320},
		}},
		{"its own primary key is held, and can be replaced", []engineStep{
			{parent, taken},
			{"CREATE TABLE c (id INT NOT NULL, a INT NOT NULL, PRIMARY KEY (id), CONSTRAINT fk FOREIGN KEY (id) REFERENCES p (id))", taken},
			{"ALTER TABLE c DROP PRIMARY KEY", 4714},
			{"ALTER TABLE c DROP PRIMARY KEY, ADD PRIMARY KEY (a)", taken},
		}},
		{"a referenced table keeps a primary key, and can replace it", []engineStep{
			{parent, taken},
			{"CREATE TABLE c (id INT PRIMARY KEY, pu INT, CONSTRAINT fk FOREIGN KEY (pu) REFERENCES p (u))", taken},
			{"ALTER TABLE p DROP PRIMARY KEY", 4714},
			{"ALTER TABLE p DROP INDEX uu", 4714},
			{"ALTER TABLE p DROP PRIMARY KEY, ADD PRIMARY KEY (id, v)", taken},
		}},
	})
}
