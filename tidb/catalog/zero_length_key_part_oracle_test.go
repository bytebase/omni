package catalog

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestZeroLengthKeyPartMatchesTiDB holds the catalog to pingcap/tidb:v8.5.5 on
// an explicit zero key-part length, KEY (c(0)): 1391 on every key path, before
// the table is looked up, while a positive prefix is taken.
func TestZeroLengthKeyPartMatchesTiDB(t *testing.T) {
	tc := startTiDBForCatalog(t)
	takeDatabaseOnSharedTiDB(t, tc)
	c := setupWithDB(t)
	setup := "CREATE TABLE zl_t (s VARCHAR(10))"
	mustExecTiDB(t, tc, setup)
	mustExec(t, c, setup)
	for _, tc2 := range []struct {
		sql  string
		code uint16
	}{
		{"CREATE TABLE zl_a (s VARCHAR(10), KEY (s(0)))", 1391},
		{"CREATE TABLE zl_b (s VARCHAR(10), KEY (nosuch(0)))", 1391},
		{"CREATE TABLE zl_c (a INT, KEY (a(0)))", 1391},
		{"ALTER TABLE zl_t ADD KEY (s(0))", 1391},
		{"CREATE INDEX zl_i ON zl_t (s(0))", 1391},
		{"CREATE INDEX zl_i ON zl_nosuch (s(0))", 1391},
		{"CREATE TABLE zl_d (s VARCHAR(10), KEY (s(5)))", 0},
	} {
		var serverCode uint16
		if _, err := tc.db.ExecContext(tc.ctx, tc2.sql); err != nil {
			var me *gomysql.MySQLError
			if !errors.As(err, &me) {
				t.Fatalf("tidb %s: %v", tc2.sql, err)
			}
			serverCode = me.Number
		}
		var catalogCode uint16
		if err := execOne(t, c, tc2.sql); err != nil {
			var ce *Error
			if !errors.As(err, &ce) {
				t.Fatalf("catalog %s: %v", tc2.sql, err)
			}
			catalogCode = uint16(ce.Code)
		}
		if serverCode != tc2.code || catalogCode != tc2.code {
			t.Errorf("%s: tidb %d, catalog %d, want %d", tc2.sql, serverCode, catalogCode, tc2.code)
		}
	}
}
