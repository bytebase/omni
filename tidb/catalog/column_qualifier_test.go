package catalog

import (
	"errors"
	"testing"
)

// TestColumnQualifier holds what a column's qualifier must name, every result read by running the
// statement on TiDB v8.5.5 on 2026-11-05. CREATE TABLE, MODIFY and CHANGE check it, the database
// first. ADD, ALTER and DROP ignore it.
func TestColumnQualifier(t *testing.T) {
	const setup = "CREATE DATABASE d; USE d; CREATE TABLE t9 (a INT, b INT, c INT);"
	for _, c := range []struct {
		sql  string
		code int
	}{
		{"CREATE TABLE t1 (t1.a INT)", 0},
		{"CREATE TABLE t1 (d.t1.a INT)", 0},
		{"CREATE TABLE t1 (T1.a INT)", 0},
		{"CREATE TABLE t1 (x.a INT)", ErrWrongTableName},
		{"CREATE TABLE t1 (x.t1.a INT)", ErrWrongDBName},
		{"CREATE TABLE t1 (d.x.a INT)", ErrWrongTableName},
		{"CREATE TABLE t1 (x.y.a INT)", ErrWrongDBName},
		{"ALTER TABLE t9 ADD COLUMN t9.d INT", 0},
		{"ALTER TABLE t9 ADD COLUMN x.d INT", 0},
		{"ALTER TABLE t9 ADD COLUMN x.t9.d INT", 0},
		{"ALTER TABLE t9 MODIFY t9.c BIGINT", 0},
		{"ALTER TABLE t9 MODIFY x.c BIGINT", ErrWrongTableName},
		{"ALTER TABLE t9 CHANGE t9.c cc INT", 0},
		{"ALTER TABLE t9 CHANGE x.c cc INT", ErrWrongTableName},
		{"ALTER TABLE t9 CHANGE c t9.cc INT", 0},
		{"ALTER TABLE t9 CHANGE c x.cc INT", ErrWrongTableName},
		{"ALTER TABLE t9 ALTER t9.a SET DEFAULT 1", 0},
		{"ALTER TABLE t9 ALTER x.a SET DEFAULT 1", 0},
		{"ALTER TABLE t9 DROP COLUMN t9.b", 0},
		{"ALTER TABLE t9 DROP COLUMN x.b", 0},
	} {
		t.Run(c.sql, func(t *testing.T) {
			cat := New()
			mustExec(t, cat, setup)
			results, err := cat.Exec(c.sql, nil)
			if err != nil {
				t.Fatal(err)
			}
			err = results[len(results)-1].Error
			if c.code == 0 {
				if err != nil {
					t.Fatalf("want no error, got %v", err)
				}
				return
			}
			var catErr *Error
			if !errors.As(err, &catErr) || catErr.Code != c.code {
				t.Fatalf("want error %d, got %v", c.code, err)
			}
		})
	}

	// The column keeps its own name.
	cat := New()
	mustExec(t, cat, setup+"CREATE TABLE t1 (t1.a INT, d.t1.b INT);")
	tbl := cat.GetDatabase("d").Tables["t1"]
	if tbl == nil || tbl.GetColumn("a") == nil || tbl.GetColumn("b") == nil {
		t.Fatalf("columns not under their own names: %+v", tbl)
	}
}
