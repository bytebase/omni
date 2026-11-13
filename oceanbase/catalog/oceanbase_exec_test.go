package catalog

import (
	"context"
	"errors"
	"fmt"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// What OceanBase answers to DDL the catalog simulates, where its answer is not MySQL's.
// Each statement below is held to the pinned server and to the catalog, with the error
// number both must answer, 0 for a statement they take. The catalog's rules cite these
// tests.

// runExecSteps runs each case's statements in a database of its own, on OceanBase and on
// a catalog, and compares the error number of every one on both.
func runExecSteps(t *testing.T, o *oracleConn, prefix string, cases []struct {
	name  string
	steps []engineStep
}) {
	ctx := context.Background()
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			db := fmt.Sprintf("%s_%02d", prefix, i)
			c := New()
			for _, s := range []string{"DROP DATABASE IF EXISTS " + db, "CREATE DATABASE " + db, "USE " + db} {
				if _, err := o.db.ExecContext(ctx, s); err != nil {
					t.Fatalf("%s: %v", s, err)
				}
				if code := omniCode(c, s); code != 0 {
					t.Fatalf("omni: %s: error %d", s, code)
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
					t.Errorf("%s\n  OceanBase answered %d (%v), recorded %d", step.sql, got, err, step.errno)
				}
				if code := omniCode(c, step.sql); code != int(step.errno) {
					t.Errorf("%s\n  the catalog answered %d, OceanBase %d", step.sql, code, step.errno)
				}
			}
		})
	}
}

// TestOracle_OceanBaseKeyErrors holds to the engine the errors of a second primary key and
// of a foreign key OceanBase will not add, and the foreign keys it takes that MySQL refuses.
func TestOracle_OceanBaseKeyErrors(t *testing.T) {
	o := oceanBaseOracle(t)
	runExecSteps(t, o, "obkey", []struct {
		name  string
		steps []engineStep
	}{
		{"a second primary key", []engineStep{
			{"CREATE TABLE t (a INT, b INT, PRIMARY KEY (a), PRIMARY KEY (b))", 1062},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT PRIMARY KEY)", 1068},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT, PRIMARY KEY (b))", 1062},
			{"CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL)", taken},
			{"ALTER TABLE t ADD PRIMARY KEY (b)", 1068},
			{"ALTER TABLE t ADD COLUMN c INT PRIMARY KEY", 1068},
		}},
		{"a foreign key to nothing", []engineStep{
			{"CREATE TABLE c (pid INT, FOREIGN KEY (pid) REFERENCES nosuch (id))", 1146},
			{"CREATE TABLE p (id INT PRIMARY KEY, v INT)", taken},
			{"CREATE TABLE c (pid INT, FOREIGN KEY (pid) REFERENCES p (nosuch))", 5031},
			{"CREATE TABLE c (pid INT, FOREIGN KEY (nosuch) REFERENCES p (id))", 5031},
			{"CREATE TABLE c (pid INT)", taken},
			{"ALTER TABLE c ADD FOREIGN KEY (pid) REFERENCES nosuch (id)", 1146},
			{"ALTER TABLE c ADD FOREIGN KEY (pid) REFERENCES p (nosuch)", 5031},
			{"ALTER TABLE c ADD FOREIGN KEY (nosuch) REFERENCES p (id)", 5031},
		}},
		// The referenced columns lead an index, unique or not, in any order.
		{"a foreign key to columns without a key", []engineStep{
			{"CREATE TABLE p (id INT, v INT, w INT, x INT, y INT, KEY kv (v), UNIQUE KEY uw (w), UNIQUE KEY uxy (x, y))", taken},
			{"CREATE TABLE c1 (pid INT, FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c2 (pid INT, FOREIGN KEY (pid) REFERENCES p (v))", taken},
			{"CREATE TABLE c3 (pid INT, FOREIGN KEY (pid) REFERENCES p (w))", taken},
			{"CREATE TABLE c4 (pid INT, FOREIGN KEY (pid) REFERENCES p (x))", taken},
			{"CREATE TABLE c5 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (x, y))", taken},
			{"CREATE TABLE c6 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (y, x))", taken},
			{"CREATE TABLE c7 (pid INT)", taken},
			{"ALTER TABLE c7 ADD FOREIGN KEY (pid) REFERENCES p (v)", taken},
			{"ALTER TABLE c7 ADD FOREIGN KEY (pid) REFERENCES p (id)", 1215},
			{"CREATE TABLE c8 (pid INT, FOREIGN KEY (pid) REFERENCES p (y))", 1215},
			{"CREATE TABLE c9 (a INT, b INT, FOREIGN KEY (a, b) REFERENCES p (v, w))", 1215},
			{"CREATE TABLE q (a INT, b INT, c INT, PRIMARY KEY (a, b), KEY kc (c, a))", taken},
			{"CREATE TABLE d1 (x INT, FOREIGN KEY (x) REFERENCES q (b))", 1215},
			{"CREATE TABLE d2 (x INT, FOREIGN KEY (x) REFERENCES q (a))", taken},
			{"CREATE TABLE d3 (x INT, y INT, FOREIGN KEY (x, y) REFERENCES q (a, c))", taken},
			{"CREATE TABLE d4 (x INT, y INT, z INT, FOREIGN KEY (x, y, z) REFERENCES q (a, b, c))", 1215},
			{"CREATE TABLE s (id INT PRIMARY KEY, pid INT, FOREIGN KEY (pid) REFERENCES s (id))", taken},
			{"CREATE TABLE s2 (id INT, pid INT, FOREIGN KEY (pid) REFERENCES s2 (id))", 1215},
		}},
		{"a foreign key between types", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (pid VARCHAR(50), FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c2 (pid BIGINT, FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c3 (pid INT UNSIGNED, FOREIGN KEY (pid) REFERENCES p (id))", 1215},
			{"CREATE TABLE c4 (pid INT)", taken},
			{"ALTER TABLE c4 ADD FOREIGN KEY (pid) REFERENCES p (id)", taken},
			{"CREATE TABLE c5 (pid BIGINT)", taken},
			{"ALTER TABLE c5 ADD FOREIGN KEY (pid) REFERENCES p (id)", 1215},
		}},
		{"SET DEFAULT is no action of a foreign key", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (pid INT DEFAULT 0, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET DEFAULT)", 1215},
			{"CREATE TABLE c2 (pid INT DEFAULT 0, FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET DEFAULT)", 1215},
			{"CREATE TABLE c3 (pid INT DEFAULT 0)", taken},
			{"ALTER TABLE c3 ADD FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET DEFAULT", 1215},
		}},
		{"SET NULL needs a nullable column", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (pid INT NOT NULL, FOREIGN KEY (pid) REFERENCES p (id) ON DELETE SET NULL)", 1830},
			{"CREATE TABLE c2 (pid INT NOT NULL, FOREIGN KEY (pid) REFERENCES p (id) ON UPDATE SET NULL)", 1830},
		}},
		{"a foreign key over a generated column", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY, v INT GENERATED ALWAYS AS (id + 1) VIRTUAL, s INT GENERATED ALWAYS AS (id + 1) STORED, UNIQUE KEY uv (v), UNIQUE KEY us (s))", taken},
			{"CREATE TABLE c1 (a INT, b INT GENERATED ALWAYS AS (a + 1) VIRTUAL, FOREIGN KEY (b) REFERENCES p (id))", taken},
			{"CREATE TABLE c2 (a INT, b INT GENERATED ALWAYS AS (a + 1) STORED, FOREIGN KEY (b) REFERENCES p (id))", taken},
			{"CREATE TABLE c3 (a INT, FOREIGN KEY (a) REFERENCES p (v))", taken},
			{"CREATE TABLE c4 (a INT, FOREIGN KEY (a) REFERENCES p (s))", taken},
		}},
		{"a foreign key's name twice", []engineStep{
			{"CREATE TABLE p (id INT PRIMARY KEY)", taken},
			{"CREATE TABLE c1 (a INT, b INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id), CONSTRAINT fk FOREIGN KEY (b) REFERENCES p (id))", 1215},
			{"CREATE TABLE c2 (a INT, b INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id))", taken},
			{"CREATE TABLE c3 (a INT, CONSTRAINT fk FOREIGN KEY (a) REFERENCES p (id))", 1022},
			{"ALTER TABLE c2 ADD CONSTRAINT fk FOREIGN KEY (b) REFERENCES p (id)", 1022},
		}},
	})
}
