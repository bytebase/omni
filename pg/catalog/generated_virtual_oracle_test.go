package catalog

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestPostgres18VirtualGeneratedDDL(t *testing.T) {
	ctr := startPGContainer(t)
	schema := ctr.freshSchema(t)
	c := New()
	execSQL(t, c, "CREATE SCHEMA "+schema+"; SET search_path TO "+schema)
	ctr.execSQL(t, "SET search_path TO "+schema)
	for _, step := range []struct {
		sql, code string
	}{
		{`CREATE TABLE t(a int, b int GENERATED ALWAYS AS (a + 1))`, ""},
		{`ALTER TABLE t ADD COLUMN c int GENERATED ALWAYS AS (a + 2) VIRTUAL`, ""},
		{`CREATE TABLE copied(LIKE t INCLUDING GENERATED)`, ""},
		{`ALTER TABLE t ALTER COLUMN b SET EXPRESSION AS (a + 3)`, ""},
		{`ALTER TABLE t ALTER COLUMN b DROP EXPRESSION`, "0A000"},
		{`ALTER TABLE t ALTER COLUMN b DROP EXPRESSION IF EXISTS`, "0A000"},
		{`ALTER TABLE t ADD CHECK(a > 0)`, ""},
		{`ALTER TABLE t ALTER COLUMN b SET EXPRESSION AS (a + 4)`, "0A000"},
	} {
		_, nativeErr := ctr.db.ExecContext(ctr.ctx, step.sql)
		var native *pgconn.PgError
		nativeCode := ""
		if nativeErr != nil {
			if !errors.As(nativeErr, &native) {
				t.Fatal(nativeErr)
			}
			nativeCode = native.Code
		}
		if nativeCode != step.code {
			t.Fatalf("native %q: %v; want %s", step.sql, nativeErr, step.code)
		}
		modeledErr := execSQLErr(c, step.sql)
		if step.code == "" {
			if modeledErr != nil {
				t.Fatalf("catalog %q: %v", step.sql, modeledErr)
			}
		} else {
			assertErrorCode(t, modeledErr, step.code)
		}
		for _, table := range []string{"t", "copied"} {
			rel := c.GetRelation(schema, table)
			if rel == nil {
				continue
			}
			for _, col := range rel.Columns {
				var kind, expr string
				if err := ctr.db.QueryRowContext(ctr.ctx, `SELECT a.attgenerated, COALESCE(pg_get_expr(d.adbin,d.adrelid),'')
					FROM pg_attribute a LEFT JOIN pg_attrdef d ON d.adrelid=a.attrelid AND d.adnum=a.attnum
					WHERE a.attrelid=$1::regclass AND a.attname=$2`, schema+"."+table, col.Name).Scan(&kind, &expr); err != nil {
					t.Fatal(err)
				}
				modeledKind := ""
				if col.Generated != 0 {
					modeledKind = string(col.Generated)
				}
				if kind != modeledKind || expr != col.GenerationExpr {
					t.Fatalf("after %q, %s.%s: native (%q,%q); catalog (%q,%q)", step.sql, table, col.Name, kind, expr, modeledKind, col.GenerationExpr)
				}
			}
		}
	}
}
