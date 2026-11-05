package catalog

import (
	"strings"
	"testing"
)

const floatDoubleForms = "CREATE TABLE fd (" +
	"a DOUBLE(10,0), b REAL(10,0), c FLOAT8(10,0), d DOUBLE PRECISION(10,0), e DOUBLE(10,0) UNSIGNED, " +
	"f DOUBLE(10,2), g DOUBLE, h FLOAT8, " +
	"i FLOAT(10,0), j FLOAT4(10,0), k FLOAT(10,2), l FLOAT(10), m FLOAT4(5), n FLOAT(24), o FLOAT(25), " +
	"p FLOAT(30), q FLOAT(0), r FLOAT, s FLOAT(10,0) ZEROFILL, t DOUBLE(10,2) ZEROFILL, u INT ZEROFILL, v DECIMAL(5,0))"

// TestFloatDoubleShowCreate: a precision without a scale does not show, so what the catalog writes
// is what the server writes, and it reads back: the parser refuses DOUBLE(10), as TiDB does.
func TestFloatDoubleShowCreate(t *testing.T) {
	tc := startTiDBForCatalog(t)
	sharedTiDB.Lock()
	defer sharedTiDB.Unlock()
	ctx := tc.ctx
	conn, err := tc.db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for _, stmt := range []string{"DROP DATABASE IF EXISTS zz_float", "CREATE DATABASE zz_float", "USE zz_float", floatDoubleForms} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}
	defer func() { _, _ = conn.ExecContext(ctx, "DROP DATABASE IF EXISTS zz_float") }()
	var table, server string
	if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE fd").Scan(&table, &server); err != nil {
		t.Fatal(err)
	}
	c := New()
	if _, err := c.Exec("CREATE DATABASE d; USE d; "+floatDoubleForms, nil); err != nil {
		t.Fatal(err)
	}
	got := c.ShowCreateTable("d", "fd")
	serverLines := strings.Split(server, "\n")
	for i, line := range strings.Split(got, "\n") {
		if strings.HasPrefix(line, "  `") && (i >= len(serverLines) || line != serverLines[i]) {
			t.Errorf("omni: %s\nserver: %s", line, serverLines[i])
		}
	}
	again := New()
	if _, err := again.Exec("CREATE DATABASE d; USE d; "+got, nil); err != nil {
		t.Fatalf("%s does not read back: %v", got, err)
	}
}
